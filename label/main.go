package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	kafka "github.com/opensourceways/kafka-lib/agent"
	"github.com/opensourceways/server-common-lib/logrusutil"
	liboptions "github.com/opensourceways/server-common-lib/options"
	"github.com/sirupsen/logrus"

	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/label/app"
	"github.com/openmerlin/merlin-server/label/infrastructure/giteaimpl"
	"github.com/openmerlin/merlin-server/models/infrastructure/modelrepositoryadapter"
	"github.com/openmerlin/merlin-server/models/messageapp"
)

const component = "merlin-label"

type options struct {
	service   liboptions.ServiceOptions
	removeCfg bool
}

func (o *options) Validate() error {
	return o.service.Validate()
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options

	o.service.AddFlags(fs)

	fs.BoolVar(
		&o.removeCfg, "rm-cfg", false,
		"whether remove the cfg file after initialized.",
	)

	if err := fs.Parse(args); err != nil {
		fs.PrintDefaults()

		logrus.Fatalf("failed to parse cmdline %s", err)
	}

	return o
}

func main() {
	logrusutil.ComponentInit(component)
	log := logrus.NewEntry(logrus.StandardLogger())

	o := gatherOptions(flag.NewFlagSet(os.Args[0], flag.ExitOnError), os.Args[1:]...)
	if err := o.Validate(); err != nil {
		logrus.Errorf("Invalid options, err:%s", err.Error())

		return
	}

	// cfg
	cfg, err := LoadConfig(o.service.ConfigFile, o.removeCfg)
	if err != nil {
		logrus.Errorf("load config failed, err:%s", err.Error())

		return
	}

	// kafka
	if err = kafka.Init(&cfg.Kafka, log, nil, "", o.removeCfg); err != nil {
		logrus.Errorf("init kafka failed, err:%s", err.Error())

		return
	}

	defer kafka.Exit()

	// postgresql
	if err = postgresql.Init(&cfg.Postgresql, o.removeCfg); err != nil {
		logrus.Errorf("init postgresql failed, err:%s", err.Error())

		return
	}

	err = modelrepositoryadapter.Init(postgresql.DB(), &cfg.Model.Tables)
	if err != nil {
		return
	}

	if err = cfg.InitPrimitive(); err != nil {
		logrus.Errorf("init primitive failed, err:%s", err.Error())

		return
	}

	// run
	run(cfg)
}

func run(cfg *Config) {
	handler := app.NewLabelHandler(giteaimpl.NewGiteaImpl(&cfg.Gitea))

	modelsApp := messageapp.NewModelAppService(modelrepositoryadapter.ModelLabelsAdapter())

	message := NewMessageServer(handler, modelsApp, cfg.UserAgent)

	err := kafka.Subscribe(
		component,
		message.handle,
		[]string{cfg.Topics.MerlinHookEvent})
	if err != nil {
		logrus.Errorf("subscribe topic failed, err:%s", err.Error())

		return
	}

	// wait
	wait()
}

func wait() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	defer wg.Wait()

	called := false
	ctx, done := context.WithCancel(context.Background())

	defer func() {
		if !called {
			called = true
			done()
		}
	}()

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()

		select {
		case <-ctx.Done():
			logrus.Info("receive done. exit normally")
			return

		case <-sig:
			logrus.Info("receive exit signal")
			called = true
			done()
			return
		}
	}(ctx)

	<-ctx.Done()
}

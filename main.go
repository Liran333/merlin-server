package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/opensourceways/community-robot-lib/logrusutil"
	"github.com/sirupsen/logrus"

	"github.com/openmerlin/merlin-server/common/infrastructure/redis"
	"github.com/openmerlin/merlin-server/config"
	"github.com/openmerlin/merlin-server/controller"
	"github.com/openmerlin/merlin-server/infrastructure/mongodb"
	"github.com/openmerlin/merlin-server/login/infrastructure/oidcimpl"
	"github.com/openmerlin/merlin-server/server"
	"github.com/openmerlin/merlin-server/user/domain"
)

type options struct {
	service     ServiceOptions
	enableDebug bool
}

type ServiceOptions struct {
	Port        int
	ConfigFile  string
	GracePeriod time.Duration
	RemoveCfg   bool
}

func (o *ServiceOptions) Validate() error {
	if o.ConfigFile == "" {
		return fmt.Errorf("missing config-file")
	}

	return nil
}

func (o *ServiceOptions) AddFlags(fs *flag.FlagSet) {
	fs.IntVar(&o.Port, "port", 8888, "Port to listen on.")
	fs.BoolVar(&o.RemoveCfg, "rm-cfg", true, "whether remove the cfg file after initialized .")

	fs.StringVar(&o.ConfigFile, "config-file", "", "Path to config file.")
	fs.DurationVar(&o.GracePeriod, "grace-period", 180*time.Second, "On shutdown, try to handle remaining events for the specified duration.")
}

func (o *options) Validate() error {
	return o.service.Validate()
}

func gatherOptions(fs *flag.FlagSet, args ...string) (options, error) {
	var o options

	o.service.AddFlags(fs)

	fs.BoolVar(
		&o.enableDebug, "enable_debug", false,
		"whether to enable debug model.",
	)

	err := fs.Parse(args)

	return o, err
}

func main() {
	logrusutil.ComponentInit("merlin")
	log := logrus.NewEntry(logrus.StandardLogger())

	o, err := gatherOptions(
		flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		os.Args[1:]...,
	)
	if err != nil {
		logrus.Fatalf("new options failed, err:%s", err.Error())
	}

	if err := o.Validate(); err != nil {
		logrus.Fatalf("Invalid options, err:%s", err.Error())
	}

	if o.enableDebug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug enabled.")
	}

	// cfg
	cfg := new(config.Config)
	if err := config.LoadConfig(o.service.ConfigFile, cfg, o.service.RemoveCfg); err != nil {
		logrus.Fatalf("load config, err:%s", err.Error())
	}

	// authing
	oidcimpl.Init(&cfg.Authing)

	// controller
	api := &cfg.API

	if err := controller.Init(api, log); err != nil {
		logrus.Fatalf("initialize api controller failed, err:%s", err.Error())
	}

	// mongo
	m := &cfg.Mongodb
	if err := mongodb.Initialize(m.DBConn, m.DBName, m.DBCert, o.service.RemoveCfg); err != nil {
		logrus.Fatalf("initialize mongodb failed, err:%s", err.Error())
	}

	defer mongodb.Close()
	//redis
	if err := redis.Init(&cfg.Redis.DB, o.service.RemoveCfg); err != nil {
		logrus.Fatalf("init redis failed, err:%s", err.Error())
	}

	// user
	domain.Init(&cfg.User)

	// run
	server.StartWebServer(o.service.Port, o.service.GracePeriod, cfg)
}

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	redisdb "github.com/opensourceways/redis-lib"
	"github.com/sirupsen/logrus"

	"github.com/openmerlin/merlin-server/common/infrastructure/gitea"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/config"
	"github.com/openmerlin/merlin-server/server"
)

type options struct {
	service     ServiceOptions
	enableDebug bool
}

type ServiceOptions struct {
	Port        int
	ConfigFile  string
	Cert        string
	Key         string
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
	fs.BoolVar(&o.RemoveCfg, "rm-cfg", false, "whether remove the cfg file after initialized .")

	fs.StringVar(&o.ConfigFile, "config-file", "", "Path to config file.")
	fs.StringVar(&o.Cert, "cert", "", "Path to tls cert file.")
	fs.StringVar(&o.Key, "key", "", "Path to tls key file.")
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

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and api Bearer.
func main() {
	o, err := gatherOptions(
		flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		os.Args[1:]...,
	)
	if err != nil {
		logrus.Errorf("new options failed, err:%s", err.Error())

		return
	}

	if err := o.Validate(); err != nil {
		logrus.Errorf("Invalid options, err:%s", err.Error())

		return
	}

	if o.enableDebug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug enabled.")
	}

	// cfg
	cfg := new(config.Config)

	if err := config.LoadConfig(o.service.ConfigFile, cfg, o.service.RemoveCfg); err != nil {
		logrus.Errorf("load config, err:%s", err.Error())

		return
	}

	//redis
	if err := redisdb.Init(&cfg.Redis, o.service.RemoveCfg); err != nil {
		logrus.Errorf("init redis failed, err:%s", err.Error())

		return
	}

	defer redisdb.Close()

	// postgresql
	if err := postgresql.Init(&cfg.Postgresql, o.service.RemoveCfg); err != nil {
		logrus.Errorf("init postgresql failed, err:%s", err.Error())

		return
	}

	// gitea
	if err := gitea.Init(&cfg.Git); err != nil {
		logrus.Errorf("init gitea failed, err:%s", err.Error())

		return
	}

	// init cfg
	if err := cfg.Init(); err != nil {
		logrus.Errorf("init cfg failed, err:%s", err.Error())
	}

	// session
	if err := cfg.InitSession(); err != nil {
		logrus.Errorf("init session failed, err:%s", err.Error())

		return
	}

	// run
	server.StartWebServer(o.service.Key, o.service.Cert, o.service.Port, o.service.GracePeriod, cfg)
}

package main

import (
	"errors"
	"fmt"

	"github.com/opensourceways/community-robot-lib/logrusutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	basegitea "github.com/openmerlin/merlin-server/common/infrastructure/gitea"
	"github.com/openmerlin/merlin-server/common/infrastructure/redis"
	"github.com/openmerlin/merlin-server/config"
	"github.com/openmerlin/merlin-server/infrastructure/mongodb"
	"github.com/openmerlin/merlin-server/user/domain"
)

var configFile string
var cfg *config.Config

func Error(cmd *cobra.Command, args []string, err error) {
	logrus.Fatalf("execute %s args:%v error:%v\n", cmd.Name(), args, err)
}

var rootCmd = &cobra.Command{
	Use:   "merlin-admin",
	Short: "merlin-admin is a admin tool for merlin server.",
	Run: func(cmd *cobra.Command, args []string) {
		Error(cmd, args, errors.New("unrecognized command"))
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
	},
}

func initServer(configFile string) {
	logrusutil.ComponentInit("admin")

	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debug("debug enabled.")

	// cfg
	cfg = new(config.Config)

	if err := config.LoadConfig(configFile, cfg, false); err != nil {
		logrus.Fatalf("load config, err:%s", err.Error())
	}

	// mongo
	m := &cfg.Mongodb
	if err := mongodb.Initialize(m.DBConn, m.DBName, m.DBCert, false); err != nil {
		logrus.Fatalf("initialize mongodb failed, err:%s", err.Error())
	}

	//redis
	if err := redis.Init(&cfg.Redis, false); err != nil {
		logrus.Fatalf("init redis failed, err:%s", err.Error())
	}

	// user
	domain.Init(&cfg.User)

	fmt.Printf("gitea: %+v\n", cfg.Git)

	// gitea
	if err := basegitea.Init(&cfg.Git); err != nil {
		logrus.Fatalf("init gitea failed, err:%s", err.Error())
	}
}

func execute() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.yaml", "config file path")

	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(tokenCmd)
	rootCmd.AddCommand(orgCmd)
	_ = rootCmd.Execute()
}

func main() {
	execute()
}

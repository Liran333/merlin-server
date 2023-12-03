package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	gitea "github.com/openmerlin/merlin-server/infrastructure/gitea"
	"github.com/opensourceways/community-robot-lib/logrusutil"
	"github.com/sirupsen/logrus"

	"github.com/openmerlin/merlin-server/common/infrastructure/redis"
	"github.com/openmerlin/merlin-server/config"
	"github.com/openmerlin/merlin-server/infrastructure/mongodb"
	"github.com/openmerlin/merlin-server/login/infrastructure/oidcimpl"

	userapp "github.com/openmerlin/merlin-server/user/app"
	"github.com/openmerlin/merlin-server/user/domain"
	usergit "github.com/openmerlin/merlin-server/user/infrastructure/git"
	userrepoimpl "github.com/openmerlin/merlin-server/user/infrastructure/repositoryimpl"
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
	username    string
	email       string
}

func (o *ServiceOptions) Validate() error {
	if o.ConfigFile == "" {
		return fmt.Errorf("missing config-file")
	}

	if o.username == "" {
		return fmt.Errorf("missing username")
	}

	if o.email == "" {
		return fmt.Errorf("missing email")
	}
	return nil
}

func (o *ServiceOptions) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&o.ConfigFile, "config-file", "", "Path to config file.")
	fs.StringVar(&o.username, "username", "", "username.")
	fs.StringVar(&o.email, "email", "", "email.")
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
	logrusutil.ComponentInit("admin")

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

	collections := &cfg.Mongodb.Collections

	// authing
	oidcimpl.Init(&cfg.Authing)

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

	fmt.Printf("gitea: %+v", cfg.Git)
	// gitea
	if err := gitea.Init(&cfg.Git); err != nil {
		logrus.Fatalf("init gitea failed, err:%s", err.Error())
	}

	acc, _ := domain.NewAccount(o.service.username)
	email, _ := domain.NewEmail(o.service.email)
	b, _ := domain.NewBio("testb")
	ava, _ := domain.NewAvatarId("1")

	user := userrepoimpl.NewUserRepo(
		mongodb.NewCollection(collections.User),
	)

	git := usergit.NewUserGit(gitea.GetClient())

	userAppService := userapp.NewUserService(
		user, git)

	_, err = userAppService.Create(&domain.UserCreateCmd{
		Email:   email,
		Account: acc,

		Bio:      b,
		AvatarId: ava,
	})
	if err != nil {
		logrus.Errorf("create user failed :%s", err.Error())
	} else {
		logrus.Info("create user successfully")
	}

}

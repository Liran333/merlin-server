/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package main is the entry point for the application.
package main

import (
	"context"
	"errors"
	"fmt"

	redisdb "github.com/opensourceways/redis-lib"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/openmerlin/merlin-server/common/domain/crypto"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/common/infrastructure/email"
	"github.com/openmerlin/merlin-server/common/infrastructure/gitea"
	"github.com/openmerlin/merlin-server/common/infrastructure/obs"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/infrastructure/giteauser"
	orgrepoimpl "github.com/openmerlin/merlin-server/organization/infrastructure/repositoryimpl"
	sessionapp "github.com/openmerlin/merlin-server/session/app"
	"github.com/openmerlin/merlin-server/session/infrastructure/loginrepositoryadapter"
	"github.com/openmerlin/merlin-server/session/infrastructure/oidcimpl"
	"github.com/openmerlin/merlin-server/session/infrastructure/sessionrepositoryadapter"
	userapp "github.com/openmerlin/merlin-server/user/app"
	"github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/user/domain/repository"
	"github.com/openmerlin/merlin-server/user/infrastructure/emailadapter"
	usergit "github.com/openmerlin/merlin-server/user/infrastructure/git"
	"github.com/openmerlin/merlin-server/user/infrastructure/obsadapter"
	userrepoimpl "github.com/openmerlin/merlin-server/user/infrastructure/repositoryimpl"
)

var userAppService userapp.UserService
var userrepo repository.User

func inittoken() {
	userrepo = userrepoimpl.NewUserRepo(
		postgresql.DAO(cfg.User.Domain.Tables.User), crypto.NewEncryption(cfg.User.Domain.Key),
	)

	member := orgrepoimpl.NewMemberRepo(postgresql.DAO(cfg.Org.Domain.Tables.Member))

	git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))
	t := userrepoimpl.NewTokenRepo(
		postgresql.DAO(cfg.User.Domain.Tables.Token),
	)

	session := sessionapp.NewSessionClearAppService(
		loginrepositoryadapter.LoginAdapter(),
		sessionrepositoryadapter.NewSessionAdapter(redisdb.DAO()),
	)

	userAppService = userapp.NewUserService(
		userrepo, member, git, t, loginrepositoryadapter.LoginAdapter(), oidcimpl.NewAuthingUser(),
		session, &cfg.User.Domain, obsadapter.NewClient(obs.Client()),
		emailadapter.NewEmailImpl(email.GetEmailInst(), cfg.User.Domain.AuditEmail, cfg.User.Domain.MailTemplates))
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "token is a admin tool for merlin server token administrator.",
	Run: func(cmd *cobra.Command, args []string) {
		Error(cmd, args, errors.New("unrecognized command"))
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
		inittoken()
	},
}

var tokenAddCmd = &cobra.Command{
	Use: "add",
	Run: func(cmd *cobra.Command, args []string) {
		acc, err := primitive.NewAccount(viper.GetString("token.create.name"))
		if err != nil {
			logrus.Fatalf("add user token failed :%s with %s", err.Error(), viper.GetString("token.create.name"))
		}
		tokenName, err := primitive.NewTokenName(viper.GetString("token.create.token_name"))
		if err != nil {
			logrus.Fatalf("add user token failed :%s with  %s", err.Error(), viper.GetString("token.create.token_name"))
		}
		tokenPerm := viper.GetString("token.create.perm")

		ctx := context.Background()
		platform, err := userAppService.GetPlatformUser(ctx, acc)
		if err != nil {
			logrus.Fatalf("failed to get platform user %s", err.Error())
		}

		perm, err := primitive.NewTokenPerm(tokenPerm)
		if err != nil {
			logrus.Fatal(err)
		}

		token, err := userAppService.CreateToken(ctx, &domain.TokenCreatedCmd{
			Account:    acc,
			Name:       tokenName,
			Permission: perm,
		}, platform)
		if err != nil {
			logrus.Fatalf("add user token failed :%s", err.Error())
		} else {
			fmt.Printf(token.Token)
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
		inittoken()
	},
}

var tokenDelCmd = &cobra.Command{
	Use: "del",
	Run: func(cmd *cobra.Command, args []string) {
		acc, err := primitive.NewAccount(viper.GetString("token.del.name"))
		if err != nil {
			logrus.Fatalf("delete user token failed :%s with %s", err.Error(), viper.GetString("token.del.name"))
		}

		tokenName, err := primitive.NewTokenName(viper.GetString("token.del.token_name"))
		if err != nil {
			logrus.Fatalf(err.Error())
		}
		ctx := context.Background()
		platform, err := userAppService.GetPlatformUser(ctx, acc)
		if err != nil {
			logrus.Fatalf("failed to get platform user , %s", err.Error())
		}

		fmt.Println("delete ", acc.Account(), tokenName)
		err = userAppService.DeleteToken(ctx, &domain.TokenDeletedCmd{
			Account: acc,
			Name:    tokenName,
		}, platform)
		if err != nil {
			logrus.Fatalf("get user token failed :%s", err.Error())
		} else {
			logrus.Infof("delete user %s token %s success", acc.Account(), tokenName)
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
		inittoken()
	},
}

var tokenListCmd = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		acc, err := primitive.NewAccount(viper.GetString("token.list.name"))
		if err != nil {
			logrus.Fatalf("get user token failed :%s with %s", err.Error(), viper.GetString("token.get.name"))
		}

		tokens, err := userAppService.ListTokens(acc)
		if err != nil {
			logrus.Fatalf("get user token failed :%s", err.Error())
		} else {
			fmt.Println("User tokens:")
			for _, token := range tokens {
				fmt.Printf("%#v\n", token)
			}
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
		inittoken()
	},
}

var tokenGetCmd = &cobra.Command{
	Use: "get",
	Run: func(cmd *cobra.Command, args []string) {
		acc, err := primitive.NewAccount(viper.GetString("token.get.name"))
		if err != nil {
			logrus.Fatalf("get user token failed :%s with %s", err.Error(), viper.GetString("token.get.name"))
		}
		name, err := primitive.NewTokenName(viper.GetString("token.get.token_name"))
		if err != nil {
			logrus.Fatalf("get user token failed :%s with %s", err.Error(), viper.GetString("token.get.name"))
		}
		token, err := userAppService.GetToken(ctx, acc, name)
		if err != nil {
			logrus.Fatalf("get user token failed :%s", err.Error())
		} else {
			fmt.Println("User tokens:")
			fmt.Printf("%#v\n", token)
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
		inittoken()
	},
}

var tokenVerifyCmd = &cobra.Command{
	Use: "verify",
	Run: func(cmd *cobra.Command, args []string) {
		token := viper.GetString("token.verify.token")

		_, err := userAppService.VerifyToken(token, primitive.NewReadPerm())
		if err != nil {
			logrus.Infof("verify user token failed, %s", err.Error())
		} else {
			logrus.Infof("verify user token success")
		}

	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
		inittoken()
	},
}

func init() {
	tokenCmd.AddCommand(tokenGetCmd)
	tokenCmd.AddCommand(tokenListCmd)
	tokenCmd.AddCommand(tokenDelCmd)
	tokenCmd.AddCommand(tokenAddCmd)
	tokenCmd.AddCommand(tokenVerifyCmd)
	// 添加命令行参数
	tokenAddCmd.Flags().StringP("name", "n", "", "user name")
	tokenAddCmd.Flags().StringP("token", "t", "", "token name")
	tokenAddCmd.Flags().StringP("perm", "p", "", "permission of then, allowed: write or read")

	if err := tokenAddCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}
	if err := tokenAddCmd.MarkFlagRequired("token"); err != nil {
		logrus.Fatal(err)
	}
	if err := tokenAddCmd.MarkFlagRequired("perm"); err != nil {
		logrus.Fatal(err)
	}

	tokenVerifyCmd.Flags().StringP("token", "t", "", "token value")
	if err := tokenVerifyCmd.MarkFlagRequired("token"); err != nil {
		logrus.Fatal(err)
	}

	if err := viper.BindPFlag("token.verify.token", tokenVerifyCmd.Flags().Lookup("token")); err != nil {
		logrus.Fatal(err)
	}

	tokenDelCmd.Flags().StringP("name", "n", "", "user name")
	tokenDelCmd.Flags().StringP("token", "t", "", "token name")
	if err := tokenDelCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}
	if err := tokenDelCmd.MarkFlagRequired("token"); err != nil {
		logrus.Fatal(err)
	}

	tokenListCmd.Flags().StringP("name", "n", "", "user name")
	if err := tokenListCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}

	tokenGetCmd.Flags().StringP("name", "n", "", "user name")
	tokenGetCmd.Flags().StringP("token", "t", "", "token name")
	if err := tokenGetCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}
	if err := tokenGetCmd.MarkFlagRequired("token"); err != nil {
		logrus.Fatal(err)
	}

	if err := viper.BindPFlag("token.get.name", tokenGetCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("token.get.token_name", tokenGetCmd.Flags().Lookup("token")); err != nil {
		logrus.Fatal(err)
	}

	if err := viper.BindPFlag("token.create.token_name", tokenAddCmd.Flags().Lookup("token")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("token.create.name", tokenAddCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("token.create.perm", tokenAddCmd.Flags().Lookup("perm")); err != nil {
		logrus.Fatal(err)
	}

	if err := viper.BindPFlag("token.del.token_name", tokenDelCmd.Flags().Lookup("token")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("token.del.name", tokenDelCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("token.list.name", tokenListCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}
}

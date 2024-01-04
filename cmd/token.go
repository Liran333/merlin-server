package main

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/common/infrastructure/gitea"
	"github.com/openmerlin/merlin-server/infrastructure/giteauser"
	"github.com/openmerlin/merlin-server/infrastructure/mongodb"
	userapp "github.com/openmerlin/merlin-server/user/app"
	"github.com/openmerlin/merlin-server/user/domain"
	usergit "github.com/openmerlin/merlin-server/user/infrastructure/git"
	userrepoimpl "github.com/openmerlin/merlin-server/user/infrastructure/repositoryimpl"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "token is a admin tool for merlin server token administrator.",
	Run: func(cmd *cobra.Command, args []string) {
		Error(cmd, args, errors.New("unrecognized command"))
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
	},
}

var tokenAddCmd = &cobra.Command{
	Use: "add",
	Run: func(cmd *cobra.Command, args []string) {
		acc, err := primitive.NewAccount(viper.GetString("token.create.name"))
		if err != nil {
			logrus.Fatalf("add user token failed :%s with %s", err.Error(), viper.GetString("token.create.name"))
		}
		tokenName := viper.GetString("token.create.token_name")
		tokenPerm := viper.GetString("token.create.perm")

		user := userrepoimpl.NewUserRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.User),
		)

		git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))
		t := userrepoimpl.NewTokenRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Token),
		)
		userAppService := userapp.NewUserService(
			user, git, t)

		platform, err := userAppService.GetPlatformUser(acc)
		if err != nil {
			logrus.Fatalf("failed to get platform user %s", err)
		}

		fmt.Println("create token", acc.Account(), tokenName)

		perm, err := primitive.NewTokenPerm(tokenPerm)
		if err != nil {
			logrus.Fatal(err)
		}

		token, err := userAppService.CreateToken(&domain.TokenCreatedCmd{
			Account:    acc,
			Name:       tokenName,
			Permission: perm,
		}, platform)
		if err != nil {
			logrus.Fatalf("add user token failed :%s", err.Error())
		} else {
			logrus.Infof("add user %s token %s success, token is %s", acc.Account(), tokenName, token.Token)
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
	},
}

var tokenDelCmd = &cobra.Command{
	Use: "del",
	Run: func(cmd *cobra.Command, args []string) {
		acc, err := primitive.NewAccount(viper.GetString("token.del.name"))
		if err != nil {
			logrus.Fatalf("delete user token failed :%s with %s", err.Error(), viper.GetString("token.del.name"))
		}

		tokenName := viper.GetString("token.del.token_name")

		user := userrepoimpl.NewUserRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.User),
		)

		git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))
		t := userrepoimpl.NewTokenRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Token),
		)
		userAppService := userapp.NewUserService(
			user, git, t)

		platform, err := userAppService.GetPlatformUser(acc)
		if err != nil {
			logrus.Fatalf("failed to get platform user , %s", err)
		}

		fmt.Println("delete ", acc.Account(), tokenName)
		err = userAppService.DeleteToken(&domain.TokenDeletedCmd{
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
	},
}

var tokenGetCmd = &cobra.Command{
	Use: "get",
	Run: func(cmd *cobra.Command, args []string) {
		acc, err := primitive.NewAccount(viper.GetString("token.get.name"))
		if err != nil {
			logrus.Fatalf("get user token failed :%s with %s", err.Error(), viper.GetString("token.get.name"))
		}

		user := userrepoimpl.NewUserRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.User),
		)

		git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))
		t := userrepoimpl.NewTokenRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Token),
		)
		userAppService := userapp.NewUserService(
			user, git, t)

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
	},
}

var tokenVerifyCmd = &cobra.Command{
	Use: "verify",
	Run: func(cmd *cobra.Command, args []string) {
		token := viper.GetString("token.verify.token")

		user := userrepoimpl.NewUserRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.User),
		)
		t := userrepoimpl.NewTokenRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Token),
		)
		git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))

		userAppService := userapp.NewUserService(
			user, git, t)

		_, b := userAppService.VerifyToken(token, primitive.NewReadPerm())
		logrus.Infof("verify user token result %t", b)

	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
	},
}

func init() {
	tokenCmd.AddCommand(tokenGetCmd)
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

	tokenGetCmd.Flags().StringP("name", "n", "", "user name")
	if err := tokenGetCmd.MarkFlagRequired("name"); err != nil {
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
	if err := viper.BindPFlag("token.get.name", tokenGetCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}
}

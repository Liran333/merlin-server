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
	"github.com/openmerlin/merlin-server/utils"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "user is a admin tool for merlin server user administrator.",
	Run: func(cmd *cobra.Command, args []string) {
		Error(cmd, args, errors.New("unrecognized command"))
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
	},
}

var userAddCmd = &cobra.Command{
	Use: "add",
	Run: func(cmd *cobra.Command, args []string) {
		acc, err := primitive.NewAccount(viper.GetString("user.create.name"))
		if err != nil {
			logrus.Fatalf("create user failed :%s", err.Error())
		}
		email, err := domain.NewEmail(viper.GetString("user.create.email"))
		if err != nil {
			logrus.Fatalf("create user failed :%s", err.Error())
		}
		bio, err := domain.NewBio("")
		if err != nil {
			logrus.Fatalf("create user failed :%s", err.Error())
		}
		ava, err := domain.NewAvatarId("")
		if err != nil {
			logrus.Fatalf("create user failed :%s", err.Error())
		}

		token := userrepoimpl.NewTokenRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Token),
		)
		user := userrepoimpl.NewUserRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.User),
		)

		git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))

		userAppService := userapp.NewUserService(
			user, git, token)

		_, err = userAppService.Create(&domain.UserCreateCmd{
			Email:    email,
			Account:  acc,
			Bio:      bio,
			AvatarId: ava,
		})
		if err != nil {
			logrus.Fatalf("create user failed :%s", err.Error())
		} else {
			logrus.Info("create user successfully")
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
	},
}

var userGetCmd = &cobra.Command{
	Use: "show",
	Run: func(cmd *cobra.Command, args []string) {
		acc, err := primitive.NewAccount(viper.GetString("user.get.name"))
		if err != nil {
			logrus.Fatalf("get user failed :%s", err.Error())
		}

		user := userrepoimpl.NewUserRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.User),
		)
		token := userrepoimpl.NewTokenRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Token),
		)
		git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))

		userAppService := userapp.NewUserService(
			user, git, token)

		u, err := userAppService.GetByAccount(acc, false)
		if err != nil {
			logrus.Fatalf("get user failed :%s", err.Error())
		} else {
			fmt.Println("User info:")
			fmt.Printf("Name: %s\n", u.Account)
			fmt.Printf("Email: %s\n", u.Email)
			fmt.Printf("Bio: %s\n", u.Bio)
			fmt.Printf("AvatarId: %s\n", u.AvatarId)
			fmt.Printf("Id: %s\n", u.Id)
			fmt.Printf("Fullname: %s\n", u.Fullname)
			fmt.Printf("Created: %s\n", utils.ToDate(u.CreatedAt))
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
	},
}

var userGetAvaCmd = &cobra.Command{
	Use: "showava",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		acc := viper.GetStringSlice("user.getava.name")
		users := make([]primitive.Account, len(acc))
		for i := range acc {
			users[i], err = primitive.NewAccount(acc[i])
			if err != nil {
				logrus.Fatalf("get user failed :%s", err.Error())
			}
		}

		user := userrepoimpl.NewUserRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.User),
		)
		token := userrepoimpl.NewTokenRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Token),
		)
		git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))

		userAppService := userapp.NewUserService(
			user, git, token)

		u, err := userAppService.GetUsersAvatarId(users)
		if err != nil {
			logrus.Fatalf("get user failed :%s", err.Error())
		} else {
			for i := range u {
				fmt.Printf("Name: %s\n", u[i].Name)
				fmt.Printf("AvatarId: %s\n", u[i].AvatarId)
			}
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
	},
}

var userDelCmd = &cobra.Command{
	Use: "del",
	Run: func(cmd *cobra.Command, args []string) {
		acc, err := primitive.NewAccount(viper.GetString("user.del.name"))
		if err != nil {
			logrus.Fatalf("delete user failed :%s with %s", err.Error(), viper.GetString("user.del.name"))
		}

		user := userrepoimpl.NewUserRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.User),
		)
		token := userrepoimpl.NewTokenRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Token),
		)
		git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))

		userAppService := userapp.NewUserService(
			user, git, token)

		err = userAppService.Delete(acc)
		if err != nil {
			logrus.Fatalf("delete user failed :%s", err.Error())
		} else {
			logrus.Info("delete user successfully")
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
	},
}

var userEditCmd = &cobra.Command{
	Use: "edit",
	Run: func(cmd *cobra.Command, args []string) {
		acc, err := primitive.NewAccount(viper.GetString("user.edit.name"))
		if err != nil {
			logrus.Fatalf("edit user failed :%s with %s", err.Error(), viper.GetString("user.edit.name"))
		}
		updateCmd := userapp.UpdateUserBasicInfoCmd{}
		avatar, err := domain.NewAvatarId(viper.GetString("user.edit.avatar"))
		if err == nil {
			updateCmd.AvatarId = avatar
		}
		bio, err := domain.NewBio(viper.GetString("user.edit.bio"))
		if err == nil {
			updateCmd.Bio = bio
		}
		email, err := domain.NewEmail(viper.GetString("user.edit.email"))
		if err == nil {
			updateCmd.Email = email
		}
		fullname := viper.GetString("user.edit.fullname")
		if fullname != "" {
			updateCmd.Fullname = fullname
		}
		user := userrepoimpl.NewUserRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.User),
		)
		token := userrepoimpl.NewTokenRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Token),
		)
		git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))

		userAppService := userapp.NewUserService(
			user, git, token)

		err = userAppService.UpdateBasicInfo(acc, updateCmd)
		if err != nil {
			logrus.Fatalf("edit user failed :%s", err.Error())
		} else {
			logrus.Info("edit user successfully")
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
	},
}

func init() {
	userCmd.AddCommand(userAddCmd)
	userCmd.AddCommand(userDelCmd)
	userCmd.AddCommand(userGetCmd)
	userCmd.AddCommand(userGetAvaCmd)
	userCmd.AddCommand(userEditCmd)
	// 添加命令行参数
	userAddCmd.Flags().StringP("name", "n", "", "user name")
	userAddCmd.Flags().StringP("email", "e", "", "user email")
	if err := userAddCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}
	if err := userAddCmd.MarkFlagRequired("email"); err != nil {
		logrus.Fatal(err)
	}

	userDelCmd.Flags().StringP("name", "n", "", "user name")
	if err := userDelCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}

	userGetCmd.Flags().StringP("name", "n", "", "user name")
	if err := userGetCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}

	userGetAvaCmd.Flags().StringSlice("name", make([]string, 0), "user name")
	if err := userGetAvaCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("user.getava.name", userGetAvaCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}

	userEditCmd.Flags().StringP("name", "n", "", "user name")
	userEditCmd.Flags().StringP("email", "e", "", "user email")
	userEditCmd.Flags().StringP("bio", "b", "", "user bio")
	userEditCmd.Flags().StringP("avatar", "a", "", "user avatar")
	userEditCmd.Flags().StringP("fullname", "f", "", "user fullname")
	if err := userEditCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}

	userEditCmd.MarkFlagsOneRequired("avatar", "bio", "email", "fullname")

	if err := viper.BindPFlag("user.create.name", userAddCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("user.create.email", userAddCmd.Flags().Lookup("email")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("user.del.name", userDelCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("user.get.name", userGetCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("user.edit.name", userEditCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("user.edit.email", userEditCmd.Flags().Lookup("email")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("user.edit.bio", userEditCmd.Flags().Lookup("bio")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("user.edit.fullname", userEditCmd.Flags().Lookup("fullname")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("user.edit.avatar", userEditCmd.Flags().Lookup("avatar")); err != nil {
		logrus.Fatal(err)
	}
}

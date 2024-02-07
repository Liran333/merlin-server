package main

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	userapp "github.com/openmerlin/merlin-server/user/app"
	"github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/user/domain/repository"
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
		inittoken()
	},
}

var userAddCmd = &cobra.Command{
	Use: "add",
	Run: func(cmd *cobra.Command, args []string) {
		acc, err := primitive.NewAccount(viper.GetString("user.create.name"))
		if err != nil {
			logrus.Fatalf("create user failed :%s", err.Error())
		}
		email, err := primitive.NewEmail(viper.GetString("user.create.email"))
		if err != nil {
			logrus.Fatalf("create user failed :%s", err.Error())
		}
		phone, err := primitive.NewPhone(viper.GetString("user.create.phone"))
		if err != nil {
			logrus.Fatalf("create user failed :%s", err.Error())
		}
		fullname, err := primitive.NewMSDFullname(viper.GetString("user.create.fullname"))
		if err != nil {
			logrus.Fatalf("create user failed :%s", err.Error())
		}
		desc, err := primitive.NewMSDDesc("")
		if err != nil {
			logrus.Fatalf("create user failed :%s", err.Error())
		}
		ava, err := primitive.NewAvatarId("")
		if err != nil {
			logrus.Fatalf("create user failed :%s", err.Error())
		}

		_, err = userAppService.Create(&domain.UserCreateCmd{
			Email:    email,
			Account:  acc,
			Desc:     desc,
			AvatarId: ava,
			Fullname: fullname,
			Phone:    phone,
		})
		if err != nil {
			logrus.Fatalf("create user failed :%s", err.Error())
		} else {
			logrus.Info("create user successfully")
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
		inittoken()
	},
}

var userGetCmd = &cobra.Command{
	Use: "show",
	Run: func(cmd *cobra.Command, args []string) {
		acc, err := primitive.NewAccount(viper.GetString("user.get.name"))
		if err != nil {
			logrus.Fatalf("get user failed :%s", err.Error())
		}
		ac, err := primitive.NewAccount(actor)
		if err != nil {
			logrus.Fatalf("invalid owner :%s", err.Error())
		}
		u, err := userAppService.GetByAccount(ac, acc)
		if err != nil {
			logrus.Fatalf("get user failed :%s", err.Error())
		} else {
			fmt.Println("User info:")
			fmt.Printf("Name: %s\n", u.Name)
			fmt.Printf("Email: %s\n", *u.Email)
			fmt.Printf("Bio: %s\n", u.Description)
			fmt.Printf("AvatarId: %s\n", u.AvatarId)
			fmt.Printf("Id: %s\n", u.Id)
			fmt.Printf("Fullname: %s\n", u.Fullname)
			fmt.Printf("Created: %s\n", utils.ToDate(u.CreatedAt))
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
		inittoken()
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
		inittoken()
	},
}

var userDelCmd = &cobra.Command{
	Use: "del",
	Run: func(cmd *cobra.Command, args []string) {
		acc, err := primitive.NewAccount(viper.GetString("user.del.name"))
		if err != nil {
			logrus.Fatalf("delete user failed :%s with %s", err.Error(), viper.GetString("user.del.name"))
		}

		err = userAppService.Delete(acc)
		if err != nil {
			logrus.Fatalf("delete user failed :%s", err.Error())
		} else {
			logrus.Info("delete user successfully")
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
		inittoken()
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
		avatar, err := primitive.NewAvatarId(viper.GetString("user.edit.avatar"))
		if err == nil {
			updateCmd.AvatarId = avatar
		}
		desc, err := primitive.NewMSDDesc(viper.GetString("user.edit.bio"))
		if err == nil {
			updateCmd.Desc = desc
		}
		email, err := primitive.NewEmail(viper.GetString("user.edit.email"))
		if err == nil {
			updateCmd.Email = email
		}
		fullname, err := primitive.NewMSDFullname(viper.GetString("user.edit.fullname"))
		if err != nil {
			updateCmd.Fullname = fullname
		}

		_, err = userAppService.UpdateBasicInfo(acc, updateCmd)
		if err != nil {
			logrus.Fatalf("edit user failed :%s", err.Error())
		} else {
			logrus.Info("edit user successfully")
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
		inittoken()
	},
}

var userListCmd = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		t := domain.UserTypeUser
		us, _, err := userrepo.ListAccount(&repository.ListOption{Type: &t})
		if err != nil {
			logrus.Fatalf("list users failed :%s", err.Error())
		}

		for _, u := range us {
			logrus.Infof("user %s", u.Account.Account())
			logrus.Infof("email %s", u.Email.Email())
			logrus.Infof("phone %s", u.Phone.PhoneNumber())
			logrus.Infof("pwd %s", u.PlatformPwd)
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
		inittoken()
	},
}

func init() {
	userCmd.AddCommand(userAddCmd)
	userCmd.AddCommand(userDelCmd)
	userCmd.AddCommand(userGetCmd)
	userCmd.AddCommand(userGetAvaCmd)
	userCmd.AddCommand(userEditCmd)
	userCmd.AddCommand(userListCmd)
	// 添加命令行参数
	userAddCmd.Flags().StringP("name", "n", "", "user name")
	userAddCmd.Flags().StringP("email", "e", "", "user email")
	userAddCmd.Flags().StringP("fullname", "f", "", "user fullname")
	userAddCmd.Flags().StringP("phone", "p", "", "user phone number")
	if err := userAddCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}
	if err := userAddCmd.MarkFlagRequired("phone"); err != nil {
		logrus.Fatal(err)
	}

	if err := viper.BindPFlag("user.create.name", userAddCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("user.create.email", userAddCmd.Flags().Lookup("email")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("user.create.phone", userAddCmd.Flags().Lookup("phone")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("user.create.fullname", userAddCmd.Flags().Lookup("fullname")); err != nil {
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

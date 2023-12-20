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
	orgapp "github.com/openmerlin/merlin-server/organization/app"
	"github.com/openmerlin/merlin-server/organization/domain"
	orgrepoimpl "github.com/openmerlin/merlin-server/organization/infrastructure/repositoryimpl"
	userapp "github.com/openmerlin/merlin-server/user/app"
	usergit "github.com/openmerlin/merlin-server/user/infrastructure/git"
	userrepoimpl "github.com/openmerlin/merlin-server/user/infrastructure/repositoryimpl"
)

var orgCmd = &cobra.Command{
	Use:   "org",
	Short: "org is a admin tool for merlin server organization administrator.",
	Run: func(cmd *cobra.Command, args []string) {
		Error(cmd, args, errors.New("unrecognized command"))
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
	},
}

var orgAddCmd = &cobra.Command{
	Use: "add",
	Run: func(cmd *cobra.Command, args []string) {
		orgName := viper.GetString("org.create.name")
		owner := viper.GetString("org.create.owner")
		fullname := viper.GetString("org.create.fullname")
		website := viper.GetString("org.create.website")
		ava := viper.GetString("org.create.avatarid")
		desc := viper.GetString("org.create.description")

		org := orgrepoimpl.NewOrgRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Organization),
		)

		member := orgrepoimpl.NewMemberRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Member),
		)

		user := userrepoimpl.NewUserRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.User),
		)
		t := userrepoimpl.NewTokenRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Token),
		)
		git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))

		userAppService := userapp.NewUserService(
			user, git, t)

		orgAppService := orgapp.NewOrgService(
			userAppService, org, member, 1209600)

		_, err := orgAppService.Create(&domain.OrgCreatedCmd{
			Name:        orgName,
			AvatarId:    ava,
			FullName:    fullname,
			Website:     website,
			Owner:       owner,
			Description: desc,
		})
		if err != nil {
			logrus.Fatalf("create org failed :%s", err.Error())
		} else {
			logrus.Info("create org successfully")
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
	},
}

var memberAddCmd = &cobra.Command{
	Use: "addmem",
	Run: func(cmd *cobra.Command, args []string) {
		orgName, err := primitive.NewAccount(viper.GetString("member.add.org"))
		if err != nil {
			logrus.Fatalf("invalid org name :%s", err.Error())
		}
		userName, err := primitive.NewAccount(viper.GetString("member.add.user"))
		if err != nil {
			logrus.Fatalf("invalid user name :%s", err.Error())
		}

		org := orgrepoimpl.NewOrgRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Organization),
		)

		member := orgrepoimpl.NewMemberRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Member),
		)
		t := userrepoimpl.NewTokenRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Token),
		)
		user := userrepoimpl.NewUserRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.User),
		)

		git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))

		userAppService := userapp.NewUserService(
			user, git, t)

		orgAppService := orgapp.NewOrgService(
			userAppService, org, member, 1209600)

		err = orgAppService.AddMember(&orgapp.OrgAddMemberCmd{
			Account: userName,
			Org:     orgName,
		})
		if err != nil {
			logrus.Fatalf("add member failed :%s", err.Error())
		} else {
			logrus.Info("add member successfully")
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
	},
}

var memberListCmd = &cobra.Command{
	Use: "listmem",
	Run: func(cmd *cobra.Command, args []string) {
		orgName, err := primitive.NewAccount(viper.GetString("member.list.org"))
		if err != nil {
			logrus.Fatalf("invalid org name :%s", err.Error())
		}

		org := orgrepoimpl.NewOrgRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Organization),
		)

		member := orgrepoimpl.NewMemberRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Member),
		)

		user := userrepoimpl.NewUserRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.User),
		)
		t := userrepoimpl.NewTokenRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Token),
		)
		git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))

		userAppService := userapp.NewUserService(
			user, git, t)

		orgAppService := orgapp.NewOrgService(
			userAppService, org, member, 1209600)

		members, err := orgAppService.ListMember(orgName)
		if err != nil {
			logrus.Fatalf("list member failed :%s", err.Error())
		} else {
			logrus.Info("Member Info:")
			for _, m := range members {
				logrus.Infof("Member org: %s\n", m.OrgName)
				logrus.Infof("Member org full name: %s\n", m.OrgFullName)
				logrus.Infof("Member user: %s\n", m.UserName)
				logrus.Infof("Member role: %s\n", m.Role)
			}
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
	},
}

var inviteSendCmd = &cobra.Command{
	Use: "sendinv",
	Run: func(cmd *cobra.Command, args []string) {
		orgName, err := primitive.NewAccount(viper.GetString("invite.add.org"))
		if err != nil {
			logrus.Fatalf("invalid org name :%s", err.Error())
		}
		userName, err := primitive.NewAccount(viper.GetString("invite.add.user"))
		if err != nil {
			logrus.Fatalf("invalid user name :%s", err.Error())
		}
		role := viper.GetString("invite.add.role")

		org := orgrepoimpl.NewOrgRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Organization),
		)

		member := orgrepoimpl.NewMemberRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Member),
		)

		user := userrepoimpl.NewUserRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.User),
		)

		git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))
		t := userrepoimpl.NewTokenRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Token),
		)
		userAppService := userapp.NewUserService(
			user, git, t)

		orgAppService := orgapp.NewOrgService(
			userAppService, org, member, 1209600)

		_, err = orgAppService.InviteMember(&orgapp.OrgInviteMemberCmd{
			Account: userName,
			Role:    role,
			Org:     orgName,
		})
		if err != nil {
			logrus.Fatalf("add invite failed :%s", err.Error())
		} else {
			logrus.Info("add invite successfully")
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
	},
}

var inviteListCmd = &cobra.Command{
	Use: "listinv",
	Run: func(cmd *cobra.Command, args []string) {
		orgName, err := primitive.NewAccount(viper.GetString("invite.list.org"))
		if err != nil {
			logrus.Fatalf("invalid org name :%s", err.Error())
		}

		org := orgrepoimpl.NewOrgRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Organization),
		)

		member := orgrepoimpl.NewMemberRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Member),
		)

		user := userrepoimpl.NewUserRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.User),
		)
		t := userrepoimpl.NewTokenRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Token),
		)
		git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))

		userAppService := userapp.NewUserService(
			user, git, t)

		orgAppService := orgapp.NewOrgService(
			userAppService, org, member, 1209600)

		members, err := orgAppService.ListInvitation(orgName)
		if err != nil {
			logrus.Fatalf("list member failed :%s", err.Error())
		} else {
			logrus.Info("Member Info:")
			for _, m := range members {
				logrus.Infof("Member org: %s\n", m.OrgName)
				logrus.Infof("Member user: %s\n", m.UserName)
				logrus.Infof("Member role: %s\n", m.Role)
				logrus.Infof("Member expire: %d\n", m.ExpiresAt)
			}
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
	},
}

var removeInviteCmd = &cobra.Command{
	Use: "rminv",
	Run: func(cmd *cobra.Command, args []string) {
		orgName, err := primitive.NewAccount(viper.GetString("invite.del.org"))
		if err != nil {
			logrus.Fatalf("invalid org name :%s", err.Error())
		}
		userName, err := primitive.NewAccount(viper.GetString("invite.del.user"))
		if err != nil {
			logrus.Fatalf("invalid user name :%s", err.Error())
		}

		org := orgrepoimpl.NewOrgRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Organization),
		)

		member := orgrepoimpl.NewMemberRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Member),
		)

		user := userrepoimpl.NewUserRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.User),
		)

		git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))
		t := userrepoimpl.NewTokenRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Token),
		)
		userAppService := userapp.NewUserService(
			user, git, t)

		orgAppService := orgapp.NewOrgService(
			userAppService, org, member, 1209600)

		_, err = orgAppService.RevokeInvite(&orgapp.OrgRemoveInviteCmd{
			Org:     orgName,
			Account: userName,
		})
		if err != nil {
			logrus.Fatalf("revoke invite failed :%s", err.Error())
		} else {
			logrus.Fatalf("revoke invite successfully")
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
	},
}

var memberEditCmd = &cobra.Command{
	Use: "editmem",
	Run: func(cmd *cobra.Command, args []string) {
		orgName, err := primitive.NewAccount(viper.GetString("member.edit.org"))
		if err != nil {
			logrus.Fatalf("invalid org name :%s", err.Error())
		}
		userName, err := primitive.NewAccount(viper.GetString("member.edit.user"))
		if err != nil {
			logrus.Fatalf("invalid user name :%s", err.Error())
		}
		role := viper.GetString("member.edit.role")

		org := orgrepoimpl.NewOrgRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Organization),
		)

		member := orgrepoimpl.NewMemberRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Member),
		)

		user := userrepoimpl.NewUserRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.User),
		)

		git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))
		t := userrepoimpl.NewTokenRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Token),
		)
		userAppService := userapp.NewUserService(
			user, git, t)

		orgAppService := orgapp.NewOrgService(
			userAppService, org, member, 1209600)

		_, err = orgAppService.EditMember(&orgapp.OrgEditMemberCmd{
			Account: userName,
			Org:     orgName,
			Role:    role,
		})
		if err != nil {
			logrus.Fatalf("edit member failed :%s", err.Error())
		} else {
			logrus.Info("edit member successfully")
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
	},
}

var memberRemoveCmd = &cobra.Command{
	Use: "rmmem",
	Run: func(cmd *cobra.Command, args []string) {
		orgName, err := primitive.NewAccount(viper.GetString("member.del.org"))
		if err != nil {
			logrus.Fatalf("invalid org name :%s", err.Error())
		}
		userName, err := primitive.NewAccount(viper.GetString("member.del.user"))
		if err != nil {
			logrus.Fatalf("invalid user name :%s", err.Error())
		}

		org := orgrepoimpl.NewOrgRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Organization),
		)

		member := orgrepoimpl.NewMemberRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Member),
		)

		user := userrepoimpl.NewUserRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.User),
		)

		git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))
		t := userrepoimpl.NewTokenRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Token),
		)
		userAppService := userapp.NewUserService(
			user, git, t)

		orgAppService := orgapp.NewOrgService(
			userAppService, org, member, 1209600)

		err = orgAppService.RemoveMember(&orgapp.OrgRemoveMemberCmd{
			Account: userName,
			Org:     orgName,
		})
		if err != nil {
			logrus.Fatalf("remove member failed :%s", err.Error())
		} else {
			logrus.Info("remove member successfully")
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
	},
}

var orgGetCmd = &cobra.Command{
	Use: "show",
	Run: func(cmd *cobra.Command, args []string) {
		acc, _ := primitive.NewAccount(viper.GetString("org.get.name"))
		owner, _ := primitive.NewAccount(viper.GetString("org.get.owner"))
		u, _ := primitive.NewAccount(viper.GetString("org.get.user"))

		org := orgrepoimpl.NewOrgRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Organization),
		)
		user := userrepoimpl.NewUserRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.User),
		)
		member := orgrepoimpl.NewMemberRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Member),
		)

		git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))
		t := userrepoimpl.NewTokenRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Token),
		)
		userAppService := userapp.NewUserService(
			user, git, t)

		orgAppService := orgapp.NewOrgService(
			userAppService, org, member, 1209600)

		if acc != nil {
			u, err := orgAppService.GetByAccount(acc)
			if err != nil {
				logrus.Fatalf("get org failed :%s", err.Error())
			} else {
				fmt.Println("Org info:")
				fmt.Printf("Name: %s\n", u.Name)
				fmt.Printf("Full name: %s\n", u.FullName)
				fmt.Printf("Website: %s\n", u.Website)
				fmt.Printf("AvatarId: %s\n", u.AvatarId)
				fmt.Printf("Id: %s\n", u.Id)
				fmt.Printf("Description: %s\n", u.Description)
				fmt.Printf("Owner: %s\n", u.Owner)
			}
		} else if owner != nil {
			orgs, err := orgAppService.GetByOwner(owner)
			if err != nil {
				logrus.Fatalf("get org by owner failed :%s", err.Error())
			} else {
				for o := range orgs {
					fmt.Println("Org info:")
					fmt.Printf("Name: %s\n", orgs[o].Name)
					fmt.Printf("Full name: %s\n", orgs[o].FullName)
					fmt.Printf("Website: %s\n", orgs[o].Website)
					fmt.Printf("AvatarId: %s\n", orgs[o].AvatarId)
					fmt.Printf("Id: %s\n", orgs[o].Id)
					fmt.Printf("Description: %s\n", orgs[o].Description)
					fmt.Printf("Owner: %s\n", orgs[o].Owner)
				}
			}
		} else if u != nil {
			orgs, err := orgAppService.GetByUser(u)
			if err != nil {
				logrus.Fatalf("get org by user failed :%s", err.Error())
			} else {
				for o := range orgs {
					fmt.Println("Org info:")
					fmt.Printf("Name: %s\n", orgs[o].Name)
					fmt.Printf("Full name: %s\n", orgs[o].FullName)
					fmt.Printf("Website: %s\n", orgs[o].Website)
					fmt.Printf("AvatarId: %s\n", orgs[o].AvatarId)
					fmt.Printf("Id: %s\n", orgs[o].Id)
					fmt.Printf("Description: %s\n", orgs[o].Description)
					fmt.Printf("Owner: %s\n", orgs[o].Owner)
				}
			}
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
	},
}

var orgDelCmd = &cobra.Command{
	Use: "del",
	Run: func(cmd *cobra.Command, args []string) {
		acc, err := primitive.NewAccount(viper.GetString("org.del.name"))
		if err != nil {
			logrus.Fatalf("del org failed :%s", err.Error())
		}

		org := orgrepoimpl.NewOrgRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Organization),
		)
		user := userrepoimpl.NewUserRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.User),
		)

		git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))
		t := userrepoimpl.NewTokenRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Token),
		)
		userAppService := userapp.NewUserService(
			user, git, t)
		member := orgrepoimpl.NewMemberRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Member),
		)
		orgAppService := orgapp.NewOrgService(
			userAppService, org, member, 1209600)

		err = orgAppService.Delete(acc)
		if err != nil {
			logrus.Fatalf("delete org failed :%s", err.Error())
		} else {
			logrus.Info("delete org successfully")
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
	},
}

var orgEditCmd = &cobra.Command{
	Use: "edit",
	Run: func(cmd *cobra.Command, args []string) {
		acc, err := primitive.NewAccount(viper.GetString("org.edit.name"))
		if err != nil {
			logrus.Fatalf("edit org failed :%s with %s", err.Error(), viper.GetString("org.edit.name"))
		}
		updateCmd := orgapp.UpdateOrgBasicInfoCmd{}
		avatar := viper.GetString("org.edit.avatar")
		if avatar != "" {
			updateCmd.AvatarId = avatar
		}
		website := viper.GetString("org.edit.website")
		if website != "" {
			updateCmd.Website = website
		}
		fullname := viper.GetString("org.edit.fullname")
		if fullname != "" {
			logrus.Infof("change full name to %s", fullname)
			updateCmd.FullName = fullname
		}
		desc := viper.GetString("org.edit.desc")
		if desc != "" {
			updateCmd.Description = desc
		}

		org := orgrepoimpl.NewOrgRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Organization),
		)
		user := userrepoimpl.NewUserRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.User),
		)
		t := userrepoimpl.NewTokenRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Token),
		)
		git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))

		userAppService := userapp.NewUserService(
			user, git, t)
		member := orgrepoimpl.NewMemberRepo(
			mongodb.NewCollection(cfg.Mongodb.Collections.Member),
		)
		orgAppService := orgapp.NewOrgService(
			userAppService, org, member, 1209600)

		_, err = orgAppService.UpdateBasicInfo(acc, &updateCmd)
		if err != nil {
			logrus.Fatalf("edit org failed :%s", err.Error())
		} else {
			logrus.Info("edit org successfully")
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
	},
}

func init() {
	orgCmd.AddCommand(orgAddCmd)
	orgCmd.AddCommand(orgDelCmd)
	orgCmd.AddCommand(orgGetCmd)
	orgCmd.AddCommand(orgEditCmd)
	orgCmd.AddCommand(memberAddCmd)
	orgCmd.AddCommand(memberListCmd)
	orgCmd.AddCommand(memberRemoveCmd)
	orgCmd.AddCommand(memberEditCmd)
	orgCmd.AddCommand(removeInviteCmd)
	orgCmd.AddCommand(inviteListCmd)
	orgCmd.AddCommand(inviteSendCmd)
	// 添加命令行参数

	orgAddCmd.Flags().StringP("name", "n", "", "org name")
	orgAddCmd.Flags().StringP("owner", "o", "", "org owner")
	orgAddCmd.Flags().StringP("website", "w", "", "org website")
	orgAddCmd.Flags().StringP("fullname", "f", "", "org fullname")
	orgAddCmd.Flags().StringP("avatar", "a", "", "org avatar")
	orgAddCmd.Flags().StringP("desc", "d", "", "org description")
	if err := orgAddCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}
	if err := orgAddCmd.MarkFlagRequired("owner"); err != nil {
		logrus.Fatal(err)
	}
	orgAddCmd.MarkFlagsOneRequired("avatar", "fullname", "website", "desc")
	if err := viper.BindPFlag("org.create.website", orgAddCmd.Flags().Lookup("website")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("org.create.fullname", orgAddCmd.Flags().Lookup("fullname")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("org.create.avatar", orgAddCmd.Flags().Lookup("avatar")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("org.create.desc", orgAddCmd.Flags().Lookup("desc")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("org.create.name", orgAddCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("org.create.owner", orgAddCmd.Flags().Lookup("owner")); err != nil {
		logrus.Fatal(err)
	}

	orgDelCmd.Flags().StringP("name", "n", "", "org name")
	if err := orgDelCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("org.del.name", orgDelCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}

	orgGetCmd.Flags().StringP("name", "n", "", "org name")
	orgGetCmd.Flags().StringP("owner", "o", "", "org owner")
	orgGetCmd.Flags().StringP("user", "u", "", "org member")

	orgGetCmd.MarkFlagsOneRequired("name", "owner", "user")
	if err := viper.BindPFlag("org.get.name", orgGetCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("org.get.owner", orgGetCmd.Flags().Lookup("owner")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("org.get.user", orgGetCmd.Flags().Lookup("user")); err != nil {
		logrus.Fatal(err)
	}

	inviteSendCmd.Flags().StringP("name", "n", "", "org name")
	inviteSendCmd.Flags().StringP("user", "u", "", "member name")
	inviteSendCmd.Flags().StringP("role", "r", "", "member role")
	if err := inviteSendCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}
	if err := inviteSendCmd.MarkFlagRequired("user"); err != nil {
		logrus.Fatal(err)
	}
	if err := inviteSendCmd.MarkFlagRequired("role"); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("invite.add.org", inviteSendCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("invite.add.user", inviteSendCmd.Flags().Lookup("user")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("invite.add.role", inviteSendCmd.Flags().Lookup("role")); err != nil {
		logrus.Fatal(err)
	}

	inviteListCmd.Flags().StringP("name", "n", "", "org name")
	if err := inviteListCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("invite.list.org", inviteListCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}

	removeInviteCmd.Flags().StringP("name", "n", "", "org name")
	removeInviteCmd.Flags().StringP("user", "u", "", "member name")
	removeInviteCmd.Flags().StringP("role", "r", "", "member role")
	if err := removeInviteCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}
	if err := removeInviteCmd.MarkFlagRequired("user"); err != nil {
		logrus.Fatal(err)
	}
	if err := removeInviteCmd.MarkFlagRequired("role"); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("invite.del.org", removeInviteCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("invite.del.user", removeInviteCmd.Flags().Lookup("user")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("invite.del.role", removeInviteCmd.Flags().Lookup("role")); err != nil {
		logrus.Fatal(err)
	}

	orgEditCmd.Flags().StringP("name", "n", "", "org name")
	orgEditCmd.Flags().StringP("website", "w", "", "org website")
	orgEditCmd.Flags().StringP("fullname", "f", "", "org fullname")
	orgEditCmd.Flags().StringP("avatar", "a", "", "org avatar")
	orgEditCmd.Flags().StringP("desc", "d", "", "org description")

	if err := orgEditCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}

	orgEditCmd.MarkFlagsOneRequired("avatar", "fullname", "website", "desc")
	if err := viper.BindPFlag("org.edit.name", orgEditCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("org.edit.website", orgEditCmd.Flags().Lookup("website")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("org.edit.fullname", orgEditCmd.Flags().Lookup("fullname")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("org.edit.avatar", orgEditCmd.Flags().Lookup("avatar")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("org.edit.desc", orgEditCmd.Flags().Lookup("desc")); err != nil {
		logrus.Fatal(err)
	}

	memberAddCmd.Flags().StringP("name", "n", "", "org name")
	memberAddCmd.Flags().StringP("user", "u", "", "member name")
	if err := memberAddCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}
	if err := memberAddCmd.MarkFlagRequired("user"); err != nil {
		logrus.Fatal(err)
	}

	if err := viper.BindPFlag("member.add.org", memberAddCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("member.add.user", memberAddCmd.Flags().Lookup("user")); err != nil {
		logrus.Fatal(err)
	}

	memberRemoveCmd.Flags().StringP("name", "n", "", "org name")
	memberRemoveCmd.Flags().StringP("user", "u", "", "member name")

	if err := memberRemoveCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}
	if err := memberRemoveCmd.MarkFlagRequired("user"); err != nil {
		logrus.Fatal(err)
	}

	if err := viper.BindPFlag("member.del.org", memberRemoveCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("member.del.user", memberRemoveCmd.Flags().Lookup("user")); err != nil {
		logrus.Fatal(err)
	}

	memberListCmd.Flags().StringP("name", "n", "", "org name")
	if err := memberListCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("member.list.org", memberListCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}

	memberEditCmd.Flags().StringP("name", "n", "", "org name")
	memberEditCmd.Flags().StringP("user", "u", "", "member name")
	memberEditCmd.Flags().StringP("role", "r", "", "member role")
	if err := memberEditCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}
	if err := memberEditCmd.MarkFlagRequired("user"); err != nil {
		logrus.Fatal(err)
	}
	if err := memberEditCmd.MarkFlagRequired("role"); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("member.edit.org", memberEditCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("member.edit.user", memberEditCmd.Flags().Lookup("user")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("member.edit.role", memberEditCmd.Flags().Lookup("role")); err != nil {
		logrus.Fatal(err)
	}

}

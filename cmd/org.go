/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package main is the entry point for the application.
package main

import (
	"errors"
	"fmt"

	redisdb "github.com/opensourceways/redis-lib"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/openmerlin/merlin-server/common/domain/crypto"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/common/infrastructure/gitea"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/infrastructure/giteauser"
	orgapp "github.com/openmerlin/merlin-server/organization/app"
	"github.com/openmerlin/merlin-server/organization/domain"
	"github.com/openmerlin/merlin-server/organization/infrastructure/messageadapter"
	orgrepoimpl "github.com/openmerlin/merlin-server/organization/infrastructure/repositoryimpl"
	sessionapp "github.com/openmerlin/merlin-server/session/app"
	"github.com/openmerlin/merlin-server/session/infrastructure/loginrepositoryadapter"
	"github.com/openmerlin/merlin-server/session/infrastructure/oidcimpl"
	"github.com/openmerlin/merlin-server/session/infrastructure/sessionrepositoryadapter"
	userapp "github.com/openmerlin/merlin-server/user/app"
	usergit "github.com/openmerlin/merlin-server/user/infrastructure/git"
	userrepoimpl "github.com/openmerlin/merlin-server/user/infrastructure/repositoryimpl"
	"github.com/openmerlin/merlin-server/utils"
)

var orgAppService orgapp.OrgService

func initorg() {
	invite := orgrepoimpl.NewInviteRepo(
		postgresql.DAO(cfg.Org.Domain.Tables.Invite),
	)

	member := orgrepoimpl.NewMemberRepo(
		postgresql.DAO(cfg.Org.Domain.Tables.Member),
	)
	p := orgapp.NewPermService(&cfg.Permission, member)

	user := userrepoimpl.NewUserRepo(
		postgresql.DAO(cfg.User.Domain.Tables.User), crypto.NewEncryption(cfg.User.Domain.Key),
	)
	t := userrepoimpl.NewTokenRepo(
		postgresql.DAO(cfg.User.Domain.Tables.Token),
	)
	git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))

	session := sessionapp.NewSessionClearAppService(
		loginrepositoryadapter.LoginAdapter(),
		sessionrepositoryadapter.NewSessionAdapter(redisdb.DAO()),
	)

	certRepo, _ := orgrepoimpl.NewCertificateImpl(
		postgresql.DAO(cfg.Org.Domain.Tables.Certificate),
		crypto.NewEncryption(cfg.User.Domain.Key),
	)

	userService := userapp.NewUserService(user, member, git, t, loginrepositoryadapter.LoginAdapter(),
		oidcimpl.NewAuthingUser(), session, &cfg.User.Domain)

	orgAppService = orgapp.NewOrgService(userService, user, member, invite, p, &cfg.Org.Domain, git,
		messageadapter.MessageAdapter(&cfg.Org.Domain.Topics), certRepo,
	)
}

var orgCmd = &cobra.Command{
	Use:   "org",
	Short: "org is a admin tool for merlin server organization administrator.",
	Run: func(cmd *cobra.Command, args []string) {
		Error(cmd, args, errors.New("unrecognized command"))
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
		initorg()
	},
}

var orgAddCmd = &cobra.Command{
	Use: "add",
	Run: func(cmd *cobra.Command, args []string) {
		orgName, err := primitive.NewAccount(viper.GetString("org.create.name"))
		if err != nil {
			logrus.Fatalf("invalid org name :%s", err.Error())
		}
		fullname, err := primitive.NewOrgFullname(viper.GetString("org.create.fullname"))
		if err != nil {
			logrus.Fatalf("invalid fullname :%s", err.Error())
		}
		var website primitive.Website
		website, err = primitive.NewOrgWebsite(viper.GetString("org.create.website"))
		if err != nil {
			logrus.Fatalf("invalid website :%s", err.Error())
		}

		ava, err := primitive.NewAvatarId(viper.GetString("org.create.avatarid"))
		if err != nil {
			logrus.Fatalf("invalid avatarid :%s", err.Error())
		}
		desc, err := primitive.NewAccountDesc(viper.GetString("org.create.description"))
		if err != nil {
			logrus.Fatalf("invalid description :%s", err.Error())
		}
		owner, err := primitive.NewAccount(actor)
		if err != nil {
			logrus.Fatalf("invalid owner :%s", err.Error())
		}

		_, err = orgAppService.Create(&domain.OrgCreatedCmd{
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
		initorg()

	},
}

var memberAddCmd = &cobra.Command{
	Use: "approve",
	Run: func(cmd *cobra.Command, args []string) {
		orgName, err := primitive.NewAccount(viper.GetString("member.approve.org"))
		if err != nil {
			logrus.Fatalf("invalid org name :%s", err.Error())
		}
		userName, err := primitive.NewAccount(viper.GetString("member.approve.user"))
		if err != nil {
			logrus.Fatalf("invalid user name :%s", err.Error())
		}

		req := &domain.OrgApproveRequestMemberCmd{}
		req.Actor = primitive.CreateAccount(actor)
		req.Org = orgName
		req.Requester = userName
		_, err = orgAppService.ApproveRequest(req)
		if err != nil {
			logrus.Fatalf("add member failed :%s", err.Error())
		} else {
			logrus.Info("add member successfully")
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
		initorg()

	},
}

var memberAcceptCmd = &cobra.Command{
	Use: "accept",
	Run: func(cmd *cobra.Command, args []string) {
		orgName, err := primitive.NewAccount(viper.GetString("member.accept.org"))
		if err != nil {
			logrus.Fatalf("invalid org name :%s", err.Error())
		}
		msg := viper.GetString("member.accept.msg")

		req := &domain.OrgAcceptInviteCmd{}
		req.Actor = primitive.CreateAccount(actor)
		req.Org = orgName
		req.Msg = msg
		req.Account = primitive.CreateAccount(actor)
		_, err = orgAppService.AcceptInvite(req)
		if err != nil {
			logrus.Fatalf("accept approve failed :%s", err.Error())
		} else {
			logrus.Info("accept approve successfully")
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
		initorg()

	},
}

var orgCheckCmd = &cobra.Command{
	Use: "check",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := primitive.NewAccount(viper.GetString("org.check.name"))
		if err != nil {
			logrus.Fatalf("invalid org name :%s", err.Error())
		}

		ok := orgAppService.CheckName(name)
		if !ok {
			logrus.Fatalf("check name failed")
		} else {
			logrus.Info("check name successfully")
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
		initorg()

	},
}

var memberListCmd = &cobra.Command{
	Use: "listmem",
	Run: func(cmd *cobra.Command, args []string) {
		orgName, err := primitive.NewAccount(viper.GetString("member.list.org"))
		if err != nil {
			logrus.Fatalf("invalid org name :%s", err.Error())
		}

		members, err := orgAppService.ListMember(&domain.OrgListMemberCmd{Org: orgName})
		if err != nil {
			logrus.Fatalf("list member failed :%s", err.Error())
		} else {
			fmt.Print("Member Info:")
			for _, m := range members {
				fmt.Printf("Member org: %s\n", m.OrgName)
				fmt.Printf("Member org full name: %s\n", m.OrgFullName)
				fmt.Printf("Member user: %s\n", m.UserName)
				fmt.Printf("Member role: %s\n", m.Role)
			}
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
		initorg()

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
		role, err := primitive.NewRole(viper.GetString("invite.add.role"))
		if err != nil {
			logrus.Fatalf("invalid user role :%s", err.Error())
		}

		_, err = orgAppService.InviteMember(&domain.OrgInviteMemberCmd{
			Actor:   primitive.CreateAccount(actor),
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
		initorg()
	},
}

var requestCmd = &cobra.Command{
	Use: "request",
	Run: func(cmd *cobra.Command, args []string) {
		orgName, err := primitive.NewAccount(viper.GetString("req.add.org"))
		if err != nil {
			logrus.Fatalf("invalid org name :%s", err.Error())
		}

		req := &domain.OrgRequestMemberCmd{}
		req.Actor = primitive.CreateAccount(actor)
		req.Org = orgName
		req.Msg = viper.GetString("req.add.msg")
		_, err = orgAppService.RequestMember(req)
		if err != nil {
			logrus.Fatalf("add member request failed :%s", err.Error())
		} else {
			logrus.Info("add member request successfully")
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
		initorg()
	},
}

var reqListCmd = &cobra.Command{
	Use: "listreq",
	Run: func(cmd *cobra.Command, args []string) {
		orgName, _ := primitive.NewAccount(viper.GetString("req.list.org"))

		user, _ := primitive.NewAccount(viper.GetString("req.list.requester"))

		req := &domain.OrgMemberReqListCmd{}
		req.Actor = primitive.CreateAccount(actor)
		req.Org = orgName
		req.Requester = user
		req.Status = domain.ApproveStatus(viper.GetString("req.list.status"))

		reqs, err := orgAppService.ListMemberReq(req)
		if err != nil {
			logrus.Fatalf("list requests failed :%s", err.Error())
		} else {
			for _, o := range reqs {
				fmt.Printf("\nOrg: %s\n", o.OrgName)
				fmt.Printf("Role: %s\n", o.Role)
				fmt.Printf("Fullname: %s\n", o.Fullname)
				fmt.Printf("Requester: %s\n", o.Username)
				fmt.Printf("Status: %s\n", o.Status)
				fmt.Printf("Msg: %s\n", o.Msg)
				fmt.Printf("By: %s\n", o.By)
				_, create := utils.DateAndTime(o.CreatedAt)
				fmt.Printf("Created at: %s\n", create)
				_, update := utils.DateAndTime(o.UpdatedAt)
				fmt.Printf("Updated at: %s\n", update)
			}
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
		initorg()

	},
}

var inviteListCmd = &cobra.Command{
	Use: "listinv",
	Run: func(cmd *cobra.Command, args []string) {
		orgName, _ := primitive.NewAccount(viper.GetString("invite.list.org"))

		inviter, _ := primitive.NewAccount(viper.GetString("invite.list.inviter"))
		invitee, _ := primitive.NewAccount(viper.GetString("invite.list.invitee"))

		var err error
		var invites []orgapp.ApproveDTO

		req := &domain.OrgInvitationListCmd{}
		req.Invitee = invitee
		req.Org = orgName
		req.Actor = primitive.CreateAccount(actor)
		req.Status = domain.ApproveStatus(viper.GetString("invite.list.status"))
		req.Inviter = inviter

		invites, err = orgAppService.ListInvitationByOrg(req.Actor, orgName, req.Status)

		if err != nil {
			logrus.Fatalf("list invites failed :%s", err.Error())
		} else {
			for _, o := range invites {
				fmt.Printf("\nOrg: %s\n", o.OrgName)
				fmt.Printf("Role: %s\n", o.Role)
				fmt.Printf("Fullname: %s\n", o.Fullname)
				fmt.Printf("Inviter: %s\n", o.Inviter)
				fmt.Printf("Invitee: %s\n", o.UserName)
				fmt.Printf("Status: %s\n", o.Status)
				_, exp := utils.DateAndTime(o.ExpiresAt)
				fmt.Printf("Expire: %s\n", exp)
				fmt.Printf("Msg: %s\n", o.Msg)
				fmt.Printf("By: %s\n", o.By)
				_, create := utils.DateAndTime(o.CreatedAt)
				fmt.Printf("Created at: %s\n", create)
				_, update := utils.DateAndTime(o.UpdatedAt)
				fmt.Printf("Updated at: %s\n", update)
			}
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
		initorg()

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
			userName = primitive.CreateAccount(actor)
		}

		_, err = orgAppService.RevokeInvite(&domain.OrgRemoveInviteCmd{
			Actor:   primitive.CreateAccount(actor),
			Org:     orgName,
			Account: userName,
			Msg:     viper.GetString("invite.del.msg"),
		})
		if err != nil {
			logrus.Fatalf("revoke invite failed :%s", err.Error())
		} else {
			logrus.Fatalf("revoke invite successfully")
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
		initorg()

	},
}

var removeReqCmd = &cobra.Command{
	Use: "rmreq",
	Run: func(cmd *cobra.Command, args []string) {
		orgName, err := primitive.NewAccount(viper.GetString("req.del.org"))
		if err != nil {
			logrus.Fatalf("invalid org name :%s", err.Error())
		}
		userName, err := primitive.NewAccount(viper.GetString("req.del.requester"))
		if err != nil {
			userName = primitive.CreateAccount(actor)
		}

		req := domain.OrgCancelRequestMemberCmd{
			Actor:     primitive.CreateAccount(actor),
			Org:       orgName,
			Requester: userName,
			Msg:       viper.GetString("req.del.msg"),
		}

		_, err = orgAppService.CancelReqMember(&req)
		if err != nil {
			logrus.Fatalf("revoke member request failed :%s", err.Error())
		} else {
			logrus.Fatalf("revoke member request successfully")
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
		initorg()

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
		role, err := primitive.NewRole(viper.GetString("member.edit.role"))
		if err != nil {
			logrus.Fatalf("invalid user role :%s", err.Error())
		}

		_, err = orgAppService.EditMember(&domain.OrgEditMemberCmd{
			Actor:   primitive.CreateAccount(actor),
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
		initorg()

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

		err = orgAppService.RemoveMember(&domain.OrgRemoveMemberCmd{
			Actor:   primitive.CreateAccount(actor),
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
		initorg()

	},
}

var orgGetCmd = &cobra.Command{
	Use: "show",
	Run: func(cmd *cobra.Command, args []string) {
		acc, _ := primitive.NewAccount(viper.GetString("org.get.name"))
		owner, _ := primitive.NewAccount(viper.GetString("org.get.owner"))
		u, _ := primitive.NewAccount(viper.GetString("org.get.user"))

		if acc != nil {
			u, err := orgAppService.GetByAccount(acc)
			if err != nil {
				logrus.Fatalf("get org failed :%s", err.Error())
			} else {
				fmt.Println("Org info:")
				fmt.Printf("Name: %s\n", u.Name)
				fmt.Printf("Full name: %s\n", u.Fullname)
				fmt.Printf("Website: %s\n", *u.Website)
				fmt.Printf("AvatarId: %s\n", u.AvatarId)
				fmt.Printf("Id: %s\n", u.Id)
				fmt.Printf("Description: %s\n", u.Description)
				fmt.Printf("Owner: %s\n", *u.Owner)
				fmt.Printf("Default role: %s\n", u.DefaultRole)
				fmt.Printf("Allow request: %t\n", *u.AllowRequest)
			}
		} else if owner != nil {
			orgs, err := orgAppService.GetByOwner(primitive.CreateAccount(actor), owner)
			if err != nil {
				logrus.Fatalf("get org by owner failed :%s", err.Error())
			} else {
				for o := range orgs {
					fmt.Println("Org info:")
					fmt.Printf("Name: %s\n", orgs[o].Name)
					fmt.Printf("Full name: %s\n", orgs[o].Fullname)
					fmt.Printf("Website: %s\n", *orgs[o].Website)
					fmt.Printf("AvatarId: %s\n", orgs[o].AvatarId)
					fmt.Printf("Id: %s\n", orgs[o].Id)
					fmt.Printf("Description: %s\n", orgs[o].Description)
					fmt.Printf("Owner: %s\n", *orgs[o].Owner)
				}
			}
		} else if u != nil {
			orgs, err := orgAppService.GetByUser(primitive.CreateAccount(actor), u)
			if err != nil {
				logrus.Fatalf("get org by user failed :%s", err.Error())
			} else {
				for o := range orgs {
					fmt.Println("Org info:")
					fmt.Printf("Name: %s\n", orgs[o].Name)
					fmt.Printf("Full name: %s\n", orgs[o].Fullname)
					fmt.Printf("Website: %s\n", *orgs[o].Website)
					fmt.Printf("AvatarId: %s\n", orgs[o].AvatarId)
					fmt.Printf("Id: %s\n", orgs[o].Id)
					fmt.Printf("Description: %s\n", orgs[o].Description)
					fmt.Printf("Owner: %s\n", *orgs[o].Owner)
				}
			}
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
		initorg()

	},
}

var orgDelCmd = &cobra.Command{
	Use: "del",
	Run: func(cmd *cobra.Command, args []string) {
		acc, err := primitive.NewAccount(viper.GetString("org.del.name"))
		if err != nil {
			logrus.Fatalf("del org failed :%s", err.Error())
		}

		err = orgAppService.Delete(&domain.OrgDeletedCmd{
			Actor: primitive.CreateAccount(actor),
			Name:  acc,
		})
		if err != nil {
			logrus.Fatalf("delete org failed :%s", err.Error())
		} else {
			logrus.Info("delete org successfully")
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
		initorg()

	},
}

var orgEditCmd = &cobra.Command{
	Use: "edit",
	Run: func(cmd *cobra.Command, args []string) {
		acc, err := primitive.NewAccount(viper.GetString("org.edit.name"))
		if err != nil {
			logrus.Fatalf("edit org failed :%s with %s", err.Error(), viper.GetString("org.edit.name"))
		}
		updateCmd := domain.OrgUpdatedBasicInfoCmd{}
		avatar := viper.GetString("org.edit.avatar")
		if avatar != "" {
			if updateCmd.AvatarId, err = primitive.NewAvatarId(avatar); err != nil {
				logrus.Fatalf("edit org failed :%s with %s", err.Error(), viper.GetString("org.edit.avatar"))
			}
		}
		website := viper.GetString("org.edit.website")
		if website != "" {
			if updateCmd.Website, err = primitive.NewOrgWebsite(website); err != nil {
				logrus.Fatalf("edit org failed :%s with %s", err.Error(), viper.GetString("org.edit.website"))
			}
		}
		fullname := viper.GetString("org.edit.fullname")
		if fullname != "" {
			if updateCmd.FullName, err = primitive.NewOrgFullname(fullname); err != nil {
				logrus.Fatalf("edit org failed :%s with %s", err.Error(), viper.GetString("org.edit.fullname"))
			}
			fmt.Printf("change full name to %s", fullname)
		}
		desc := viper.GetString("org.edit.desc")
		if desc != "" {
			if updateCmd.Description, err = primitive.NewAccountDesc(desc); err != nil {
				logrus.Fatalf("edit org failed :%s with %s", err.Error(), viper.GetString("org.edit.desc"))
			}
		}
		*updateCmd.AllowRequest = viper.GetBool("org.edit.allowrequest")
		updateCmd.DefaultRole, err = primitive.NewRole(viper.GetString("org.edit.defaultrole"))
		if err != nil {
			logrus.Fatalf("edit org failed :%s", err.Error())
		}

		updateCmd.Actor = primitive.CreateAccount(actor)
		updateCmd.OrgName = acc
		_, err = orgAppService.UpdateBasicInfo(&updateCmd)
		if err != nil {
			logrus.Fatalf("edit org failed :%s", err.Error())
		} else {
			logrus.Info("edit org successfully")
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		initServer(configFile)
		initorg()

	},
}

func init() {
	orgCmd.AddCommand(orgAddCmd)
	orgCmd.AddCommand(orgDelCmd)
	orgCmd.AddCommand(orgGetCmd)
	orgCmd.AddCommand(orgEditCmd)
	orgCmd.AddCommand(reqListCmd)
	orgCmd.AddCommand(memberAddCmd)
	orgCmd.AddCommand(memberListCmd)
	orgCmd.AddCommand(memberRemoveCmd)
	orgCmd.AddCommand(memberEditCmd)
	orgCmd.AddCommand(removeInviteCmd)
	orgCmd.AddCommand(removeReqCmd)
	orgCmd.AddCommand(inviteListCmd)
	orgCmd.AddCommand(inviteSendCmd)
	orgCmd.AddCommand(orgCheckCmd)
	orgCmd.AddCommand(memberAcceptCmd)
	orgCmd.AddCommand(requestCmd)

	orgAddCmd.Flags().StringP("name", "n", "", "org name")
	orgAddCmd.Flags().StringP("website", "w", "", "org website")
	orgAddCmd.Flags().StringP("fullname", "f", "", "org fullname")
	orgAddCmd.Flags().StringP("avatar", "a", "", "org avatar")
	orgAddCmd.Flags().StringP("desc", "d", "", "org description")
	if err := orgAddCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}
	if err := orgAddCmd.MarkFlagRequired("fullname"); err != nil {
		logrus.Fatal(err)
	}

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

	orgDelCmd.Flags().StringP("name", "n", "", "org name")
	if err := orgDelCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("org.del.name", orgDelCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}

	orgCheckCmd.Flags().StringP("name", "n", "", "org name")
	if err := orgCheckCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("org.check.name", orgCheckCmd.Flags().Lookup("name")); err != nil {
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

	requestCmd.Flags().StringP("name", "n", "", "org name")
	requestCmd.Flags().StringP("msg", "m", "", "req msg")
	if err := requestCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}

	if err := viper.BindPFlag("req.add.org", requestCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("req.add.msg", requestCmd.Flags().Lookup("msg")); err != nil {
		logrus.Fatal(err)
	}

	removeReqCmd.Flags().StringP("name", "n", "", "org name")
	removeReqCmd.Flags().StringP("requester", "r", "", "requester")
	removeReqCmd.Flags().StringP("msg", "m", "", "msg")

	if err := removeReqCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}

	if err := viper.BindPFlag("req.del.org", removeReqCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("req.del.requester", removeReqCmd.Flags().Lookup("requester")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("req.del.msg", removeReqCmd.Flags().Lookup("msg")); err != nil {
		logrus.Fatal(err)
	}

	inviteListCmd.Flags().StringP("name", "n", "", "org name")
	inviteListCmd.Flags().StringP("inviter", "r", "", "inviter name")
	inviteListCmd.Flags().StringP("invitee", "e", "", "invitee name")
	inviteListCmd.Flags().StringP("status", "s", "", "invitee name")

	inviteListCmd.MarkFlagsOneRequired("name", "invitee", "inviter")
	if err := viper.BindPFlag("invite.list.org", inviteListCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("invite.list.invitee", inviteListCmd.Flags().Lookup("invitee")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("invite.list.inviter", inviteListCmd.Flags().Lookup("inviter")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("invite.list.status", inviteListCmd.Flags().Lookup("status")); err != nil {
		logrus.Fatal(err)
	}

	reqListCmd.Flags().StringP("name", "n", "", "org name")
	reqListCmd.Flags().StringP("requester", "r", "", "requester name")
	reqListCmd.Flags().StringP("status", "s", "", "request status")
	reqListCmd.MarkFlagsOneRequired("name", "requester")
	if err := viper.BindPFlag("req.list.org", reqListCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("req.list.requester", reqListCmd.Flags().Lookup("requester")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("req.list.status", reqListCmd.Flags().Lookup("status")); err != nil {
		logrus.Fatal(err)
	}

	removeInviteCmd.Flags().StringP("name", "n", "", "org name")
	removeInviteCmd.Flags().StringP("user", "u", "", "member name")
	removeInviteCmd.Flags().StringP("msg", "m", "", "msg")

	if err := removeInviteCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}

	if err := viper.BindPFlag("invite.del.org", removeInviteCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("invite.del.user", removeInviteCmd.Flags().Lookup("user")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("invite.del.msg", removeInviteCmd.Flags().Lookup("msg")); err != nil {
		logrus.Fatal(err)
	}
	orgEditCmd.Flags().StringP("name", "n", "", "org name")
	orgEditCmd.Flags().StringP("website", "w", "", "org website")
	orgEditCmd.Flags().StringP("fullname", "f", "", "org fullname")
	orgEditCmd.Flags().StringP("avatar", "a", "", "org avatar")
	orgEditCmd.Flags().StringP("desc", "d", "", "org description")
	orgEditCmd.Flags().Bool("allowrequest", false, "whether allow request member")
	orgEditCmd.Flags().StringP("default_role", "r", "", "default role when request member")

	if err := orgEditCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}

	orgEditCmd.MarkFlagsOneRequired("avatar", "fullname", "website", "desc", "allowrequest", "default_role")
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
	if err := viper.BindPFlag("org.edit.allowrequest", orgEditCmd.Flags().Lookup("allowrequest")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("org.edit.defaultrole", orgEditCmd.Flags().Lookup("default_role")); err != nil {
		logrus.Fatal(err)
	}

	memberAddCmd.Flags().StringP("name", "n", "", "org name")
	memberAddCmd.Flags().StringP("user", "u", "", "member name")
	memberAddCmd.Flags().StringP("msg", "m", "", "msg")
	if err := memberAddCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}
	if err := memberAddCmd.MarkFlagRequired("user"); err != nil {
		logrus.Fatal(err)
	}

	if err := viper.BindPFlag("member.approve.org", memberAddCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("member.approve.user", memberAddCmd.Flags().Lookup("user")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("member.approve.msg", memberAddCmd.Flags().Lookup("msg")); err != nil {
		logrus.Fatal(err)
	}
	memberAcceptCmd.Flags().StringP("name", "n", "", "org name")
	memberAcceptCmd.Flags().StringP("msg", "m", "", "msg")
	if err := memberAcceptCmd.MarkFlagRequired("name"); err != nil {
		logrus.Fatal(err)
	}

	if err := viper.BindPFlag("member.accept.org", memberAcceptCmd.Flags().Lookup("name")); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.BindPFlag("member.accept.msg", memberAcceptCmd.Flags().Lookup("msg")); err != nil {
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

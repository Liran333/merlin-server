/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package server provides functionality for setting up and configuring a server for handling code repo operations.
package server

import (
	"github.com/gin-gonic/gin"

	commonapp "github.com/openmerlin/merlin-server/common/app"
	"github.com/openmerlin/merlin-server/common/domain/crypto"
	"github.com/openmerlin/merlin-server/common/infrastructure/email"
	"github.com/openmerlin/merlin-server/common/infrastructure/gitea"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/config"
	"github.com/openmerlin/merlin-server/infrastructure/giteauser"
	"github.com/openmerlin/merlin-server/organization/app"
	"github.com/openmerlin/merlin-server/organization/controller"
	"github.com/openmerlin/merlin-server/organization/infrastructure/emailimpl"
	"github.com/openmerlin/merlin-server/organization/infrastructure/messageadapter"
	orgrepoimpl "github.com/openmerlin/merlin-server/organization/infrastructure/repositoryimpl"
	usergit "github.com/openmerlin/merlin-server/user/infrastructure/git"
	userrepoimpl "github.com/openmerlin/merlin-server/user/infrastructure/repositoryimpl"
)

// initOrg depends on initUser
func initOrg(cfg *config.Config, services *allServices) error {
	org := userrepoimpl.NewUserRepo(postgresql.DAO(cfg.User.Domain.Tables.User), crypto.NewEncryption(cfg.User.Domain.Key))

	orgMember := orgrepoimpl.NewMemberRepo(postgresql.DAO(cfg.Org.Domain.Tables.Member))

	invitation := orgrepoimpl.NewInviteRepo(postgresql.DAO(cfg.Org.Domain.Tables.Invite))

	permission := app.NewPermService(&cfg.Permission, orgMember)

	git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))

	certRepo, err := orgrepoimpl.NewCertificateImpl(
		postgresql.DAO(cfg.Org.Domain.Tables.Certificate),
		crypto.NewEncryption(cfg.User.Domain.Key),
	)
	if err != nil {
		return err
	}

	services.orgCertificateApp = app.NewOrgCertificateService(
		permission, emailimpl.NewEmailImpl(email.GetEmailInst(), cfg.Org.Domain.CertificateEmail), certRepo,
	)

	services.orgApp = app.NewOrgService(
		services.userApp, org, orgMember,
		invitation, permission, &cfg.Org.Domain, git,
		messageadapter.MessageAdapter(&cfg.Org.Domain.Topics), certRepo,
	)

	services.npuGatekeeper = app.NewPrivilegeOrgService(services.orgApp, cfg.PrivilegeOrg.Npu, app.AllocNpu)
	services.disable = app.NewPrivilegeOrgService(services.orgApp, cfg.PrivilegeOrg.Disable, app.Disable)
	services.permissionApp = commonapp.NewResourcePermissionAppService(permission, services.disable)

	return nil
}

func setRouterOfOrg(v1 *gin.RouterGroup, cfg *config.Config, services *allServices) {
	controller.AddRouterForOrgController(
		v1,
		services.orgApp,
		services.orgCertificateApp,
		services.userApp,
		services.operationLog,
		services.securityLog,
		services.userMiddleWare,
		services.rateLimiterMiddleWare,
		services.privacyCheck,
		services.npuGatekeeper,
		services.disable,
	)
}

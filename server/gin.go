package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/community-robot-lib/interrupts"
	redisdb "github.com/opensourceways/redis-lib"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/openmerlin/merlin-server/common/controller/middleware"
	"github.com/openmerlin/merlin-server/common/infrastructure/gitea"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"

	"github.com/openmerlin/merlin-server/api"
	"github.com/openmerlin/merlin-server/config"
	"github.com/openmerlin/merlin-server/controller"
	"github.com/openmerlin/merlin-server/infrastructure/giteauser"
	"github.com/openmerlin/merlin-server/infrastructure/mongodb"

	userapp "github.com/openmerlin/merlin-server/user/app"
	userctl "github.com/openmerlin/merlin-server/user/controller"
	userrepo "github.com/openmerlin/merlin-server/user/domain/repository"
	usergit "github.com/openmerlin/merlin-server/user/infrastructure/git"
	userrepoimpl "github.com/openmerlin/merlin-server/user/infrastructure/repositoryimpl"

	orgapp "github.com/openmerlin/merlin-server/organization/app"
	orgrepo "github.com/openmerlin/merlin-server/organization/domain/repository"
	orgrepoimpl "github.com/openmerlin/merlin-server/organization/infrastructure/repositoryimpl"

	modelapp "github.com/openmerlin/merlin-server/models/app"
	modelctl "github.com/openmerlin/merlin-server/models/controller"
	modelrepo "github.com/openmerlin/merlin-server/models/domain/repository"
	"github.com/openmerlin/merlin-server/models/infrastructure/modelrepositoryadapter"

	spaceapp "github.com/openmerlin/merlin-server/space/app"
	spacectl "github.com/openmerlin/merlin-server/space/controller"
	spacerepo "github.com/openmerlin/merlin-server/space/domain/repository"
	"github.com/openmerlin/merlin-server/space/infrastructure/spacerepositoryadapter"

	sessionapp "github.com/openmerlin/merlin-server/session/app"
	sessionctl "github.com/openmerlin/merlin-server/session/controller"
	"github.com/openmerlin/merlin-server/session/infrastructure/csrftokenrepositoryadapter"
	"github.com/openmerlin/merlin-server/session/infrastructure/loginrepositoryadapter"
	"github.com/openmerlin/merlin-server/session/infrastructure/oidcimpl"

	coderepoapp "github.com/openmerlin/merlin-server/coderepo/app"
	coderepoctl "github.com/openmerlin/merlin-server/coderepo/controller"
	"github.com/openmerlin/merlin-server/coderepo/infrastructure/coderepoadapter"
	"github.com/openmerlin/merlin-server/coderepo/infrastructure/coderepofileadapter"
)

func StartWebServer(port int, timeout time.Duration, cfg *config.Config) {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(logRequest())
	engine.TrustedPlatform = "x-real-ip"

	middleware.Init()

	services, err := initServices(cfg)
	if err != nil {
		logrus.Error(err)

		return
	}

	// web api
	setRouterOfWeb("/web", engine, cfg, &services)

	// restful api
	setRouterOfRestful("/api", engine, cfg, &services)

	engine.UseRawPath = true

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           engine,
		ReadHeaderTimeout: time.Duration(cfg.ReadHeaderTimeout) * time.Second,
	}

	defer interrupts.WaitForGracefulShutdown()

	interrupts.ListenAndServe(srv, timeout)
}

type allServices struct {
	userApp          userapp.UserService
	userRepo         userrepo.User
	orgMember        orgrepo.OrgMember
	sessionApp       sessionapp.SessionAppService
	permission       orgapp.Permission
	codeRepoApp      coderepoapp.CodeRepoAppService
	userMiddleWare   middleware.UserMiddleWare
	codeRepoFileApp  coderepoapp.CodeRepoFileAppService
	modelRepoAdapter modelrepo.ModelRepositoryAdapter
	spaceRepoAdapter spacerepo.SpaceRepositoryAdapter
}

func initServices(cfg *config.Config) (services allServices, err error) {
	err = modelrepositoryadapter.Init(postgresql.DB(), &cfg.Model.Tables)
	if err != nil {
		return
	}

	err = spacerepositoryadapter.Init(postgresql.DB(), &cfg.Space.Tables)
	if err != nil {
		return
	}

	collections := &cfg.Mongodb.Collections

	git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))

	token := userrepoimpl.NewTokenRepo(mongodb.NewCollection(collections.Token))

	userRepo := userrepoimpl.NewUserRepo(mongodb.NewCollection(collections.User))

	services.userApp = userapp.NewUserService(userRepo, git, token)

	services.userRepo = userRepo

	services.orgMember = orgrepoimpl.NewMemberRepo(mongodb.NewCollection(collections.Member))

	services.sessionApp = sessionAppService(cfg, services.userApp)

	services.permission = orgapp.NewPermService(&cfg.Permission, services.orgMember)

	services.codeRepoApp = coderepoapp.NewCodeRepoAppService(
		coderepoadapter.NewRepoAdapter(gitea.Client()),
	)

	services.codeRepoFileApp = codeRepoFileAppService(cfg)

	services.modelRepoAdapter = modelrepositoryadapter.ModelAdapter()

	services.spaceRepoAdapter = spacerepositoryadapter.SpaceAdapter()

	return
}

func codeRepoFileAppService(cfg *config.Config) coderepoapp.CodeRepoFileAppService {
	return coderepoapp.NewCodeRepoFileAppService(
		coderepofileadapter.NewCodeRepoFileAdapter(gitea.Client()))
}

func sessionAppService(cfg *config.Config, userApp userapp.UserService) sessionapp.SessionAppService {
	return sessionapp.NewSessionAppService(
		oidcimpl.NewAuthingUser(),
		userApp,
		cfg.Session.Domain.MaxSessionNum,
		loginrepositoryadapter.LoginAdapter(),
		csrftokenrepositoryadapter.NewCSRFTokenAdapter(redisdb.DAO()),
	)
}

// setRouter init router
func setRouterOfWeb(prefix string, engine *gin.Engine, cfg *config.Config, services *allServices) {
	api.SwaggerInfo.BasePath = prefix
	api.SwaggerInfo.Title = "merlin"
	api.SwaggerInfo.Description = "set header: 'PRIVATE-TOKEN=xxx'"

	rg := engine.Group(api.SwaggerInfo.BasePath)

	services.userMiddleWare = sessionctl.WebAPIMiddleware(services.sessionApp)

	// set routers
	setRouterOfSession(rg, services)

	setRouterOfModelWeb(rg, services)

	setRouterOfSpaceWeb(rg, services)

	setRouterOfUserAndOrg(rg, cfg, services)

	setRouteOfCodeRepoFile(rg, services)

	rg.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

// setRouter init router
func setRouterOfRestful(prefix string, engine *gin.Engine, cfg *config.Config, services *allServices) {
	api.SwaggerInfo.BasePath = prefix
	api.SwaggerInfo.Title = "merlin"
	api.SwaggerInfo.Description = "set header: 'PRIVATE-TOKEN=xxx'"

	rg := engine.Group(api.SwaggerInfo.BasePath)

	services.userMiddleWare = userctl.RestfulAPI(services.userApp)

	// set routers
	setRouterOfModelRestful(rg, services)

	setRouterOfSpaceRestful(rg, services)

	setRouterOfUserAndOrg(rg, cfg, services)

	setRouteOfCodeRepoFile(rg, services)

	rg.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

// user and org router
func setRouterOfUserAndOrg(v1 *gin.RouterGroup, cfg *config.Config, services *allServices) {
	org := orgrepoimpl.NewOrgRepo(
		mongodb.NewCollection(cfg.Mongodb.Collections.User),
	)

	invitation := orgrepoimpl.NewInviteRepo(
		mongodb.NewCollection(cfg.Mongodb.Collections.Invitation),
	)

	request := orgrepoimpl.NewMemberRequestRepo(
		mongodb.NewCollection(cfg.Mongodb.Collections.MemberRequest),
	)

	/*
		session := session.NewSessionAppService(nil)

	*/

	orgAppService := orgapp.NewOrgService(
		services.userApp, org, services.orgMember,
		invitation, request, services.permission, &cfg.Org,
	)

	controller.AddRouterForUserController(
		v1, services.userApp, services.userRepo,
	)

	controller.AddRouterForOrgController(
		v1, orgAppService, services.userApp,
	)
}

// session router
func setRouterOfSession(rg *gin.RouterGroup, services *allServices) {
	sessionctl.AddRouterForSessionController(
		rg, services.sessionApp, services.userMiddleWare,
	)
}

// code repo router
func setRouteOfCodeRepoFile(rg *gin.RouterGroup, services *allServices) {
	coderepoctl.AddRouterForCodeRepoFileController(
		rg,
		services.codeRepoFileApp,
		services.userMiddleWare,
	)
}

// model router
func setRouterOfModelRestful(rg *gin.RouterGroup, services *allServices) {
	modelctl.AddRouteForModelRestfulController(
		rg,
		modelapp.NewModelAppService(
			services.permission,
			services.codeRepoApp,
			services.modelRepoAdapter,
		),
		services.userMiddleWare,
	)
}

func setRouterOfSpaceRestful(rg *gin.RouterGroup, services *allServices) {
	spacectl.AddRouteForSpaceRestfulController(
		rg,
		spaceapp.NewSpaceAppService(
			services.permission,
			services.codeRepoApp,
			services.spaceRepoAdapter,
		),
		services.userMiddleWare,
	)
}

func setRouterOfModelWeb(rg *gin.RouterGroup, services *allServices) {
	modelctl.AddRouteForModelWebController(
		rg,
		modelapp.NewModelAppService(
			services.permission,
			services.codeRepoApp,
			services.modelRepoAdapter,
		),
		services.userMiddleWare,
		services.userApp,
	)
}

func setRouterOfSpaceWeb(rg *gin.RouterGroup, services *allServices) {
	spacectl.AddRouteForSpaceWebController(
		rg,
		spaceapp.NewSpaceAppService(
			services.permission,
			services.codeRepoApp,
			services.spaceRepoAdapter,
		),
		services.userMiddleWare,
		services.userApp,
	)
}

func logRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()

		l := controller.GetOperateLog(c)
		logrus.Infof(
			"| %d | %d | %s | %s | %s",
			c.Writer.Status(),
			endTime.Sub(startTime),
			c.Request.Method,
			c.Request.RequestURI,
			l,
		)
	}
}

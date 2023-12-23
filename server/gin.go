package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/community-robot-lib/interrupts"
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
	"github.com/openmerlin/merlin-server/login/infrastructure/oidcimpl"
	session "github.com/openmerlin/merlin-server/session/app"
	sessionrepo "github.com/openmerlin/merlin-server/session/infrastructure"

	userapp "github.com/openmerlin/merlin-server/user/app"
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
	permission       orgapp.Permission
	codeRepoApp      coderepoapp.CodeRepoAppService
	userMiddleWare   middleware.UserMiddleWare
	modelRepoAdapter modelrepo.ModelRepositoryAdapter
	codeRepoFileApp coderepoapp.CodeRepoFileAppService
}

func initServices(cfg *config.Config) (services allServices, err error) {
	err = modelrepositoryadapter.Init(postgresql.DB(), &cfg.Model.Tables)
	if err != nil {
		return
	}

	services.codeRepoApp = coderepoapp.NewCodeRepoAppService(
		coderepoadapter.NewRepoAdapter(gitea.Client()),
	)

	services.modelRepoAdapter = modelrepositoryadapter.ModelAdapter()

	collections := &cfg.Mongodb.Collections

	services.orgMember = orgrepoimpl.NewMemberRepo(
		mongodb.NewCollection(collections.Member),
	)

	services.permission = orgapp.NewPermService(&cfg.Permission, services.orgMember)

	services.userRepo = userrepoimpl.NewUserRepo(mongodb.NewCollection(collections.User))

	git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))

	token := userrepoimpl.NewTokenRepo(mongodb.NewCollection(collections.Token))

	services.userApp = userapp.NewUserService(services.userRepo, git, token)

	services.codeRepoFileApp = codeRepoFileAppService(cfg)
	return
}

// setRouter init router
func setRouterOfWeb(prefix string, engine *gin.Engine, cfg *config.Config, services *allServices) {
	api.SwaggerInfo.BasePath = prefix
	api.SwaggerInfo.Title = "merlin"
	api.SwaggerInfo.Description = "set header: 'PRIVATE-TOKEN=xxx'"

	rg := engine.Group(api.SwaggerInfo.BasePath)

	services.userMiddleWare = middleware.WebAPI()

	// set routers
	setRouterOfModelWeb(rg, services)

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

	services.userMiddleWare = middleware.RestfulAPI()

	// set routers
	setRouterOfModelRestful(rg, services)

	setRouterOfUserAndOrg(rg, cfg, services)

	rg.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func setRouterOfUserAndOrg(v1 *gin.RouterGroup, cfg *config.Config, services *allServices) {
	collections := &cfg.Mongodb.Collections

	org := orgrepoimpl.NewOrgRepo(
		mongodb.NewCollection(cfg.Mongodb.Collections.Organization),
	)
	sessrepo := sessionrepo.NewSessionRepository(
		sessionrepo.NewSessionStore(
			mongodb.NewCollection(collections.Session),
		),
	)

	session := session.NewSessionService(sessrepo)

	authingUser := oidcimpl.NewAuthingUser()

	var orgAppService orgapp.OrgService

	orgAppService = orgapp.NewOrgService(
		services.userApp, org, services.orgMember,
		services.permission, cfg.Org.InviteExpiry,
	)

	controller.AddRouterForUserController(
		v1, services.userApp, services.userRepo, authingUser,
	)

	controller.AddRouterForOrgController(
		v1, orgAppService, services.userApp,
	)

	controller.AddRouterForLoginController(
		v1, services.userApp, authingUser, session,
	)
}

func codeRepoFileAppService(cfg *config.Config) coderepoapp.CodeRepoFileAppService {
	return coderepoapp.NewCodeRepoFileAppService(
		coderepofileadapter.NewCodeRepoFileAdapter(gitea.Client()))
}

func setRouterOfModelRestful(rg *gin.RouterGroup, services *allServices) {
	modelctl.AddRouteForModelRestfulController(
		rg,
		modelapp.NewModelAppService(
			services.codeRepoApp, services.modelRepoAdapter,
			services.permission,
		),
		services.userMiddleWare,
	)
}

func setRouteOfCodeRepoFile(rg *gin.RouterGroup, services *allServices) {
	coderepoctl.AddRouterForCodeRepoFileController(
		rg,
		services.codeRepoFileApp,
		services.userMiddleWare,
	)
}

func setRouterOfModelWeb(rg *gin.RouterGroup, services *allServices) {
	modelctl.AddRouteForModelWebController(
		rg,
		modelapp.NewModelAppService(
			services.codeRepoApp, services.modelRepoAdapter,
			services.permission,
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

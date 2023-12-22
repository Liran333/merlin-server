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
	usergit "github.com/openmerlin/merlin-server/user/infrastructure/git"
	userrepoimpl "github.com/openmerlin/merlin-server/user/infrastructure/repositoryimpl"

	orgapp "github.com/openmerlin/merlin-server/organization/app"
	orgrepoimpl "github.com/openmerlin/merlin-server/organization/infrastructure/repositoryimpl"

	modelapp "github.com/openmerlin/merlin-server/models/app"
	modelctl "github.com/openmerlin/merlin-server/models/controller"
	modelrepo "github.com/openmerlin/merlin-server/models/domain/repository"
	"github.com/openmerlin/merlin-server/models/infrastructure/modelrepositoryadapter"

	coderepoapp "github.com/openmerlin/merlin-server/coderepo/app"
	"github.com/openmerlin/merlin-server/coderepo/infrastructure/coderepoadapter"
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
	services.userMiddleWare = middleware.WebAPI()

	setRouter("/web", engine, cfg, &services)

	// restfull api
	services.userMiddleWare = middleware.RestfullAPI()

	setRouter("/api", engine, cfg, &services)

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
	modelRepoAdapter modelrepo.ModelRepositoryAdapter

	codeRepoApp coderepoapp.CodeRepoAppService

	userMiddleWare middleware.UserMiddleWare
}

func initServices(cfg *config.Config) (services allServices, err error) {
	err = modelrepositoryadapter.Init(postgresql.DB(), &cfg.Model.Tables)
	if err != nil {
		return
	}

	services.modelRepoAdapter = modelrepositoryadapter.ModelAdapter()
	services.codeRepoApp = codeRepoAppService(cfg)

	return
}

// setRouter init router
func setRouter(prefix string, engine *gin.Engine, cfg *config.Config, services *allServices) {
	api.SwaggerInfo.BasePath = prefix
	api.SwaggerInfo.Title = "merlin"
	api.SwaggerInfo.Description = "set header: 'PRIVATE-TOKEN=xxx'"

	rg := engine.Group(api.SwaggerInfo.BasePath)

	// set routers
	setRouterOfModel(rg, services)

	setRouterOfUserAndOrg(rg, cfg)

	rg.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func setRouterOfUserAndOrg(v1 *gin.RouterGroup, cfg *config.Config) {
	collections := &cfg.Mongodb.Collections

	user := userrepoimpl.NewUserRepo(
		mongodb.NewCollection(collections.User),
	)

	token := userrepoimpl.NewTokenRepo(
		mongodb.NewCollection(collections.Token),
	)

	member := orgrepoimpl.NewMemberRepo(
		mongodb.NewCollection(cfg.Mongodb.Collections.Member),
	)
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

	git := usergit.NewUserGit(giteauser.GetClient(gitea.Client()))

	var orgAppService orgapp.OrgService

	permission := orgapp.NewPermService(
		&cfg.Permission, member)

	userAppService := userapp.NewUserService(
		user, git, token)

	orgAppService = orgapp.NewOrgService(
		userAppService, org, member, permission, cfg.Org.InviteExpiry,
	)

	{

		controller.AddRouterForUserController(
			v1, userAppService, user,
			authingUser,
		)

		controller.AddRouterForLoginController(
			v1, userAppService, authingUser, session,
		)

		controller.AddRouterForOrgController(
			v1, orgAppService, userAppService,
		)

	}
}

func codeRepoAppService(cfg *config.Config) coderepoapp.CodeRepoAppService {
	return coderepoapp.NewCodeRepoAppService(
		coderepoadapter.NewRepoAdapter(gitea.Client()),
	)
}

func setRouterOfModel(rg *gin.RouterGroup, services *allServices) {
	modelctl.AddRouteForModelController(
		rg,
		modelapp.NewModelAppService(services.codeRepoApp, services.modelRepoAdapter),
		services.userMiddleWare,
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

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

	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"

	"github.com/openmerlin/merlin-server/api"
	"github.com/openmerlin/merlin-server/config"
	"github.com/openmerlin/merlin-server/controller"
	"github.com/openmerlin/merlin-server/infrastructure/gitea"
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
	"github.com/openmerlin/merlin-server/models/infrastructure/modelrepositoryadapter"

	coderepoapp "github.com/openmerlin/merlin-server/coderepo/app"
)

func StartWebServer(port int, timeout time.Duration, cfg *config.Config) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logRequest())
	r.TrustedPlatform = "x-real-ip"

	if err := setRouter(r, cfg); err != nil {
		logrus.Error(err)

		return
	}

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadHeaderTimeout: time.Duration(cfg.ReadHeaderTimeout) * time.Second,
		Handler:           r,
	}

	defer interrupts.WaitForGracefulShutdown()

	interrupts.ListenAndServe(srv, timeout)
}

// setRouter init router
func setRouter(engine *gin.Engine, cfg *config.Config) error {
	api.SwaggerInfo.BasePath = "/api"
	api.SwaggerInfo.Title = "merlin"
	api.SwaggerInfo.Description = "set header: 'PRIVATE-TOKEN=xxx'"

	v1 := engine.Group(api.SwaggerInfo.BasePath)

	codeRepoApp := codeRepoAppService(v1, cfg)

	if err := setRouterOfModel(v1, cfg, codeRepoApp); err != nil {
		return err
	}

	setRouterOfUserAndOrg(v1, cfg)

	engine.UseRawPath = true
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return nil
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

	git := usergit.NewUserGit(gitea.GetClient())

	userAppService := userapp.NewUserService(
		user, git, token)

	orgAppService := orgapp.NewOrgService(
		userAppService, org, member, cfg.Org.InviteExpiry,
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
			v1, orgAppService,
		)

	}
}

func codeRepoAppService(
	v1 *gin.RouterGroup, cfg *config.Config,
) coderepoapp.CodeRepoAppService {
	return coderepoapp.NewCodeRepoAppService()
}

func setRouterOfModel(
	v1 *gin.RouterGroup, cfg *config.Config,
	codeRepoApp coderepoapp.CodeRepoAppService,
) error {
	repo, err := modelrepositoryadapter.NewModelAdapter(
		postgresql.DB(), &cfg.Model.Tables,
	)
	if err != nil {
		return err
	}

	modelctl.AddRouteForModelController(
		v1,
		modelapp.NewModelAppService(codeRepoApp, repo),
	)

	return nil
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

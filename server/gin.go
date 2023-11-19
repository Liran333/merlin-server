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

	"github.com/openmerlin/merlin-server/api"
	"github.com/openmerlin/merlin-server/config"
	"github.com/openmerlin/merlin-server/controller"
	"github.com/openmerlin/merlin-server/infrastructure/mongodb"
	"github.com/openmerlin/merlin-server/login/infrastructure/oidcimpl"
	session "github.com/openmerlin/merlin-server/session/app"
	sessionrepo "github.com/openmerlin/merlin-server/session/infrastructure"

	userapp "github.com/openmerlin/merlin-server/user/app"
	userrepoimpl "github.com/openmerlin/merlin-server/user/infrastructure/repositoryimpl"
)

func StartWebServer(port int, timeout time.Duration, cfg *config.Config) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logRequest())

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

	collections := &cfg.Mongodb.Collections

	user := userrepoimpl.NewUserRepo(
		mongodb.NewCollection(collections.User),
	)

	sessrepo := sessionrepo.NewSessionRepository(
		sessionrepo.NewSessionStore(
			mongodb.NewCollection(collections.Session),
		),
	)

	session := session.NewSessionService(sessrepo)

	authingUser := oidcimpl.NewAuthingUser()

	v1 := engine.Group(api.SwaggerInfo.BasePath)

	userAppService := userapp.NewUserService(
		user)

	{

		controller.AddRouterForUserController(
			v1, userAppService, user,
			authingUser,
		)

		controller.AddRouterForLoginController(
			v1, userAppService, authingUser, session,
		)

	}

	engine.UseRawPath = true
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

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

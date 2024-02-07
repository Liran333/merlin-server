package server

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/server-common-lib/interrupts"
	"github.com/sirupsen/logrus"

	"github.com/openmerlin/merlin-server/api"
	"github.com/openmerlin/merlin-server/config"
)

const (
	version         = "development" // program version for this build
	apiDesc         = "Modelfoundry server APIs"
	apiTitle        = "Modelfoundry"
	waitServerStart = 3 // 3s
)

func StartWebServer(key, cert string, removeCfg bool, port int, timeout time.Duration, cfg *config.Config) {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(logRequest())
	engine.UseRawPath = true
	engine.TrustedPlatform = "x-real-ip"

	api.SwaggerInfo.Title = apiTitle
	api.SwaggerInfo.Version = version
	api.SwaggerInfo.Description = apiDesc

	// init services
	services, err := initServices(cfg)
	if err != nil {
		logrus.Error(err)

		return
	}

	// web api
	setRouterOfWeb("/web", engine, cfg, &services)

	// restful api
	setRouterOfRestful("/api", engine, cfg, &services)

	// internal service api
	setRouterOfInternal("/internal", engine, cfg, &services)

	// start server
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           engine,
		ReadHeaderTimeout: time.Duration(cfg.ReadHeaderTimeout) * time.Second,
	}

	defer interrupts.WaitForGracefulShutdown()

	if key != "" && cert != "" {
		srv.TLSConfig = &tls.Config{
			MinVersion:               tls.VersionTLS12, // tls1.3 cipher suite is not configurable
			MaxVersion:               tls.VersionTLS13,
			PreferServerCipherSuites: true,
		}
		interrupts.ListenAndServeTLS(srv, cert, key, timeout)
		// wait server start
		time.Sleep(time.Duration(waitServerStart) * time.Second)
		if removeCfg {
			if err := os.Remove(cert); err != nil {
				logrus.Errorf("remove cert file: %s", err)
			}

			if err := os.Remove(key); err != nil {
				logrus.Errorf("remove key file: %s", err)
			}
		}

	} else {
		interrupts.ListenAndServe(srv, timeout)

	}
}

func logRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()

		errmsg := ""
		for _, ginErr := range c.Errors {
			if errmsg != "" {
				errmsg += ","
			}
			errmsg = fmt.Sprintf("%s%s", errmsg, ginErr.Error())
		}

		log := fmt.Sprintf(
			"| %d | %d | %s | %s ",
			c.Writer.Status(),
			endTime.Sub(startTime),
			c.Request.Method,
			c.Request.RequestURI,
		)
		if errmsg != "" {
			log += fmt.Sprintf("| %s ", errmsg)
		}

		logrus.Info(log)
	}
}

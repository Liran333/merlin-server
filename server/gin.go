/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package server provides functionality for setting up and configuring a server for handling code repo operations.
package server

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/server-common-lib/interrupts"
	"github.com/sirupsen/logrus"

	_ "github.com/openmerlin/merlin-server/api"
	"github.com/openmerlin/merlin-server/config"
)

const (
	version         = "development" // program version for this build
	apiDesc         = "Modelfoundry server APIs"
	apiTitle        = "Modelfoundry"
	waitServerStart = 3 // 3s
)

// StartWebServer starts a web server with the given configuration.
// It initializes the services, sets up the routers for different APIs, and starts the server.
// If TLS key and certificate are provided, it will use HTTPS.
// If removeCfg is true, it will remove the key and certificate files after starting the server.
func StartWebServer(key, cert string, removeCfg bool, port int, timeout time.Duration, cfg *config.Config) {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(logRequest())
	engine.UseRawPath = true
	engine.TrustedPlatform = "x-real-ip"

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
				logrus.Fatalf("remove cert file: %s", err.Error())
			}

			if err := os.Remove(key); err != nil {
				logrus.Fatalf("remove key file: %s", err.Error())
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

		for _, ginErr := range c.Errors {
			logrus.Errorf("error on %s %s:\n%+v", c.Request.Method, c.Request.RequestURI, ginErr.Unwrap())
		}

		if strings.Contains(c.Request.RequestURI, "/swagger/") ||
			strings.Contains(c.Request.RequestURI, "/internal/heartbeat") {
			return
		}

		log := fmt.Sprintf(
			"| %d | %d | %s | %s ",
			c.Writer.Status(),
			endTime.Sub(startTime),
			c.Request.Method,
			c.Request.RequestURI,
		)

		logrus.Info(log)
	}
}

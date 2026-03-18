package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rshby/go-event-ticketing/config"
	"github.com/rshby/go-event-ticketing/internal/database"
	"github.com/rshby/go-event-ticketing/tracing"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	ServerCmd = &cobra.Command{
		Use: "server",
		Run: server,
	}
)

func init() {
	RootCmd.AddCommand(ServerCmd)
}

func server(cmd *cobra.Command, args []string) {
	logrus.Infof("running Server🚀")

	// connect to open telemetry
	tracerProvider, err := tracing.ConnectOTLPTrace(context.Background())
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = tracerProvider.Shutdown(context.Background())
	}()

	// connect to Redis
	redisClient := database.ConnectRedis()
	defer func() {
		_ = redisClient.Close()
	}()

	// connect to PostgreSql
	db, err := database.ConnectPostgreSql()
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		if db != nil {
			if sqlDb, errDb := db.DB(); errDb == nil {
				_ = sqlDb.Close()
				logrus.Infof("close database🛑")
			}
		}
	}()

	app := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	// create http server
	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%s", config.AppPort()),
		Handler:           app,
		ReadHeaderTimeout: config.HttpServerReadHeaderTimeout(),
		ReadTimeout:       config.HttpServerReadTimeout(),
		WriteTimeout:      config.HttpServerWriteTimeout(),
		IdleTimeout:       config.HttpServerIdleTimeout(),
	}

	waitChan := make(chan struct{})
	errChan := make(chan error, 1)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// running http server
	go func() {
		logrus.Infof("Running Http Server on Port %s ⌛️", config.AppPort())
		errHttp := httpServer.ListenAndServe()
		if errHttp != nil && !errors.Is(errHttp, http.ErrServerClosed) {
			logrus.Error(errHttp)
			errChan <- errHttp
			return
		}
	}()

	go func() {
		select {
		case <-errChan:
			gracefullShutdown(httpServer)
			waitChan <- struct{}{}
			return
		case <-signalChan:
			logrus.Infof("receive signal⚠️")
			gracefullShutdown(httpServer)
			waitChan <- struct{}{}
			return
		}
	}()

	<-waitChan
}

// gracefullShutdown gracefull
func gracefullShutdown(httpServer *http.Server) {
	if httpServer != nil {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancelFunc()

		logrus.Info("shutdown http server")
		if err := httpServer.Shutdown(ctx); err != nil {
			logrus.Error(err)
			logrus.Info("force close http server")
			_ = httpServer.Close()
		}
	}
}

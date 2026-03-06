package cmd

import (
	"context"

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
}

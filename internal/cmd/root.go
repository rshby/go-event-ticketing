package cmd

import (
	"github.com/rshby/go-event-ticketing/config"
	"github.com/rshby/go-event-ticketing/internal/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:   "go-event-ticketing",
		Short: "go-event-ticketing",
	}
)

func init() {
	logger.SetupLogger()
	config.LoadConfig()
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

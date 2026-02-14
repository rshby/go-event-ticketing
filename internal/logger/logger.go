package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func SetupLogger() {
	logrus.SetOutput(os.Stdout)
}

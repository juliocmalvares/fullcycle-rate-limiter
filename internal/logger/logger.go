package logger

import (
	"github.com/sirupsen/logrus"
)
var Logger *logrus.Logger

func Init() {
	Logger = logrus.New()
	Logger.SetOutput(logrus.StandardLogger().Writer())
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	Logger.SetLevel(logrus.InfoLevel)
}

func Error(msg string, err error) {
	Logger.WithFields(logrus.Fields{
		"error": err,
	}).Error(msg)
}
func Info(msg string) {
	Logger.Info(msg)
}
func Debug(msg string) {
	Logger.Debug(msg)
}

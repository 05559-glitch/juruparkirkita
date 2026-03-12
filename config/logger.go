package config

import (
	"os"
	"github.com/sirupsen/logrus"
)

func InitLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true, 
	})

	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)
	
	return log
}

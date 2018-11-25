package main

import (
	"os"

	"github.com/sirupsen/logrus"

	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var baseLog = logrus.New()

type logPrefixes struct {
	web *logrus.Entry
	lib *logrus.Entry
}

var log logPrefixes

func init() {
	baseLog.Formatter = &prefixed.TextFormatter{}
	baseLog.Level = logrus.InfoLevel
	baseLog.SetOutput(os.Stdout)

	log = logPrefixes{
		web: baseLog.WithField("prefix", "web"),
		lib: baseLog.WithField("prefix", "lib"),
	}

}

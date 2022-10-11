package logging

import (
	"github.com/sirupsen/logrus"
	"os"
)

var (
	Log            = logrus.New()
	failLogger     = logrus.New()
	stableLogger   = logrus.New()
	restoredLogger = logrus.New()
)

func LoggingSetup() {

	Log.SetFormatter(&logrus.JSONFormatter{})
	failLogger.SetFormatter(&logrus.JSONFormatter{})
	stableLogger.SetFormatter(&logrus.JSONFormatter{})
	restoredLogger.SetFormatter(&logrus.JSONFormatter{})

	file, err := os.OpenFile("telmon.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Log.Out = file
		failLogger.Out = file
		stableLogger.Out = file
		restoredLogger.Out = file
	} else {
		Log.Info("Failed to log to file, using default stderr")
	}
}

var (
	FailLogger     = failLogger.WithFields(logrus.Fields{"kind": "failed"})
	StableLogger   = stableLogger.WithFields(logrus.Fields{"kind": "stable"})
	RestoredLogger = restoredLogger.WithFields(logrus.Fields{"kind": "restored"})
)

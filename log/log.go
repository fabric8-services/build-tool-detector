/*

Package log is used to log problems
with the service using sentry and logorus.
If a sentry_dsn  is not provided, the
log output will be sent to stdout.

*/
package log

import (
	"os"

	"github.com/evalphobia/logrus_sentry"
	"github.com/fabric8-services/build-tool-detector/config"
	"github.com/sirupsen/logrus"
)

const (
	// SentryDSN for sentry
	SentryDSN = "SENTRY_DSN"
)

const (
	buildToolDetector = "build-tool-detector"
	applicationName   = "applicationName"
)

// Logger something
func Logger() *logrus.Entry {
	var configuration config.Configuration
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.WarnLevel)

	// TODO: have env variable to specify we are running tests which
	// will return a different logger.
	if configuration.GetSentryDSN() != "" {
		hook, err := logrus_sentry.NewSentryHook(configuration.GetSentryDSN(), []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		})

		if err != nil {
			panic(err)
		}

		// Add sentry hook
		logrus.AddHook(hook)
	} else {
		logrus.SetOutput(os.Stdout)
	}
	return logrus.WithField(applicationName, buildToolDetector)
}

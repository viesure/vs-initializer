package logging

import (
	"strings"

	"github.com/hirosassa/zerodriver"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"viesure.io/vs-initializer/pkg/config"
)

func InitLogger() {

	logger := zerodriver.NewProductionLogger()

	logLevel := zerolog.InfoLevel
	logLevelConfig := strings.ToUpper(viper.GetString(config.AppLogLevelVarName))
	switch logLevelConfig {
	case "DEBUG":
		logLevel = zerolog.DebugLevel
	case "INFO":
		logLevel = zerolog.InfoLevel
	case "WARN":
		logLevel = zerolog.WarnLevel
	case "ERROR":
		logLevel = zerolog.ErrorLevel
	default:
		logLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(logLevel)

	viper.Set(config.LoggerVarName, logger)

}

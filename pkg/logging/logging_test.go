package logging

import (
	"testing"

	"github.com/hirosassa/zerodriver"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"viesure.io/vs-initializer/pkg/config"
)

type LoggingSuite struct {
	suite.Suite
}

func TestLogging(t *testing.T) {

	suite.Run(t, &LoggingSuite{})

}

func (ls *LoggingSuite) SetupSuite() {

	viper.Set(config.AppLogLevelVarName, config.DefaultAppLogLevel)

}

func (ls *LoggingSuite) TestLogger() {

	InitLogger()

	logger := viper.Get(config.LoggerVarName)

	ls.NotNil(logger)

	ls.IsType(&zerodriver.Logger{}, logger)

	ls.IsType(zerolog.InfoLevel, logger.(*zerodriver.Logger).GetLevel())

}

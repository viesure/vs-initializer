package files

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"testing"

// 	"github.com/hirosassa/zerodriver"
// 	"github.com/rs/zerolog"
// 	"github.com/spf13/viper"
// 	"viesure.io/vs-initializer/pkg/config"
// 	"viesure.io/vs-initializer/pkg/svc/gcp"
// )

// var logger *zerodriver.Logger

// func TestNewFakeSecretManagerClient(t *testing.T) {

// 	//asserts := assert.New(t)

// 	viper.Set(config.AppLogLevelVarName, "DEBUG")

// 	// initialize logging
// 	logger = zerodriver.NewProductionLogger()

// 	var logLevel zerolog.Level
// 	logLevelConfig := strings.ToUpper(viper.GetString(config.AppLogLevelVarName))
// 	switch logLevelConfig {
// 	case "DEBUG":
// 		logLevel = zerolog.DebugLevel
// 	case "INFO":
// 		logLevel = zerolog.InfoLevel
// 	case "WARN":
// 		logLevel = zerolog.WarnLevel
// 	case "ERROR":
// 		logLevel = zerolog.ErrorLevel
// 	default:
// 		logLevel = zerolog.InfoLevel
// 	}

// 	zerolog.SetGlobalLevel(logLevel)

// 	appContext := context.Background()
// 	appContext = context.WithValue(appContext, config.LoggerContextKey{}, logger)

// 	client, err := gcp.NewSecretManagerClient(appContext)
// 	if err != nil {
// 		logger.Error().Msgf("error creating gcp secret manager client: %v", err)
// 		os.Exit(1)
// 	}

// 	secretManagerService, err := gcp.NewSecretManagerService(appContext, client)
// 	if err != nil {
// 		logger.Error().Msgf("error creating gcp secret manager service: %v", err)
// 		os.Exit(1)
// 	}
// 	viper.Set(config.SecretManagerServiceComponentName, secretManagerService)

// 	basePath, err := filepath.Abs("./../../..")
// 	if err != nil {
// 		t.Errorf("error while checking base path: %v", err)
// 	}
// 	logger.Debug().Msgf("base path: %s", basePath)

// 	templatePath := fmt.Sprintf("%s%s%s", basePath, string(os.PathSeparator), "tests/template")
// 	viper.Set(config.TemplateDirVarName, templatePath)

// 	outputPath := fmt.Sprintf("%s%s%s", basePath, string(os.PathSeparator), "tests.out")
// 	viper.Set(config.OutputDirVarName, outputPath)

// 	if _, err := os.Stat(outputPath); err != nil {
// 		err := os.MkdirAll(outputPath, 0755)
// 		if err != nil {
// 			t.Errorf("error creating output directory: %v", err)
// 		}
// 	}

// 	fileService := NewFileService(appContext)
// 	fileService.ReplaceSecretsInFiles()

// }

package main

import (
	"context"
	"os"

	"github.com/hirosassa/zerodriver"
	"github.com/spf13/viper"
	"viesure.io/vs-initializer/pkg/config"
	"viesure.io/vs-initializer/pkg/logging"
	"viesure.io/vs-initializer/pkg/svc/files"
	"viesure.io/vs-initializer/pkg/svc/gcp"
	"viesure.io/vs-initializer/pkg/svc/k8s"
	"viesure.io/vs-initializer/pkg/utils"
)

var Version = "latest"
var appContext context.Context
var logger *zerodriver.Logger

func init() {

	appContext = context.Background()

	// config variables
	viper.Set(config.VersionVarName, Version)
	viper.Set(config.AppLogLevelVarName, utils.GetEnvOrDefault(config.AppLogLevelVarName, config.DefaultAppLogLevel))
	viper.Set(config.EnvSecretVarName, utils.GetEnvOrDefault(config.EnvSecretVarName, config.DefaultEnvSecret))
	viper.Set(config.AppNamespaceVarName, utils.GetEnvOrDefault(config.AppNamespaceVarName, config.DefaultAppNamespace))
	viper.Set(config.TemplateDirVarName, utils.NormalizeDirectoryPath(utils.GetEnvOrDefault(config.TemplateDirVarName, config.DefaultTemplateDir)))
	viper.Set(config.OutputDirVarName, utils.NormalizeDirectoryPath(utils.GetEnvOrDefault(config.OutputDirVarName, config.DefaultOutputDir)))

	// initialize logging
	logging.InitLogger()
	logger = viper.Get(config.LoggerVarName).(*zerodriver.Logger)

	// config information
	logger.Info().Msg("----------- configuration -----------")
	for key, value := range viper.AllSettings() {
		if key == config.LoggerVarName {
			continue
		}
		logger.Info().Msgf("%s: '%v'", key, value)
	}
	logger.Info().Msg("-------------------------------------")

	// services
	viper.Set(config.FileServiceComponentName, files.NewFileService())

	client, err := gcp.NewSecretManagerClient(appContext)
	if err != nil {
		logger.Error().Msgf("error creating gcp secret manager client: %v", err)
		os.Exit(1)
	}

	secretManagerService, err := gcp.NewSecretManagerService(client)
	if err != nil {
		logger.Error().Msgf("error creating gcp secret manager service: %v", err)
		os.Exit(1)
	}
	viper.Set(config.SecretManagerServiceComponentName, secretManagerService)

	secretsService, err := k8s.NewSecretsService(appContext)
	if err != nil {
		logger.Error().Msgf("error creating k8s secrets service: %v", err)
		os.Exit(1)
	}
	viper.Set(config.SecretsServiceComponentName, secretsService)

}

func main() {
	logger.Info().Msgf("starting initializer %s", viper.GetString(config.VersionVarName))

	secretManagerService := viper.Get(config.SecretManagerServiceComponentName).(*gcp.SecretManagerService)
	defer secretManagerService.Shutdown()

	fileService := viper.Get(config.FileServiceComponentName).(*files.FileService)

	err := fileService.ReplaceSecretsInFiles()
	if err != nil {
		logger.Error().Msgf("error while replacing secrets in files: %v", err)
	}

	logger.Info().Msg("finished processing, good bye and thanks for all the fish!")

}

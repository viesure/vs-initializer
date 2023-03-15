package config

import (
	"viesure.io/vs-initializer/pkg/utils"
)

const (
	VersionVarName = "VERSION"

	AppLogLevelVarName = "APP_LOG_LEVEL"
	DefaultAppLogLevel = "INFO"

	TemplateDirVarName = "TEMPLATE_DIR"
	DefaultTemplateDir = "/data.tmpl"

	OutputDirVarName = "OUTPUT_DIR"
	DefaultOutputDir = "/data"

	EnvSecretVarName = "ENV_SECRET"
	DefaultEnvSecret = "app-env"

	AppNamespaceVarName = "APP_NAMESPACE"

	QuoteChars = "'\""

	EnvFileName = ".env"

	SecretManagerProtocol = "sm"
	UrlSearchPrefix       = SecretManagerProtocol + "://"

	InClusterNamespaceLocation = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

	SecretManagerServiceComponentName = "secretManagerService"
	FileServiceComponentName          = "fileService"
	SecretsServiceComponentName       = "secretsService"

	LoggerVarName = "logger"
)

var (
	DefaultAppNamespace = utils.ReadFileOrEmpty(InClusterNamespaceLocation)
)

package files

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/hirosassa/zerodriver"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"viesure.io/vs-initializer/pkg/config"
	"viesure.io/vs-initializer/pkg/svc/gcp"
	"viesure.io/vs-initializer/pkg/svc/k8s"
)

type (
	FileServiceInterface interface {
		ReplaceSecretsInFiles() error
		readFiles(directory string) ([]string, error)
		processEnvironmentFile(sourceFilename string) error
		processFile(sourceFilename string, destinationFilename string) error
		replaceUrlsInLine(line string, stripComments bool) string
	}

	FileService struct {
		logger *zerodriver.Logger
	}
)

func NewFileService() *FileService {

	return &FileService{
		logger: viper.Get(config.LoggerVarName).(*zerodriver.Logger),
	}

}

func (svc *FileService) ReplaceSecretsInFiles() error {

	files, err := svc.readFiles(viper.GetString(config.TemplateDirVarName))
	if err != nil {
		return err
	}

	for _, file := range files {

		sourcefile := viper.GetString(config.TemplateDirVarName) + string(os.PathSeparator) + file
		svc.logger.Debug().Msgf("source file: %s", sourcefile)

		if file == config.EnvFileName {
			err := svc.processEnvironmentFile(sourcefile)
			if err != nil {
				svc.logger.Error().Msgf("error while replacing secrets in environment file '%s': %v", sourcefile, err)
				os.Exit(1)
			}
			continue
		}

		destinationfile := viper.GetString(config.OutputDirVarName) + string(os.PathSeparator) + file
		svc.logger.Debug().Msgf("destination file: %s", destinationfile)

		err := svc.processFile(sourcefile, destinationfile)
		if err != nil {
			svc.logger.Error().Msgf("error while replacing secrets in file '%s': %v", sourcefile, err)
			os.Exit(1)
		}

	}

	return nil

}

func (svc *FileService) readFiles(directory string) ([]string, error) {

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return nil, err
	}

	files := []string{}

	f, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, file := range f {
		if file.IsDir() {
			continue
		}
		if len(file.Name()) < 2 {
			continue
		}
		if file.Name()[0:2] == ".." {
			continue
		}

		files = append(files, file.Name())
	}

	return files, nil

}

// processEnvironmentFile parses the file specified as `sourceFilename`, e.g.:
//
//	"/srv/data.tmpl/.env"
//
// and puts the content to a k8s secret named as the environment variable
// `ENV_SECRET` (or it's default `app-env`).
// It returns an error if one occurs.
func (svc *FileService) processEnvironmentFile(sourceFilename string) error {
	svc.logger.Info().Msgf("searching and replacing secret manager urls inside environment file '%s'", sourceFilename)

	file, err := os.Open(sourceFilename)
	if err != nil {
		return err
	}
	defer file.Close()

	envMap := map[string]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		svc.logger.Debug().Msgf("read line from file '%s': '%s'", sourceFilename, line)
		if !strings.Contains(line, "=") {
			svc.logger.Debug().Msg("skipping line because it does not contain an '='")
			continue
		}
		if strings.TrimSpace(line)[0:1] == "#" {
			svc.logger.Debug().Msg("skipping line because it is commented out")
			continue
		}
		line = svc.replaceUrlsInLine(line, true)
		if strings.Contains(line, "=") {
			lineParts := strings.SplitN(line, "=", 2)
			envMap[lineParts[0]] = lineParts[1]
		}
	}
	svc.logger.Debug().Msgf("map for environment variables has %d entries", len(envMap))

	k8sSecretsService := viper.Get(config.SecretsServiceComponentName).(*k8s.SecretsService)

	svc.logger.Debug().Msgf("getting k8s secret '%s/%s'", viper.GetString(config.AppNamespaceVarName), viper.GetString(config.EnvSecretVarName))
	secret, err := k8sSecretsService.Get(viper.GetString(config.EnvSecretVarName), metav1.GetOptions{})
	if err != nil {
		return err
	}
	svc.logger.Debug().Msgf("successfully got k8s secret '%s/%s'", secret.Namespace, secret.Name)

	svc.logger.Debug().Msgf("updating k8s secret '%s/%s'", secret.Namespace, secret.Name)
	secret.StringData = envMap
	newSecret, err := k8sSecretsService.Update(secret, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	svc.logger.Debug().Msgf("successfully updated k8s secret '%s/%s'", newSecret.Namespace, newSecret.Name)

	svc.logger.Info().Msgf("done processing environment file '%s'", sourceFilename)

	return nil
}

// processFile parses and replaces secret manager URLs with the actual secret
// value inside the file specified as `sourceFilename`, e.g.:
//
//	"/srv/data.tmpl/app-config.json"
//
// and puts the processed content to a file specified as
// `destinationFilename`, e.g.:
//
//	"/srv/data/app-config.json"
//
// It returns an error if one occurs.
func (svc *FileService) processFile(sourceFilename string, destinationFilename string) error {
	svc.logger.Info().Msgf("searching and replacing secrets inside '%s' -> '%s'", sourceFilename, destinationFilename)

	sourceFile, err := os.Open(sourceFilename)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	if _, err := os.Stat(destinationFilename); err == nil {
		err := os.Remove(destinationFilename)
		if err != nil {
			return err
		}
	}

	destinationFile, err := os.Create(destinationFilename)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	scanner := bufio.NewScanner(sourceFile)
	for scanner.Scan() {
		line := scanner.Text()
		svc.logger.Debug().Msgf("read line from file '%s': '%s'", sourceFilename, line)
		processLine := true
		if strings.TrimSpace(line) == "" {
			svc.logger.Debug().Msg("skip processing line because it is empty")
			processLine = false
		}
		if !processLine || strings.TrimLeft(line, " ")[0:1] == "#" {
			svc.logger.Debug().Msg("skip processing line because it is commented out")
			processLine = false
		}
		if processLine {
			line = fmt.Sprintln(svc.replaceUrlsInLine(line, false))
		} else {
			line = fmt.Sprintln(line)
		}
		destinationFile.Write([]byte(line))
	}
	svc.logger.Info().Msgf("done processing file '%s'", sourceFilename)

	return nil

}

func (svc *FileService) replaceUrlsInLine(line string, stripComments bool) string {

	// filter out empty and commented lines
	if len(strings.TrimSpace(line)) == 0 || strings.TrimSpace(line)[0:1] == "#" {
		if stripComments {
			return ""
		} else {
			return line
		}
	}

	// filter out lines without URLs
	if !strings.Contains(line, config.UrlSearchPrefix) {
		svc.logger.Debug().Msgf("skipping line because URL prefix '%s' was not found", config.UrlSearchPrefix)
		return line
	}

	// line passed sanity checks, let's replace some URLs!
	index := 0
	for {
		index = strings.Index(line, config.UrlSearchPrefix)
		if index < 0 || (strings.Contains(line, "#") && strings.Index(line, "#") < index) {
			break
		}
		svc.logger.Debug().Msgf("line part in front of match: '%s'", line[0:index])
		remainingLine := line[index:]
		svc.logger.Debug().Msgf("found URL, remaining line: '%s'", remainingLine)
		quoteCharacter := ""
		for _, character := range config.QuoteChars {
			if index > 0 && strings.Contains(line[index-1:index], string(character)) {
				quoteCharacter = string(character)
				break
			}
		}
		url := ""
		if quoteCharacter != "" {
			svc.logger.Debug().Msgf("url is quoted by '%s'", quoteCharacter)
			url = remainingLine[0:strings.Index(remainingLine, quoteCharacter)]
		} else {
			svc.logger.Debug().Msgf("url is not quoted")
			if strings.Contains(remainingLine, " ") {
				url = remainingLine[0:strings.Index(remainingLine, " ")]
			} else {
				url = remainingLine
			}
		}
		svc.logger.Debug().Msgf("url is '%s'", url)

		secretManagerService := viper.Get(config.SecretManagerServiceComponentName).(*gcp.SecretManagerService)

		secret, err := secretManagerService.GetSecret(url)
		if err != nil {
			svc.logger.Error().Msgf("error getting secret for url '%s': %v", url, err)
			os.Exit(1)
		}

		svc.logger.Debug().Msgf("replacing url '%s' with secret", url)
		line = strings.ReplaceAll(line, url, secret.SecretValue)

	}

	if stripComments && strings.Contains(line, "#") {
		svc.logger.Debug().Msgf("stripping comments and TrimRight line")
		line = line[0:strings.Index(line, "#")]
		line = strings.TrimRight(line, " ")
	}

	return line

}

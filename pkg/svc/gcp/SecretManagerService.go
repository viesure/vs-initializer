package gcp

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/hirosassa/zerodriver"
	"github.com/spf13/viper"
	"google.golang.org/api/iterator"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"viesure.io/vs-initializer/pkg/config"
	"viesure.io/vs-initializer/pkg/models"
)

type (
	SecretManagerServiceInterface interface {
		GetSecret(smUrl string) (*models.Secret, error)
		Shutdown()
		getSecretsOfProject(project string) error
		getSecretValue(secret *models.Secret, version string) error
	}

	SecretManagerService struct {
		secrets             map[string][]*models.Secret
		secretManagerClient *SecretManagerClient
		logger              *zerodriver.Logger
	}
)

func NewSecretManagerService(secretManagerClient *SecretManagerClient) (*SecretManagerService, error) {

	secretManagerClient, err := NewSecretManagerClient(secretManagerClient.appContext)
	if err != nil {
		return nil, err
	}
	return &SecretManagerService{
		secrets:             make(map[string][]*models.Secret),
		secretManagerClient: secretManagerClient,
		logger:              viper.Get(config.LoggerVarName).(*zerodriver.Logger),
	}, nil

}

func (svc *SecretManagerService) GetSecret(smUrl string) (*models.Secret, error) {

	u, err := url.Parse(smUrl)
	if err != nil {
		return nil, err
	}
	project := u.Host
	if _, ok := svc.secrets[project]; !ok {
		err := svc.getSecretsOfProject(project)
		if err != nil {
			return nil, err
		}
	}

	pathParts := strings.Split(u.Path[1:], "/")

	labels := map[string]string{
		"secret-name": pathParts[0],
	}

	for key, value := range u.Query() {
		labels[key] = value[0]
	}

	var s *models.Secret
	for _, secret := range svc.secrets[project] {
		if reflect.DeepEqual(secret.Labels, labels) {
			svc.logger.Debug().Msgf("secret '%v' found in cache", secret.Name)
			s = secret
			break
		}
	}

	if s == nil {
		return nil, errors.New("did not find any secret that matches the labels")
	}

	if s.SecretValue == "" {
		svc.logger.Debug().Msgf("no secret value for secret '%v' in cache, fetching data for secret ...", s.Name)
		version := "latest"
		if len(pathParts) > 1 {
			version = pathParts[1]
		}
		err := svc.getSecretValue(s, version)
		if err != nil {
			return nil, err
		}
	}

	return s, nil

}

// getSecrets retrieves the attributes of all secrets in the project,
// given by the `project`, e.g.:
//
//	"my-project"
//
// and stores them into the local variable `secrets` as cache. It returns
// an error if one occurs.
//
// Note: this function does not fetch the secret value of the secret as this
// could be very time-consuming.
func (svc *SecretManagerService) getSecretsOfProject(project string) error {

	svc.logger.Debug().Msgf("fetching secrets for project '%v'", project)
	req := &secretmanagerpb.ListSecretsRequest{
		Parent: "projects/" + project,
	}

	it := svc.secretManagerClient.ListSecrets(req)

	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			svc.logger.Info().Msgf("failed to query secret list: %v", err)
			return err
		}

		nameParts := strings.Split(resp.Name, "/")
		shortName := nameParts[len(nameParts)-1]

		labels := make(map[string]string)
		for key, value := range resp.Labels {
			labels[key] = value
		}

		secret := models.Secret{
			Name:      resp.Name,
			ShortName: shortName,
			Labels:    labels,
		}

		svc.secrets[project] = append(svc.secrets[project], &secret)

	}

	return nil

}

// getSecretValue retrieves the secret value of the secret given as `secret`
// and the version given by `version`, e.g.:
//
//	"latest"
//	"3"
//
// and stores it into the parameter `secret`. It returns
// an error if one occurs.
func (svc *SecretManagerService) getSecretValue(secret *models.Secret, version string) error {

	versionName := fmt.Sprintf("%s/versions/%s", secret.Name, version)
	accessSecretVersionRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: versionName,
	}

	result, err := svc.secretManagerClient.AccessSecretVersion(accessSecretVersionRequest)
	if err != nil {
		svc.logger.Info().Msgf("failed to access secret value: %v", err)
		return err
	}

	secret.VersionName = result.Name
	secret.SecretValue = string(result.Payload.Data)

	return nil

}

func (svc *SecretManagerService) Shutdown() {
	defer svc.secretManagerClient.Close()
}

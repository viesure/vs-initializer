package gcp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"gopkg.in/yaml.v3"
	"viesure.io/vs-initializer/pkg/models"
)

type (
	FakeSecretManagerClientInterface interface {
		ListSecrets(req *secretmanagerpb.ListSecretsRequest) *secretmanager.SecretIterator
		AccessSecretVersion(req *secretmanagerpb.AccessSecretVersionRequest) (*secretmanagerpb.AccessSecretVersionResponse, error)
		Close()
	}
	FakeSecretManagerClient struct {
		appContext context.Context
		secrets    []*secretmanagerpb.Secret
	}
)

const (
	TESTSECRETSFILENAME = "tests/secrets.yaml"
)

func NewFakeSecretManagerClient(appContext context.Context) (*FakeSecretManagerClient, error) {

	secrets, err := readSecrets()
	if err != nil {
		return nil, err
	}

	return &FakeSecretManagerClient{
		appContext: appContext,
		secrets:    secrets,
	}, nil

}

func readSecrets() ([]*secretmanagerpb.Secret, error) {

	testFile, err := filepath.Abs(fmt.Sprintf("../../../%s", TESTSECRETSFILENAME))
	if err != nil {
		return nil, err
	}

	contents, err := os.ReadFile(testFile)
	if err != nil {
		return nil, err
	}

	var secrets []models.Secret
	err = yaml.Unmarshal(contents, &secrets)
	if err != nil {
		return nil, err
	}

	var smSecrets []*secretmanagerpb.Secret
	for _, v := range secrets {
		secret := &secretmanagerpb.Secret{
			Name:   v.Name,
			Labels: v.Labels,
		}
		smSecrets = append(smSecrets, secret)
	}

	return smSecrets, nil

}

func (c *FakeSecretManagerClient) ListSecrets(req *secretmanagerpb.ListSecretsRequest) *secretmanager.SecretIterator {

	it := &secretmanager.SecretIterator{
		Response: secretmanagerpb.ListSecretsResponse{
			Secrets:   c.secrets,
			TotalSize: int32(len(c.secrets)),
		},
	}

	return it

}

func (c *FakeSecretManagerClient) AccessSecretVersion(ctx context.Context, req secretmanagerpb.AccessSecretVersionRequest) (*secretmanagerpb.AccessSecretVersionResponse, error) {

	resp := &secretmanagerpb.AccessSecretVersionResponse{}

	return resp, nil

}

func (c *FakeSecretManagerClient) Close() {}

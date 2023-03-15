package gcp

import (
	"context"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type (
	SecretManagerClientInterface interface {
		ListSecrets(req *secretmanagerpb.ListSecretsRequest) *secretmanager.SecretIterator
		AccessSecretVersion(req *secretmanagerpb.AccessSecretVersionRequest) (*secretmanagerpb.AccessSecretVersionResponse, error)
		Close()
	}

	SecretManagerClient struct {
		appContext     context.Context
		internalClient *secretmanager.Client
	}
)

func NewSecretManagerClient(appContext context.Context) (*SecretManagerClient, error) {

	client, err := secretmanager.NewClient(appContext)
	if err != nil {
		return nil, err
	}

	return &SecretManagerClient{
		appContext:     appContext,
		internalClient: client,
	}, nil

}

func (c *SecretManagerClient) ListSecrets(req *secretmanagerpb.ListSecretsRequest) *secretmanager.SecretIterator {

	return c.internalClient.ListSecrets(c.appContext, req)

}

func (c *SecretManagerClient) AccessSecretVersion(req *secretmanagerpb.AccessSecretVersionRequest) (*secretmanagerpb.AccessSecretVersionResponse, error) {

	return c.internalClient.AccessSecretVersion(c.appContext, req)

}

func (c *SecretManagerClient) Close() {
	defer c.internalClient.Close()
}

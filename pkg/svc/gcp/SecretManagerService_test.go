package gcp_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"viesure.io/vs-initializer/pkg/svc/gcp"
)

type SecretManagerServiceSuite struct {
	suite.Suite
	client *gcp.FakeSecretManagerClient
}

func TestSecretManagerServiceSuite(t *testing.T) {
	suite.Run(t, &SecretManagerServiceSuite{})
}

func (ts *SecretManagerServiceSuite) SetupSuite() {

	testContext := context.Background()

	client, err := gcp.NewFakeSecretManagerClient(testContext)
	if err != nil {
		ts.FailNowf("error while creating client", err.Error())
	}

	ts.client = client

}

func (ts *SecretManagerServiceSuite) TearDownSuite() {

	ts.client.Close()

}

func (ts *SecretManagerServiceSuite) TestClientNotNil() {

	ts.NotNil(ts.client)

	ts.IsType(&gcp.FakeSecretManagerClient{}, ts.client)

}

package gcp_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"viesure.io/vs-initializer/pkg/svc/gcp"
)

type FakeSecretManagerClientSuite struct {
	suite.Suite
	client *gcp.FakeSecretManagerClient
}

var appContext context.Context

func TestFakeSecretManagerClientSuite(t *testing.T) {
	suite.Run(t, &FakeSecretManagerClientSuite{})
}

func (ts *FakeSecretManagerClientSuite) SetupSuite() {

	appContext = context.Background()

	client, err := gcp.NewFakeSecretManagerClient(appContext)
	if err != nil {
		ts.FailNow("error while creating client", err.Error())
	}

	ts.client = client

}

func (ts *FakeSecretManagerClientSuite) TearDownSuite() {

	ts.client.Close()

}

func (ts *FakeSecretManagerClientSuite) TestNewFakeSecretManagerClient() {

	// ts.Equal(len(ts.client.secrets), 2)

}

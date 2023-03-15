package utils_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
	"viesure.io/vs-initializer/pkg/utils"
)

const (
	ENV_VAR_NAME     = "TEST_ENVVAR"
	ENV_VAR_VALUE    = "TEST_ENVVAR_VALUE"
	TESTREADFILENAME = "testread.txt"
)

type UtilsSuite struct {
	suite.Suite
	tempDir string
}

func TestUtils(t *testing.T) {
	suite.Run(t, &UtilsSuite{})
}

func (us *UtilsSuite) SetupSuite() {
	err := os.Setenv(ENV_VAR_NAME, ENV_VAR_VALUE)
	us.Nil(err)

	tempDir, err := os.MkdirTemp("", "vs-initializer")
	us.NoError(err)

	us.tempDir = tempDir
}

func (us *UtilsSuite) TearDownSuite() {
	os.Unsetenv(ENV_VAR_NAME)
	os.RemoveAll(us.tempDir)
}

func (us *UtilsSuite) TestGetEnvOrDefault() {

	// test with existing env variable
	result := utils.GetEnvOrDefault(ENV_VAR_NAME, "DUMMY_ENVVAR_VALUE")
	us.Equal(ENV_VAR_VALUE, result)

	// test with non-existent env variable
	result = utils.GetEnvOrDefault("TEST_ENVVAR2", "DUMMY_ENVVAR_VALUE")
	us.Equal("DUMMY_ENVVAR_VALUE", result)

}

func (us *UtilsSuite) TestReadFileOrEmpty() {

	tempfile := filepath.Join(us.tempDir, TESTREADFILENAME)
	err := os.WriteFile(tempfile, []byte("Dummy text\n"), 0664)
	us.NoError(err)

	// test with existing file
	content := utils.ReadFileOrEmpty(tempfile)
	us.Equal("Dummy text\n", content)

	// test with non-existing file
	content = utils.ReadFileOrEmpty(fmt.Sprintf("123%s", TESTREADFILENAME))
	us.Equal("", content)

}

func (us *UtilsSuite) TestNormalizeDirectoryPath() {

	// test with trailing "/"
	result := utils.NormalizeDirectoryPath("/1/2/3/")
	us.Equal("/1/2/3", result)

	// test without trailing "/"
	result = utils.NormalizeDirectoryPath("/1/2/3")
	us.Equal("/1/2/3", result)

	// test with empty string
	result = utils.NormalizeDirectoryPath("")
	us.Equal("", result)

}

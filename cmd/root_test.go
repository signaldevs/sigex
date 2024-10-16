package cmd_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/signaldevs/sigex/cmd"
	sigex "github.com/signaldevs/sigex/pkg"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RootCmdTestSuite struct {
	suite.Suite
	command  *cobra.Command
	osHelper *osHelperMock
}

type osHelperMock struct {
	mock.Mock
}

func (o *osHelperMock) LookPath(path string) (string, error) {
	args := o.Called(path)
	return args.String(0), args.Error(1)
}

func (o *osHelperMock) Exec(argv0 string, argv []string, envv []string) error {
	args := o.Called(argv0, argv, envv)
	return args.Error(0)
}

// Setup test suite
func (suite *RootCmdTestSuite) SetupTest() {
	suite.command = &cobra.Command{Use: "sigex", RunE: cmd.RootCmdRunE}
	suite.osHelper = new(osHelperMock)
	sigex.SetOSHelper(suite.osHelper)
}

func (suite *RootCmdTestSuite) TestNoArguments() {
	_, err := execute(suite.T(), suite.command, []string{}...)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), errors.New("no command argument was provided"), err)
}

func (suite *RootCmdTestSuite) TestRootCmdEchoCommand() {
	suite.osHelper.On("LookPath", mock.Anything).Return("/bin/echo", nil)
	suite.osHelper.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	args := []string{"echo", "hello"}
	_, err := execute(suite.T(), suite.command, args...)
	assert.Nil(suite.T(), err)
}

// Helper execute function
func execute(t *testing.T, c *cobra.Command, args ...string) (string, error) {
	t.Helper()

	buf := new(bytes.Buffer)
	c.SetOut(buf)
	c.SetErr(buf)
	c.SetArgs(args)

	err := c.Execute()
	return strings.TrimSpace(buf.String()), err
}

func TestRootCmd(t *testing.T) {
	suite.Run(t, new(RootCmdTestSuite))
}

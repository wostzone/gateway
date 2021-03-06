package hub_test

import (
	"flag"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wostzone/hub/core/hub"
	"github.com/wostzone/wostlib-go/pkg/certsetup"
)

// testing takes place using the test folder on localhost
var homeFolder string
var certsFolder string

var hostnames = []string{"localhost"}

// TestMain sets the project test folder as the home folder and makes sure the neccesary
// certificates exist.
func TestMain(m *testing.M) {
	cwd, _ := os.Getwd()
	homeFolder = path.Join(cwd, "../../test")
	certsFolder = path.Join(homeFolder, "certs")
	certsetup.CreateCertificateBundle(hostnames, certsFolder)

	result := m.Run()

	os.Exit(result)
}

// Setup to run each test
func setup() {
	// Reset args to prevent 'flag redefined' error
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	os.Args = append(os.Args[0:1], strings.Split("", " ")...)

	// hubConfig, _ = hubconfig.SetupConfig(homeFolder, pluginID, customConfig)
}

func TestStartHubNoPlugins(t *testing.T) {
	setup()
	err := hub.StartHub(homeFolder, false)
	require.NoError(t, err)

	time.Sleep(3 * time.Second)
	hub.StopHub()
}

func TestStartHubWithPlugins(t *testing.T) {
	setup()
	err := hub.StartHub(homeFolder, true)
	assert.NoError(t, err)

	time.Sleep(3 * time.Second)
	hub.StopHub()
}

func TestStartHubBadHome(t *testing.T) {
	setup()
	err := hub.StartHub("/notahomefolder", true)
	assert.Error(t, err)

}

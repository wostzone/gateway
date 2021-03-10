package config

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// HubConfigName the configuration file name of the hub
const HubConfigName = "hub.yaml"

// HubLogFile the file name of the hub logging
const HubLogFile = "hub.log"

// DefaultCertsFolder with the location of certificates
const DefaultCertsFolder = "./certs"

// DefaultSmbHost with the default server address and port
const DefaultSmbHost = "localhost:9678"

// ConfigArgs configuration commandline arguments
type ConfigArgs struct {
	name         string
	defaultValue string
	description  string
}

// HubConfig with hub configuration parameters
type HubConfig struct {
	Logging struct {
		Loglevel   string `yaml:"logLevel"`   // debug, info, warning, error. Default is warning
		LogFile    string `yaml:"logFile"`    // hub logging to file
		TimeFormat string `yaml:"timeFormat"` // go default ISO8601 ("2006-01-02T15:04:05.000-0700")
	} `yaml:"logging"`

	// Messenger configuration of hub plugin messaging
	Messenger struct {
		CertFolder string `yaml:"certFolder"` // location of certificates when using TLS. Default is ./certs
		HostPort   string `yaml:"hostname"`   // hostname:port or ip:port to listen on of message bus
		Protocol   string `yaml:"protocol"`   // internal, MQTT, default internal
		Timeout    int    `yaml:"timeout"`    // Client connection timeout in seconds. 0 for indefinite
	} `yaml:"messenger"`

	Home         string   `yaml:"home"`         // application home directory. Default is parent of executable.
	ConfigFolder string   `yaml:"configFolder"` // location of configuration files. Default is ./config
	PluginFolder string   `yaml:"pluginFolder"` // location of plugin binaries. Default is ./bin
	Plugins      []string `yaml:"plugins"`      // names of plugins to start
	// internal
}

// CreateDefaultHubConfig with default values
// homeFolder is the home of the application, log and configuration folders.
// Use "" for default: parent of application binary
// When relative path is given, it is relative to the application binary
func CreateDefaultHubConfig(homeFolder string) *HubConfig {
	appBin, _ := os.Executable()
	binFolder := path.Dir(appBin)
	if homeFolder == "" {
		homeFolder = path.Dir(binFolder)

		// for running within the project tests use the test folder as application root folder
		// if path.Base(binFolder) != "bin" {
		// 	// appFolder = path.Join(appFolder, "../test")
		// 	cwd, _ := os.Getwd()
		// 	appFolder = path.Join(cwd, "../../test")
		// }
		// logrus.Infof("appBin: %s. CWD=%s", appBin, cwd)
	} else if !path.IsAbs(homeFolder) {
		// turn relative home folder in absolute path
		// cwd, _ := os.Getwd()
		// homeFolder = path.Join(cwd, homeFolder)
		homeFolder = path.Join(binFolder, homeFolder)
	}
	logrus.Infof("AppBin is: %s; Home is: %s", appBin, homeFolder)
	config := &HubConfig{
		// ConfigFolder: path.Join(homeFolder, "config"),
		Home:         homeFolder,
		ConfigFolder: path.Join(homeFolder, "config"),
		Plugins:      make([]string, 0),
		PluginFolder: path.Join(homeFolder, "./bin"),
	}
	// config.Messenger.CertsFolder = path.Join(homeFolder, "certs")
	config.Messenger.CertFolder = path.Join(homeFolder, DefaultCertsFolder)
	config.Messenger.HostPort = DefaultSmbHost // use default "localhost:9678"
	config.Messenger.Protocol = ""             // use default
	config.Logging.Loglevel = "warning"
	// config.Logging.LogFile = path.Join(homeFolder, "logs/"+HubLogFile)
	config.Logging.LogFile = path.Join(homeFolder, "./logs/"+HubLogFile)
	return config
}

// LoadConfig loads the configuration from file into the given config
// Returns nil if successful
func LoadConfig(configFile string, config interface{}) error {
	var err error
	var rawConfig []byte
	rawConfig, err = ioutil.ReadFile(configFile)
	if err != nil {
		logrus.Warningf("Unable to load config file: %s", err)
		return err
	}
	logrus.Infof("Loaded config file '%s'", configFile)

	err = yaml.Unmarshal(rawConfig, config)
	if err != nil {
		logrus.Errorf("Error parsing config file '%s': %s", configFile, err)
		return err
	}
	return nil
}

// ValidateConfig checks if values in the hub configuration are correct
// Returns an error if the config is invalid
func ValidateConfig(config *HubConfig) error {
	if _, err := os.Stat(config.Home); os.IsNotExist(err) {
		logrus.Errorf("Home folder '%s' not found\n", config.Home)
		return err
	}
	if _, err := os.Stat(config.ConfigFolder); os.IsNotExist(err) {
		logrus.Errorf("Configuration folder '%s' not found\n", config.ConfigFolder)
		return err
	}

	loggingFolder := path.Dir(config.Logging.LogFile)
	if _, err := os.Stat(loggingFolder); os.IsNotExist(err) {
		logrus.Errorf("Logging folder '%s' not found\n", loggingFolder)
		return err
	}

	if _, err := os.Stat(config.Messenger.CertFolder); os.IsNotExist(err) {
		logrus.Errorf("TLS certificate folder '%s' not found\n", config.Messenger.CertFolder)
		return err
	}
	// Pluginfolder is either empty or must exist
	if config.PluginFolder != "" {
		if _, err := os.Stat(config.PluginFolder); os.IsNotExist(err) {
			logrus.Errorf("Plugins folder '%s' not found\n", config.PluginFolder)
			return err
		}
	}

	return nil
}
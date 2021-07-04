package idprovpb

import (
	"github.com/wostzone/idprov-go/pkg/idprov"
	"github.com/wostzone/idprov-go/pkg/idprovserver"
	"github.com/wostzone/wostlib-go/pkg/hubconfig"
)

const PluginID = "idprov-pb"

// IDProvPBConfig Protocol binding configuration
type IDProvPBConfig struct {
	IdpAddress      string `yaml:"idpAddress"`      // listening address, default is automatic
	IdpPort         uint   `yaml:"idpPort"`         // idprov listening port
	IdpCerts        string `yaml:"idpCerts"`        // folder to store client certificates
	ClientID        string `yaml:"clientID"`        // unique service instance client ID
	EnableDiscovery bool   `yaml:"enableDiscovery"` // DNS-SD disco
	ValidForDays    uint   `yaml:"validForDays"`    // Nr days certificates are valid for
}

// IDProv provisioning protocol binding service
// Configure and start IDProv server
type IDProvPB struct {
	config    IDProvPBConfig
	hubConfig hubconfig.HubConfig
	server    *idprovserver.IDProvServer
}

// Start the IDProv service
func (pb *IDProvPB) Start() error {

	err := pb.server.Start()
	return err
}

// Stop the IDProv service
func (pb *IDProvPB) Stop() {
	if pb.server != nil {
		pb.server.Stop()
	}
}

func NewIDProvPB(config IDProvPBConfig, hubConfig hubconfig.HubConfig) *IDProvPB {
	// use default values if config is incomplete
	if config.IdpAddress == "" {
		config.IdpAddress = idprovserver.GetOutboundIP("").String()
	}
	if config.IdpPort == 0 {
		config.IdpPort = 43776
	}
	if config.IdpCerts == "" {
		config.IdpCerts = "./idpcerts"
	}
	if config.ValidForDays <= 0 {
		config.ValidForDays = 3
	}
	server := idprovserver.NewIDProvServer(
		config.ClientID,
		config.IdpAddress,
		config.IdpPort,
		hubConfig.CertsFolder, config.IdpCerts,
		config.ValidForDays, idprov.IdprovServiceDiscoveryType)

	return &IDProvPB{
		config:    config,
		hubConfig: hubConfig,
		server:    server,
	}
}
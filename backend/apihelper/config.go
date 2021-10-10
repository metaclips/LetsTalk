package apihelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/pion/webrtc/v3"
)

type configuration struct {
	TLS struct {
		CertPath string
		KeyPath  string
	}

	ICEServers []struct {
		URLs       []string
		Username   string
		AuthType   string
		AuthSecret string
	}

	peerConnectionICEServers []webrtc.ICEServer

	CorsAllowedOrigins []string

	AdminEmail               string
	AdminPassword            string
	Email                    string
	EmailPassword            string
	DbHost                   string
	DbName                   string
	DropboxToken             string
	Port                     string
	EnableClassSessionRecord bool // If EnableClassSessionUpload is set to true and no token is provided, files are saved to DB using GridFS.
}

// Config contains application environment variables.
var Config configuration

var credetialType map[string]webrtc.ICECredentialType = map[string]webrtc.ICECredentialType{
	"Password": webrtc.ICECredentialTypePassword,
	"Oauth":    webrtc.ICECredentialTypePassword,
}

// LoadConfiguration loads all application environment variables.
func LoadConfiguration(configPath string) error {
	file, err := os.Open(configPath) // For read access.
	if err != nil {
		return err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		return err
	}

	if len(Config.ICEServers) == 0 {
		Config.peerConnectionICEServers = []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
		}
	}

	for _, config := range Config.ICEServers {
		if len(config.URLs) == 0 {
			return errors.New("user did not specify ICE server")
		}

		credential, ok := credetialType[config.AuthType]
		if !ok {
			return fmt.Errorf("invalid webrtc credential type %s only AuthType Password and Oauth are allowed", config.AuthType)
		}

		iceServer := webrtc.ICEServer{
			URLs:           config.URLs,
			Username:       config.Username,
			Credential:     config.AuthSecret,
			CredentialType: credential,
		}

		Config.peerConnectionICEServers = append(Config.peerConnectionICEServers, iceServer)
	}

	return nil
}

func (c *configuration) GetPeerConfig() webrtc.Configuration {
	return webrtc.Configuration{
		SDPSemantics: webrtc.SDPSemanticsUnifiedPlanWithFallback,
		ICEServers:   c.peerConnectionICEServers,
	}
}

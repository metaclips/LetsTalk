package values

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	err := LoadConfiguration("../../config.json")
	if err != nil {
		t.Error("could not load config", err)
	}

	assert.Equal(t, "mongodb://localhost:27017", Config.DbHost)
	assert.Equal(t, "LetsTalkDB", Config.DbName)
	assert.Equal(t, true, Config.EnableClassSessionRecord)
	assert.Equal(t, len(Config.ICEServers), 1)
	assert.Equal(t, []string{"stun:stun.l.google.com:19302"}, Config.ICEServers[0].URLs)
	assert.Equal(t, "cert.pem", Config.TLS.CertPath)
	assert.Equal(t, "key.pem", Config.TLS.KeyPath)
}

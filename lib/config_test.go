package lib

import (
	"github.com/alexcesaro/log"
	"github.com/alexcesaro/log/golog"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

var configReader *ConfigReader

func init() {
	logger := *golog.New(os.Stderr, log.Warning)
	configReader = NewConfigReader(logger)
}

func TestConfig_NotificationEmailYml(t *testing.T) {
	setup()
	defer cleanup()

	yamlStr := "" +
		"notifications:\n" +
		"   email: test@email.de"
	ioutil.WriteFile(".autoupgrade.yml", []byte(yamlStr), 0700)

	config, err := configReader.ReadConfig()

	assert.Nil(t, err)
	assert.Equal(t, config.Notification.Email, "test@email.de")
}

func TestConfig_NotificationEmailYaml(t *testing.T) {
	setup()
	defer cleanup()

	yamlStr := "" +
		"notifications:\n" +
		"   email: test@email.de"
	ioutil.WriteFile(".autoupgrade.yaml", []byte(yamlStr), 0700)

	config, err := configReader.ReadConfig()

	assert.Nil(t, err)
	assert.Equal(t, config.Notification.Email, "test@email.de")
}

func TestConfig_MissingConfigFile(t *testing.T) {
	setup()
	defer cleanup()

	config, err := configReader.ReadConfig()

	assert.Nil(t, err)
	assert.Equal(t, config, &Config{})
}

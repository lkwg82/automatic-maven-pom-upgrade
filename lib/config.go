package lib

import (
	"github.com/alexcesaro/log/golog"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
)

type Config struct {
	Notification ConfigNotification `json:"notifications"`
}

type ConfigNotification struct {
	Email string `json:"email"`
}

type ConfigReader struct {
	logger golog.Logger
}

func NewConfigReader(logger golog.Logger) *ConfigReader {
	return &ConfigReader{
		logger: logger,
	}
}

func (cr *ConfigReader) ReadConfig() (*Config, error) {
	configFile := ".autoupgrade.yml"
	cr.logger.Debugf("looking for config file '%s'", configFile)
	if _, err := os.Stat(configFile); err != nil {
		configFile = ".autoupgrade.yaml"
		cr.logger.Debugf("looking for config file '%s'", configFile)
		if _, err := os.Stat(configFile); err != nil {
			cr.logger.Debug("no config file found")
			return &Config{}, nil
		}
	}

	cr.logger.Debugf("reading config from %s", configFile)

	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

package config

import (
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Config which would be read from the config.yaml
type Config struct {
	SyncInterval    string       `yaml:"syncInterval"`
	RemoveUnknownIp bool         `yaml:"removeUnknownIp"`
	IpProviders     []IpProvider `yaml:"ipProviders"`
	Provider        Provider     `yaml:"provider"`
	Filter          Filter       `yaml:"filter"`
}

// IpProvider that the controller will be using to gather whitelist IPs
type IpProvider struct {
	Name   string                      `yaml:"name"`
	Params map[interface{}]interface{} `yaml:"params"`
}

// Provider that the controller will be using to update to allow access
type Provider struct {
	Name   string                      `yaml:"name"`
	Params map[interface{}]interface{} `yaml:"params"`
}

// Filter that will be used to filter resources on the provider
type Filter struct {
	LabelName  string `yaml:"labelName"`
	LabelValue string `yaml:"labelValue"`
}

// ReadConfig function that reads the yaml file
func ReadConfig(filePath string) (Config, error) {
	var config Config
	// Read YML
	source, err := ioutil.ReadFile(filePath)
	if err != nil {
		return config, err
	}

	// Unmarshall
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

// GetConfiguration gets the yaml configuration for the controller
func GetConfiguration() Config {
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if len(configFilePath) == 0 {
		//Default config file is placed in configs/ folder
		configFilePath = "configs/config.yaml"
	}
	configuration, err := ReadConfig(configFilePath)
	if err != nil {
		log.Panic(err)
	}
	return configuration
}

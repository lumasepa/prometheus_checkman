package main

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Check struct {
	Name      string `yaml:"Name"`
	Command   string `yaml:"Command"`
	Frequency int `yaml:"Frequency"`
	Environment map[string]string `yaml:"Environment"`
	Labels map[string]string `yaml:"Labels"`
	Help string `yaml:"Help"`
}

type Configuration struct {
	Checks []Check `yaml:"checks"`
	ListenIP string `yaml:"listen_ip"`
	ListenPort int `yaml:"listen_port"`
	ExporterPath string `yaml:"exporter_path"`
}


func ParseConf(filePath string) (*Configuration, error) {
	confFileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var configuration *Configuration
	err = yaml.Unmarshal(confFileContent, &configuration)

	if err != nil {
		return nil, err
	}

	return configuration, nil
}

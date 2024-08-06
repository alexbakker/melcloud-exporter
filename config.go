package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type MELCloudConfig struct {
	Email           string `yaml:"email"`
	Password        string `yaml:"password"`
	RefreshInterval int    `yaml:"refresh_interval"`
}

type PrometheusConfig struct {
	Addr string `yaml:"addr"`
}

type Config struct {
	Prometheus *PrometheusConfig `yaml:"prometheus"`
	MELCloud   *MELCloudConfig   `yaml:"melcloud"`
}

func ReadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err = yaml.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

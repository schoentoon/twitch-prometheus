package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ClientID     string   `yaml:"client_id"`
	ClientSecret string   `yaml:"client_secret"`
	ListenPort   int      `yaml:"port"`
	Followers    []string `yaml:"followers"`
}

func ReadConfig(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	out := &Config{}
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&out)
	if err != nil {
		return nil, err
	}

	return out, err
}

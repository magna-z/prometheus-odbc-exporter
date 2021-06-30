package config

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type MetricsOptions struct {
	Usage       string `yaml:"usage"`
	Description string `yaml:"description"`
}

type Metric struct {
	Query    string           `yaml:"query"`
	CacheTTL time.Duration    `yaml:"cache_ttl"`
	Metrics  []MetricsOptions `yaml:"metrics"`
}

type Config map[string]Metric

func GetConfig(path string) (Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

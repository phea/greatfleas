package core

import (
	"encoding/json"
	"io/ioutil"
)

// Config holds the configuration option.
type Config struct {
	ArticlePath  string `json:"article_path"`
	TemplatePath string `json:"template_path"`
	PublicPath   string `json:"public_path"`
	Port         string `json:"port"`
}

func loadConfig(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var conf Config
	if err := json.Unmarshal(data, &conf); err != nil {
		return nil, err
	}

	return &conf, err
}

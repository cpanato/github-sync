package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Config is the information needed to communicate with Mattermost.
type Config struct {
	AuthToken string `json:"authToken"`
	Org       string `json:"org"`
}

// LoadConfig loads a Config from a JSON file.
func LoadConfig(path string) (*Config, error) {
	config := &Config{}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("couldn't open file: %v", err)
	}
	if err := json.Unmarshal(content, &config); err != nil {
		return config, fmt.Errorf("couldn't parse config: %v", err)
	}

	return config, nil
}

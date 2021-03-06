package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	valid "monitor/utils/validators"
)

// Config represents a configuration set for the program to run
type Config struct {
	CheckTimer struct {
		Interval int
		Timeout  int
	}
	URLs []URL
}

// CreateConfigurationFromFile returns a new configuration loaded from a file
func CreateConfigurationFromFile(configFile string) (Config, error) {
	config := Config{
		CheckTimer: struct {
			Interval int
			Timeout  int
		}{
			Interval: 30,
			Timeout:  5,
		},
		URLs: []URL{},
	}

	if _, err := os.Stat(configFile); !os.IsNotExist(err) {
		log.Printf("Found configuration file: %s\n", configFile)

		configData, err := ioutil.ReadFile(configFile)
		if err != nil {
			return config, fmt.Errorf("Error reading configuation file: %v", err)
		}

		err = json.Unmarshal(configData, &config)
		if err != nil {
			return config, fmt.Errorf("Error decoding configuration: %v", err)
		}

		log.Print("Configuration loaded successfully\n")
	}

	if len(config.URLs) == 0 {
		return config, fmt.Errorf("No urls found in the configuration")
	}

	for _, url := range config.URLs {
		if isValid, _ := valid.IsValidURL(url.URL); !isValid {
			return config, fmt.Errorf("Invalid url format in URLs section of config: %v", url.URL)
		}
	}

	maxTotalTimeout := config.CheckTimer.Timeout * len(config.URLs)

	if maxTotalTimeout >= config.CheckTimer.Interval {
		return config, fmt.Errorf("Timeout value (times monitor count) cannot be greater or equal than interval in the configuration file")
	}

	return config, nil
}

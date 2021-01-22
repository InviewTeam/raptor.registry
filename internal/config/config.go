package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type RabbitConfig struct {
	Address       string `json:"address"`
	WorkerQueue   string `json:"worker_queue"`
	AnalyzerQueue string `json:"analyzer_queue"`
}

type DatabaseConfig struct {
	Address    string `json:"address"`
	Database   string `json:"database"`
	Collection string `json:"collection"`
}

type Settings struct {
	RegistryAddress string         `json:"registry_address"`
	Database        DatabaseConfig `json:"database"`
	Rabbit          RabbitConfig   `json:"rabbit"`
}

func Init(cfgFile string) (*Settings, error) {
	conf := &Settings{}

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		pwd, _ := os.Getwd()
		viper.SetConfigName("configs/example_config")
		viper.AddConfigPath(pwd)
		viper.AutomaticEnv()
		viper.SetConfigType("json")
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("cannot read config file %s: %w", cfgFile, err)
	}

	if err := viper.Unmarshal(conf); err != nil {
		return nil, fmt.Errorf("cannot parse config file %s: %w", cfgFile, err)
	}

	return conf, nil
}

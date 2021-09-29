package settings

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Settings map[string]interface{}

var (
	settings Settings
)

const (
	EnvPrefix              = "j3"
	DefaultConfigFileName  = "application"
	DefaultConfigLocations = ".,./config"
)

// Parse and get settings
func Parse() {
	// Environment management
	viper.AutomaticEnv()
	viper.SetEnvPrefix(EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	// default config file name
	viper.SetConfigName(DefaultConfigFileName)

	// default config paths.
	for _, path := range strings.Split(DefaultConfigLocations, ",") {
		viper.AddConfigPath(path)
	}

	// Check if we have a specific file to load, and if so set it.
	pflag.StringP("configuration-file", "c", "", "path to the configuration file")
	pflag.Parse()
	configurationFilePath, _ := pflag.CommandLine.GetString("configuration-file")
	if configurationFilePath != "" {
		viper.SetConfigFile(configurationFilePath)
	}

	// read configuration
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(fmt.Errorf("unable to load settings: %s", err))
		}
	}

	// convert to settings and return the result.
	settings = viper.AllSettings()
}

// Parse and get settings
func ParseAndGet() Settings {
	Parse()
	return Get()
}

// Returns the current settings.
func Get() Settings {
	return settings
}

// Return the key of the settings
func GetKey(key string) map[string]interface{} {
	if settings == nil || settings[key] == nil {
		return make(map[string]interface{})
	}
	return settings[key].(map[string]interface{})
}

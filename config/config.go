package config

import (
	"cloud-commis/logger"
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var ParsedData = koanf.New(".")

const (
	defaultConfigPath = "/etc/cloud-commis/config.yaml"
)

func Read(config *koanf.Koanf) {

	Flags.Read()

	// default config
	err := config.Load(confmap.Provider(map[string]interface{}{
		"configPath":       defaultConfigPath,
		"loglevel":         "info",
		"scanIntervalMin":  60,
		"disable_ui":       false,
		"storage":          "local",
		"localStoragePath": "/data/cloud-commis",
		"s3BucketName":     "",
		"s3BucketPath":     "",
		"httpPort":         8080,
	}, "."), nil)
	if err != nil {
		logger.Log.Error(err.Error())
	}

	// default config is overwritten by yaml config
	if _, err := os.Stat(Flags.configFile); err == nil {
		if err := config.Load(file.Provider(Flags.configFile), yaml.Parser()); err != nil {
			logger.Log.Error(err.Error())
		} else {
			err := config.Set("configPath", Flags.configFile)
			if err != nil {
				logger.Log.Error(err.Error())
			}
		}
	} else {
		logger.Log.Info("No config file found")
	}

	// yaml config is overwrittent by env
	err = config.Load(env.Provider("CC_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "CC_")), "_", ".", -1)
	}), nil)
	if err != nil {
		logger.Log.Error(err.Error())
	}

}

package config

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
)

// build infos from compiler args
var Version string
var BuildDate string

type flagConfig struct {
	configFile string
}

var Flags flagConfig

func (flagConfig *flagConfig) Read() {
	configFile := flag.String("config", defaultConfigPath, "config file path")
	version := flag.Bool("version", false, "show cloud commis version")

	flag.Parse()

	Flags.configFile = *configFile

	if *version {
		if Version == "" {
			Version = "devel"
		}

		fmt.Println("cloud-commis " + Version)
		fmt.Println("Build \t" + BuildDate)

		buildInfo, _ := debug.ReadBuildInfo()

		fmt.Println("GOLANG \t" + buildInfo.GoVersion)

		for setting := range buildInfo.Settings {
			if strings.HasPrefix(buildInfo.Settings[setting].Key, "GOARCH") ||
				strings.HasPrefix(buildInfo.Settings[setting].Key, "GOOS") {
				fmt.Println(buildInfo.Settings[setting].Key + " \t" + buildInfo.Settings[setting].Value)
			}

		}

		os.Exit(0)
	}
}

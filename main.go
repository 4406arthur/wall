package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	core "wall/cmd"
	"wall/utils/logger"

	"github.com/spf13/viper"
)

var usageStr = `
Developer Manager BK Service

Server Options:
    -c, --config <file>              Configuration file path
    -h, --help                       Show this message
    -v, --version                    Show version
`

// usage will print out the flag options for the server.
func usage() {
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}

func setup(path string) *viper.Viper {
	v := viper.New()
	v.SetConfigType("json")
	v.SetConfigName("config")
	if path != "" {
		v.AddConfigPath(path)
	} else {
		v.AddConfigPath("./config/")
	}

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	return v
}

var version string

func printVersion() {

	fmt.Printf(`JWT token issuer service %s, Compiler: %s %s, Copyright (C) 2021 E.SUN BANK, Inc.`,
		version,
		runtime.Compiler,
		runtime.Version())
	fmt.Println()
}

func main() {
	//from env
	var configFile string
	var showVersion bool
	version = "1.0.0"
	flag.BoolVar(&showVersion, "v", false, "Print version information.")
	flag.StringVar(&configFile, "c", "", "Configuration file path.")
	flag.Usage = usage
	flag.Parse()

	if showVersion {
		printVersion()
		os.Exit(0)
	}
	viperConfig := setup(configFile)
	log := logger.InitLogger(viperConfig.GetInt("server_config.log_level"))

	core.InitRouter(log, viperConfig)
}

package main

import (
	"flag"
	"log"
	"project/inetrnal/app/projectname"
)

var (
	configPath    string
	profileString string
)

func init() {
	flag.StringVar(&configPath, "configs-path", "./configs/config.toml", "Path to config files")
	flag.StringVar(&profileString, "profile", "local", "Profile where app runs, accept two options docker/local")
}

func main() {
	flag.Parse()

	config, err := projectname.ParseConfig(configPath, profileString)
	if err != nil {
		log.Fatal(err)
	}

	if err := projectname.Start(config); err != nil {
		log.Fatal(err)
	}
}

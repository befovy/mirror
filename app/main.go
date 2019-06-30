package main

import (
	"flag"
	"github.com/befovy/mirror"
	_ "github.com/befovy/mirror/issues"
	//"github.com/gohugoio/hugo/commands"
)

func main() {
	var config string
	flag.StringVar(&config, "config", "./mirror.yaml", "the file path of mirror config")
	flag.Parse()
	mirror.SyncWithConfig(config)
	//commands.Execute()
}

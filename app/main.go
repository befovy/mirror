package main

import (
	"github.com/befovy/mirror"
	_ "github.com/befovy/mirror/issues"
	//"github.com/gohugoio/hugo/commands"
)

func main() {
	mirror.Sync()
	//commands.Execute()
}

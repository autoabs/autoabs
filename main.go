package main

import (
	"flag"
	"github.com/autoabs/autoabs/cmd"
	"github.com/autoabs/autoabs/requires"
)

func main() {
	flag.Parse()

	requires.Init()

	switch flag.Arg(0) {
	case "app":
		cmd.App()
	case "set":
		cmd.Settings()
	case "sync":
		cmd.Sync()
	case "build":
		cmd.Build()
	}
}

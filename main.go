package main

import (
	"flag"
	"github.com/autoabs/autoabs/cmd"
	"github.com/autoabs/autoabs/config"
	"github.com/autoabs/autoabs/requires"
)

func main() {
	flag.Parse()

	err := config.Load()
	if err != nil {
		panic(err)
	}

	requires.Init()

	switch flag.Arg(0) {
	case "app":
		cmd.App()
	case "set":
		cmd.Settings()
	}
}

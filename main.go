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
	case "builds":
		switch flag.Arg(1) {
		case "sync":
			cmd.Sync()
		case "build":
			cmd.Build()
		case "clear":
			switch flag.Arg(2) {
			case "all":
				cmd.ClearAll()
			case "pending":
				cmd.ClearPending()
			case "failed":
				cmd.ClearFailed()
			}
		}
	case "genkey":
		cmd.GenKey()
	}
}

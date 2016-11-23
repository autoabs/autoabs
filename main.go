package main

import (
	"flag"
	"fmt"
	"github.com/autoabs/autoabs/cmd"
	"github.com/autoabs/autoabs/requires"
	"os"
)

func main() {
	flag.Parse()

	requires.Init()

	switch flag.Arg(0) {
	case "app":
		cmd.App()
		return
	case "set":
		cmd.Settings()
		return
	case "builds":
		switch flag.Arg(1) {
		case "sync":
			cmd.Sync()
			return
		case "build":
			cmd.Build()
			return
		case "clear":
			switch flag.Arg(2) {
			case "all":
				cmd.ClearAll()
				return
			case "pending":
				cmd.ClearPending()
				return
			case "failed":
				cmd.ClearFailed()
				return
			}
		}
	case "genkey":
		cmd.GenKey()
		return
	}

	fmt.Fprintln(os.Stderr, "Unknown command")
}

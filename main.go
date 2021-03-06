package main

import (
	"flag"
	"fmt"
	"github.com/autoabs/autoabs/cmd"
	"github.com/autoabs/autoabs/logger"
	"github.com/autoabs/autoabs/requires"
	"os"
)

func main() {
	flag.Parse()

	requires.Init()
	logger.Init()

	switch flag.Arg(0) {
	case "set":
		cmd.Settings()
		return
	case "node":
		switch flag.Arg(1) {
		case "web":
			cmd.WebNode()
			return
		case "builder":
			cmd.BuilderNode()
			return
		case "storage":
			cmd.StorageNode()
			return
		}
	case "builds":
		switch flag.Arg(1) {
		case "sync":
			cmd.Sync()
			return
		case "sync-state":
			cmd.SyncState()
			return
		case "upload":
			cmd.Upload()
			return
		case "clean":
			cmd.Clean()
			return
		case "build":
			cmd.Build()
			return
		case "retry":
			cmd.RetryFailed()
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

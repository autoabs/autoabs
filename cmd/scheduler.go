package cmd

import (
	"github.com/autoabs/autoabs/node"
	"github.com/autoabs/autoabs/scheduler"
	"github.com/autoabs/autoabs/utils"
)

func StorageScheduler() {
	nde := node.Node{
		Id:   utils.RandName(),
		Type: "storage",
	}
	nde.Keepalive()

	sch := scheduler.Storage{}

	sch.Start()
}

func BuildScheduler() {
	nde := node.Node{
		Id:   utils.RandName(),
		Type: "builder",
	}
	nde.Keepalive()

	sch := scheduler.Build{}

	sch.Start()
}

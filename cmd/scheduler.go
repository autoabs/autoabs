package cmd

import (
	"github.com/autoabs/autoabs/scheduler"
)

func StorageScheduler() {
	sch := scheduler.Storage{}

	sch.Start()
}

func BuildScheduler() {
	sch := scheduler.Build{}

	sch.Start()
}

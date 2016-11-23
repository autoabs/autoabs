package cmd

import (
	"github.com/autoabs/autoabs/build"
	"github.com/autoabs/autoabs/queue"
)

func Sync() {
	que := queue.Queue{}

	err := que.Queue()
	if err != nil {
		panic(err)
	}
}

func Build() {
	que := queue.Queue{}

	err := que.Build()
	if err != nil {
		panic(err)
	}
}

func ClearAll() {
	err := build.ClearAllBuilds()
	if err != nil {
		panic(err)
	}
}

func ClearPending() {
	err := build.ClearPendingBuilds()
	if err != nil {
		panic(err)
	}
}

func ClearFailed() {
	err := build.ClearFailedBuilds()
	if err != nil {
		panic(err)
	}
}

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

func Upload() {
	que := queue.Queue{}

	err := que.Upload()
	if err != nil {
		panic(err)
	}
}

func Clean() {
	que := queue.Queue{}

	err := que.Clean()
	if err != nil {
		panic(err)
	}
}

func RetryFailed() {
	err := build.RetryFailed()
	if err != nil {
		panic(err)
	}
}

func ClearAll() {
	err := build.ClearAll()
	if err != nil {
		panic(err)
	}
}

func ClearPending() {
	err := build.ClearPending()
	if err != nil {
		panic(err)
	}
}

func ClearFailed() {
	err := build.ClearFailed()
	if err != nil {
		panic(err)
	}
}

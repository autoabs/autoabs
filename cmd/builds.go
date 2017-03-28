package cmd

import (
	"github.com/autoabs/autoabs/build"
	"github.com/autoabs/autoabs/builder"
	"github.com/autoabs/autoabs/queue"
)

func Sync() {
	que := queue.Queue{}

	err := que.Sync()
	if err != nil {
		panic(err)
	}
}

func SyncState() {
	que := queue.Queue{}

	err := que.SyncState()
	if err != nil {
		panic(err)
	}
}

func Build() {
	bilder := builder.Builder{}

	err := bilder.Start()
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

package cmd

import (
	"github.com/autoabs/autoabs/queue"
)

func Sync() {
	que := queue.Queue{}

	err := que.Queue()
	if err != nil {
		panic(err)
	}
}

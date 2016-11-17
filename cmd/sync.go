package cmd

import (
	"github.com/autoabs/autoabs/queue"
)

func Sync() {
	que := queue.Queue{}

	err := que.Scan()
	if err != nil {
		panic(err)
	}

	err = que.Queue()
	if err != nil {
		panic(err)
	}
}

package cmd

import (
	"github.com/autoabs/autoabs/queue"
)

func Build() {
	que := queue.Queue{}

	err := que.Build()
	if err != nil {
		panic(err)
	}
}

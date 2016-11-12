package main

import (
	"github.com/autoabs/autoabs/config"
	"github.com/autoabs/autoabs/queue"
)

func main() {
	err := config.Load()
	if err != nil {
		panic(err)
	}

	que := &queue.Queue{}

	err = que.Queue()
	if err != nil {
		panic(err)
	}

	err = que.Build()
	if err != nil {
		panic(err)
	}
}

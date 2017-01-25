package scheduler

import (
	"github.com/Sirupsen/logrus"
	"github.com/autoabs/autoabs/queue"
	"time"
)

type Build struct{}

func (b *Build) build() (err error) {
	logrus.Info("scheduler: Building")

	que := queue.Queue{}

	err = que.Build()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("scheduler: Building failed")
		return
	}

	return
}

func (b *Build) Start() {
	for {
		b.build()
		time.Sleep(1 * time.Second)
	}
}

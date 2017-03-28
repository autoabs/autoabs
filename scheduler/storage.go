package scheduler

import (
	"github.com/Sirupsen/logrus"
	"github.com/autoabs/autoabs/queue"
	"time"
)

type Storage struct{}

func (s *Storage) sync() (err error) {
	logrus.Info("scheduler: Syncing")

	que := &queue.Queue{}

	err = que.Sync()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("scheduler: Syncing failed")
		return
	}

	return
}

func (s *Storage) syncState() (err error) {
	logrus.Info("scheduler: Syncing state")

	que := &queue.Queue{}

	err = que.SyncState()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("scheduler: Syncing state failed")
		return
	}

	return
}

func (s *Storage) upload() (err error) {
	logrus.Info("scheduler: Uploading")

	que := &queue.Queue{}

	err = que.Upload()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("scheduler: Uploading failed")
		return
	}

	return
}

func (s *Storage) clean() (err error) {
	logrus.Info("scheduler: Cleaning")

	que := &queue.Queue{}

	err = que.Clean()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("scheduler: Cleaning failed")
		return
	}

	return
}

func (s *Storage) runAll() {
	s.sync()
	time.Sleep(1 * time.Second)
	s.syncState()
	time.Sleep(1 * time.Second)
	s.upload()
	time.Sleep(1 * time.Second)
	s.clean()
}

func (s *Storage) Start() {
	for {
		s.runAll()
		time.Sleep(1 * time.Second)
	}
}

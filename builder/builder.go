package builder

import (
	"github.com/Sirupsen/logrus"
	"github.com/autoabs/autoabs/build"
	"github.com/autoabs/autoabs/database"
	"github.com/autoabs/autoabs/node"
	"sync"
	"time"
)

type Builder struct {
	lock    sync.Mutex
	waiters sync.WaitGroup
	running int
	Count   int
}

func (b *Builder) getConcurrency() int {
	settings := node.Self.Settings.(*node.BuilderSettings)
	return settings.Concurrency
}

func (b *Builder) acquire() {
	for {
		b.lock.Lock()
		if b.running < b.getConcurrency() {
			b.running += 1
			node.Self.Stats.(*node.BuilderStats).Active = b.running
			b.lock.Unlock()

			break
		} else {
			b.lock.Unlock()
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func (b *Builder) build(bild *build.Build) {
	b.waiters.Add(1)

	go func() {
		defer func() {
			b.lock.Lock()
			b.running -= 1
			node.Self.Stats.(*node.BuilderStats).Active = b.running
			b.lock.Unlock()
			b.waiters.Done()
		}()

		db := database.GetDatabase()
		defer db.Close()

		err := bild.Build(db)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("builder: Build failed")
			return
		}
	}()
}

func (b *Builder) Start() (err error) {
	db := database.GetDatabase()
	defer db.Close()

	for {
		b.acquire()

		bild, e := build.GetQueued(db)
		if e != nil {
			err = e
			return
		}

		if bild == nil {
			break
		}

		b.build(bild)
	}

	b.waiters.Wait()

	return
}

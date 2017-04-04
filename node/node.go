package node

import (
	"github.com/Sirupsen/logrus"
	"github.com/autoabs/autoabs/database"
	"github.com/autoabs/autoabs/utils"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Node struct {
	Id        string      `bson:"_id" json:"id"`
	Type      string      `bson:"type" json:"type"`
	Timestamp time.Time   `bson:"timestamp" json:"timestamp"`
	Memory    float64     `bson:"memory" json:"memory"`
	Load1     float64     `bson:"load1" json:"load1"`
	Load5     float64     `bson:"load5" json:"load5"`
	Load15    float64     `bson:"load15" json:"load15"`
	Data      interface{} `bson:"-" json:"data"`
}

func (n *Node) LoadData() (err error) {
	if n.Type == "builder" {
		n.Data = &BuilderData{
			Concurrency: 4, // TODO
		}
	}

	return
}

func (n *Node) keepalive() {
	db := database.GetDatabase()
	coll := db.Nodes()

	for {
		n.Timestamp = time.Now()

		mem, err := utils.MemoryUsed()
		if err != nil {
			n.Memory = 0

			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("node: Failed to get memory")
		} else {
			n.Memory = mem
		}

		load, err := utils.LoadAverage()
		if err != nil {
			n.Load1 = 0
			n.Load5 = 0
			n.Load15 = 0

			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("node: Failed to get load")
		} else {
			n.Load1 = load.Load1
			n.Load5 = load.Load5
			n.Load15 = load.Load15
		}

		coll.Upsert(&bson.M{
			"_id": n.Id,
		}, n)

		time.Sleep(10 * time.Second)
	}
}

func (n *Node) Keepalive() (err error) {
	go n.keepalive()

	return
}

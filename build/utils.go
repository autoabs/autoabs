package build

import (
	"github.com/autoabs/autoabs/database"
	"gopkg.in/mgo.v2/bson"
)

func GetQueuedBuilds(db *database.Database) (builds []*Build, err error) {
	builds = []*Build{}
	coll := db.Builds()

	cursor := coll.Find(&bson.M{
		"state": "pending",
	}).Iter()

	bild := &Build{}
	for cursor.Next(bild) {
		builds = append(builds, bild)
		bild = &Build{}
	}

	err = cursor.Close()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
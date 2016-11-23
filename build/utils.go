package build

import (
	"github.com/autoabs/autoabs/database"
	"gopkg.in/mgo.v2/bson"
)

func GetQueued(db *database.Database) (builds []*Build, err error) {
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

func RetryFailed() (err error) {
	db := database.GetDatabase()
	defer db.Close()
	coll := db.Builds()

	cursor := coll.Find(&bson.M{
		"state": "failed",
	}).Iter()

	bild := &Build{}
	for cursor.Next(bild) {
		err = bild.Retry(db)
		if err != nil {
			return
		}
	}

	return
}

func ClearAll() (err error) {
	db := database.GetDatabase()
	defer db.Close()
	coll := db.Builds()

	cursor := coll.Find(&bson.M{}).Iter()

	bild := &Build{}
	for cursor.Next(bild) {
		err = bild.Remove(db)
		if err != nil {
			return
		}
	}

	return
}

func ClearPending() (err error) {
	db := database.GetDatabase()
	defer db.Close()
	coll := db.Builds()

	cursor := coll.Find(&bson.M{
		"state": "pending",
	}).Iter()

	bild := &Build{}
	for cursor.Next(bild) {
		err = bild.Remove(db)
		if err != nil {
			return
		}
	}

	return
}

func ClearFailed() (err error) {
	db := database.GetDatabase()
	defer db.Close()
	coll := db.Builds()

	cursor := coll.Find(&bson.M{
		"state": "failed",
	}).Iter()

	bild := &Build{}
	for cursor.Next(bild) {
		err = bild.Remove(db)
		if err != nil {
			return
		}
	}

	return
}

package build

import (
	"github.com/autoabs/autoabs/database"
	"github.com/autoabs/autoabs/event"
	"github.com/autoabs/autoabs/utils"
	"gopkg.in/mgo.v2/bson"
)

func PublishChange(db *database.Database) (err error) {
	evt := &event.Dispatch{
		Type: "build.change",
	}

	err = event.Publish(db, "dispatch", evt)
	if err != nil {
		return
	}

	return
}

func Get(db *database.Database, buildId bson.ObjectId) (
	bild *Build, err error) {

	coll := db.Builds()

	bild = &Build{}
	err = coll.FindId(buildId).One(bild)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetKey(db *database.Database, name, version, release,
	repo, arch string) (bild *Build, err error) {

	coll := db.Builds()

	bild = &Build{}
	err = coll.Find(&bson.M{
		"name":    name,
		"version": version,
		"release": release,
		"repo":    repo,
		"arch":    arch,
	}).One(bild)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetAll(db *database.Database, index int) (
	builds []*Build, queryIndex, count int, err error) {

	builds = []*Build{}
	coll := db.Builds()

	count, err = coll.Count()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	queryIndex = utils.Min(index, utils.Max(0, count-500))

	cursor := coll.Find(&bson.M{}).Sort(
		"state_rank").Skip(queryIndex).Limit(500).Iter()

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

func GetReady(db *database.Database) (builds []*Build, err error) {
	builds = []*Build{}
	coll := db.Builds()

	cursor := coll.Find(&bson.M{
		"state":    "completed",
		"uploaded": false,
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
		err = bild.Rebuild(db)
		if err != nil {
			return
		}
	}

	err = cursor.Close()
	if err != nil {
		err = database.ParseError(err)
		return
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

	err = cursor.Close()
	if err != nil {
		err = database.ParseError(err)
		return
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

	err = cursor.Close()
	if err != nil {
		err = database.ParseError(err)
		return
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

	err = cursor.Close()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func GetLog(db *database.Database, buildId bson.ObjectId) (
	log []string, err error) {

	coll := db.BuildsLog()
	log = []string{}
	reversed := []string{}

	cursor := coll.Find(&bson.M{
		"b": buildId,
	}).Sort("-t").Limit(1000).Iter()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	entry := &BuildLog{}
	for cursor.Next(entry) {
		reversed = append(reversed, entry.Log)
	}

	for i := len(reversed) - 1; i >= 0; i-- {
		log = append(log, reversed[i])
	}

	err = cursor.Err()
	if err != nil {
		return
	}

	return
}

func ClearLog(db *database.Database, buildId bson.ObjectId) (err error) {
	coll := db.BuildsLog()

	_, err = coll.RemoveAll(&bson.M{
		"b": buildId,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

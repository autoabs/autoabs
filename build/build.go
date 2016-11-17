package build

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Build struct {
	Id          bson.ObjectId   `bson:"_id"`
	Start       time.Time       `bson:"start"`
	Stop        time.Time       `bson:"stop"`
	State       string          `bson:"state"`
	Name        string          `bson:"name"`
	Version     string          `bson:"version"`
	Release     string          `bson:"release"`
	Repo        string          `bson:"core"`
	Arch        string          `bson:"arch"`
	PkgIds      []bson.ObjectId `bson:"pkg_ids"`
	PkgBuildIds []bson.ObjectId `bson:"pkg_build_ids"`
}

package build

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Build struct {
	Id          bson.ObjectId   `bson:"_id"`
	Name        string          `bson:"name"`
	Start       time.Time       `bson:"start,omitempty"`
	Stop        time.Time       `bson:"stop,omitempty"`
	State       string          `bson:"state"`
	Version     string          `bson:"version"`
	Release     string          `bson:"release"`
	Repo        string          `bson:"core"`
	Arch        string          `bson:"arch"`
	PkgIds      []bson.ObjectId `bson:"pkg_ids"`
	PkgBuildIds []bson.ObjectId `bson:"pkg_build_ids"`
}

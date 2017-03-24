package queue

import (
	"github.com/autoabs/autoabs/build"
	"github.com/autoabs/autoabs/database"
	"github.com/autoabs/autoabs/pkg"
	"github.com/autoabs/autoabs/source"
	"github.com/dropbox/godropbox/container/set"
	"gopkg.in/mgo.v2/bson"
)

type Queue struct {
	sources      map[string]*source.Source
	sourcesKeys  set.Set
	packages     map[string]*pkg.Package
	packagesKeys set.Set
	oldPackages  []*pkg.Package
	add          set.Set // *source.Source
	remove       set.Set // *pkg.Package
	fix          set.Set // *pkg.Package
	update       set.Set // *source.Source
}

func (q *Queue) Scan() (err error) {
	q.add = set.NewSet()
	q.remove = set.NewSet()
	q.fix = set.NewSet()
	q.update = set.NewSet()

	q.packages, q.oldPackages, q.packagesKeys, err = pkg.GetAll()
	if err != nil {
		return
	}

	q.sources, q.sourcesKeys, err = source.GetAll()
	if err != nil {
		return
	}

	remPackagesKeys := q.packagesKeys.Copy()
	remPackagesKeys.Subtract(q.sourcesKeys)
	for key := range remPackagesKeys.Iter() {
		q.remove.Add(q.packages[key.(string)])
	}

	for _, pk := range q.oldPackages {
		q.remove.Add(pk)
		q.fix.Add(q.packages[pk.Key()])
	}

	addPackagesKeys := q.sourcesKeys.Copy()
	addPackagesKeys.Subtract(q.packagesKeys)
	for key := range addPackagesKeys.Iter() {
		q.add.Add(q.sources[key.(string)])
	}

	for key, pk := range q.packages {
		src, ok := q.sources[key]
		if !ok {
			continue
		}

		if src.Version != pk.Version || src.Release != pk.Release {
			q.update.Add(src)
		}
	}

	return
}

func (q *Queue) Sync() (err error) {
	db := database.GetDatabase()
	defer db.Close()

	err = q.Scan()
	if err != nil {
		return
	}

	for pkInf := range q.fix.Iter() {
		pk := pkInf.(*pkg.Package)

		err = pk.Fix()
		if err != nil {
			return
		}
	}

	for pkInf := range q.remove.Iter() {
		pk := pkInf.(*pkg.Package)
		pk.Remove()
	}

	changed := false

	for srcInf := range q.add.Iter() {
		src := srcInf.(*source.Source)

		queued, e := src.Queue(db, true)
		if e != nil {
			err = e
			return
		}

		if queued {
			changed = true
		}
	}

	for srcInf := range q.update.Iter() {
		src := srcInf.(*source.Source)

		queued, e := src.Queue(db, false)
		if e != nil {
			err = e
			return
		}

		if queued {
			changed = true
		}
	}

	if q.remove.Len() != 0 || len(q.oldPackages) != 0 ||
		q.fix.Len() != 0 || changed {

		build.PublishChange(db)
	}

	return
}

func (q *Queue) SyncState() (err error) {
	stateId := bson.NewObjectId()
	db := database.GetDatabase()
	defer db.Close()

	err = q.Scan()
	if err != nil {
		return
	}

	for _, pk := range q.packages {
		err = pk.SyncState(db, stateId)
		if err != nil {
			return
		}
	}

	coll := db.Builds()

	_, err = coll.UpdateAll(&bson.M{
		"repo_state": &bson.M{
			"$ne": stateId,
		},
	}, &bson.M{
		"$set": &bson.M{
			"repo_state": nil,
		},
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func (q *Queue) Build() (err error) {
	db := database.GetDatabase()
	defer db.Close()

	builds, err := build.GetQueued(db)
	if err != nil {
		return
	}

	for _, bild := range builds {
		err = bild.Build(db)
		if err != nil {
			return
		}
	}

	return
}

func (q *Queue) Upload() (err error) {
	db := database.GetDatabase()
	defer db.Close()

	err = q.Scan()
	if err != nil {
		return
	}

	builds, err := build.GetReady(db)
	if err != nil {
		return
	}

	for _, bild := range builds {
		pk, ok := q.packages[bild.Key()]
		if ok {
			if !utils.VersionNewer(bild.FullVersion(), pk.FullVersion()) {
				continue
			}
		}

		err = bild.Upload(db, false)
		if err != nil {
			return
		}
	}

	return
}

func (q *Queue) Clean() (err error) {
	return
}

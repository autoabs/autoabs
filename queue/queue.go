package queue

import (
	"github.com/autoabs/autoabs/build"
	"github.com/autoabs/autoabs/database"
	"github.com/autoabs/autoabs/pkg"
	"github.com/dropbox/godropbox/container/set"
)

type Queue struct {
	curPackages     map[string]*pkg.Package
	curPackagesKeys set.Set
	newPackages     map[string]*pkg.Package
	newPackagesKeys set.Set
	addPackages     []*pkg.Package
	remPackages     []*pkg.Package
	updatePackages  []*pkg.Package
	buildPackages   []*pkg.Package
}

func (q *Queue) Scan() (err error) {
	q.curPackages = map[string]*pkg.Package{}
	q.curPackagesKeys = set.NewSet()
	q.newPackages = map[string]*pkg.Package{}
	q.newPackagesKeys = set.NewSet()
	q.addPackages = []*pkg.Package{}
	q.remPackages = []*pkg.Package{}
	q.updatePackages = []*pkg.Package{}

	curPkgs, err := getCurPackages()
	if err != nil {
		return
	}

	for _, pk := range curPkgs {
		key := pk.Key()
		q.curPackages[key] = pk
		q.curPackagesKeys.Add(key)
	}

	newPkgs, err := getNewPackages()
	if err != nil {
		return
	}

	for _, pk := range newPkgs {
		key := pk.Key()
		q.newPackages[key] = pk
		q.newPackagesKeys.Add(key)
	}

	remPackagesKeys := q.curPackagesKeys.Copy()
	remPackagesKeys.Subtract(q.newPackagesKeys)
	for key := range remPackagesKeys.Iter() {
		q.remPackages = append(q.remPackages, q.curPackages[key.(string)])
	}

	addPackagesKeys := q.newPackagesKeys.Copy()
	addPackagesKeys.Subtract(q.curPackagesKeys)
	for key := range addPackagesKeys.Iter() {
		q.addPackages = append(q.addPackages, q.newPackages[key.(string)])
	}

	for key, pk := range q.curPackages {
		newPkg, ok := q.newPackages[key]
		if !ok {
			continue
		}

		if newPkg.Version != pk.Version || newPkg.Release != pk.Release {
			q.updatePackages = append(q.updatePackages, pk)
		}
	}

	return
}

func (q *Queue) Queue() (err error) {
	err = q.Scan()
	if err != nil {
		return
	}

	for _, pk := range q.remPackages {
		pk.Remove()
	}

	queued := set.NewSet()

	for _, pk := range q.addPackages {
		key := pk.IdKey()
		if queued.Contains(key) {
			continue
		}
		queued.Add(key)

		err = pk.QueueBuild(false)
		if err != nil {
			return
		}
	}

	for _, pk := range q.updatePackages {
		key := pk.IdKey()
		if queued.Contains(key) {
			continue
		}
		queued.Add(key)

		err = pk.QueueBuild(false)
		if err != nil {
			return
		}
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

func (q *Queue) Clean() (err error) {
	err = q.Scan()
	if err != nil {
		return
	}

	curPkgs, err := getCurPackages()
	if err != nil {
		return
	}

	for _, pk := range curPkgs {
		newPkg, ok := q.newPackages[pk.Key()]
		if !ok {
			continue
		}

		if newPkg.Version != pk.Version || newPkg.Release != pk.Release {
			pk.Remove()
		}
	}

	return
}

package queue

import (
	"github.com/autoabs/autoabs/config"
	"github.com/autoabs/autoabs/pkg"
	"github.com/autoabs/autoabs/utils"
	"github.com/dropbox/godropbox/container/set"
	"path"
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
	hasDir, err := utils.ContainsDir(
		path.Join(config.Config.RootPath, "builds"))
	if err != nil {
		return
	}

	if hasDir {
		return
	}

	err = q.Scan()
	if err != nil {
		return
	}

	for _, pk := range q.remPackages {
		pk.Print()
		pk.Remove()
	}

	queued := set.NewSet()

	for _, pk := range q.addPackages {
		key := pk.IdKey()
		if queued.Contains(key) {
			continue
		}
		queued.Add(key)

		pk.Print()

		err = pk.QueueBuild()
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

		pk.Print()

		err = pk.QueueBuild()
		if err != nil {
			return
		}
	}

	return
}

func (q *Queue) Build() (err error) {
	q.buildPackages, err = getBuildPackages()
	if err != nil {
		return
	}

	for _, pk := range q.buildPackages {
		pk.Print()

		err = pk.Build()
		if err != nil {
			return
		}
	}

	return
}

func (q *Queue) Clean() (err error) {
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

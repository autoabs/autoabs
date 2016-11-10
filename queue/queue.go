package queue

import (
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
	dupPackages     []*pkg.Package
	updatePackages  []*pkg.Package
}

func (q *Queue) Queue() (err error) {
	q.curPackages = map[string]*pkg.Package{}
	q.curPackagesKeys = set.NewSet()
	q.newPackages = map[string]*pkg.Package{}
	q.newPackagesKeys = set.NewSet()
	q.addPackages = []*pkg.Package{}
	q.remPackages = []*pkg.Package{}
	q.dupPackages = []*pkg.Package{}
	q.updatePackages = []*pkg.Package{}

	curPkgs, err := getCurPackages()
	if err != nil {
		return
	}

	for _, pk := range curPkgs {
		key := pk.Key()
		q.curPackages[key] = pk

		if q.curPackagesKeys.Contains(key) {
			q.dupPackages = append(q.dupPackages, pk)
		} else {
			q.curPackagesKeys.Add(key)
		}
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

	for key, pkg := range q.curPackages {
		prevPkg, ok := q.newPackages[key]
		if !ok {
			continue
		}

		if prevPkg.Version != pkg.Version || prevPkg.Release != pkg.Release {
			pkg.Previous = prevPkg
			q.updatePackages = append(q.updatePackages, pkg)
		}
	}

	return
}

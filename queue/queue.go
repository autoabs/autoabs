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
}

func (q *Queue) Build() (err error) {
	q.curPackages = map[string]*pkg.Package{}
	q.curPackagesKeys = set.NewSet()
	q.newPackages = map[string]*pkg.Package{}
	q.newPackagesKeys = set.NewSet()
	q.addPackages = []*pkg.Package{}
	q.remPackages = []*pkg.Package{}

	curPkgs, err := getCurPackages()
	if err != nil {
		return
	}

	for _, pk := range curPkgs {
		key := pk.Key()
		q.curPackages[key] = pk
		q.curPackagesKeys.Add(key)
	}

	return
}

package queue

import (
	"github.com/autoabs/autoabs/pkg"
)

type Queue struct {
	curPackages map[string]*pkg.Package
	newPackages map[string]*pkg.Package
}

func Build() (err error) {
	curPkgs, err := getCurPackages()
	if err != nil {
		return
	}

	_ = curPkgs

	return
}

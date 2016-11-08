package queue

import (
	"fmt"
)

type Package struct {
	Name    string
	Version string
	Release string
	Repo    string
	Arch    string
}

func (p *Package) Key() string {
	return p.Name + "-" + p.Repo + "-" + p.Arch + "-" + p.Version
}

func (p *Package) Print() {
	fmt.Printf("%s %s-%s: %s:%s\n",
		p.Name, p.Version, p.Release, p.Repo, p.Arch)
}

type Queue struct {
	curPackages map[string]*Package
	newPackages map[string]*Package
}

func Build() (err error) {
	curPkgs, err := getCurPackages()
	if err != nil {
		return
	}

	_ = curPkgs

	return
}

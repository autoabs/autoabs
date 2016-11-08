package queue

import (
	"fmt"
	"github.com/autoabs/autoabs/config"
	"github.com/autoabs/autoabs/errortypes"
	"github.com/autoabs/autoabs/utils"
	"github.com/dropbox/godropbox/errors"
	"io/ioutil"
	"os/exec"
	"path"
	"strings"
)

type Package struct {
	Name    string
	Version string
	Release string
	Repo    string
	Arch    string
}

type Queue struct {
	curPackages map[string]*Package
	newPackages map[string]*Package
}

func scanRepos(pkgName, pth string) (pkgs []*Package, err error) {
	pkgs = []*Package{}

	repos, err := ioutil.ReadDir(pth)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrapf(err, "queue: Failed to read dir %s", pth),
		}
		return
	}

	for _, entry := range repos {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()

		split := strings.Split(name, "-")
		if len(split) != 2 {
			continue
		}

		repo := split[0]
		arch := split[1]

		if !config.Config.TargetRepos.Contains(repo) ||
			!config.Config.TargetArchs.Contains(arch) {

			continue
		}

		pkgBuildPath := path.Join(pth, name, "PKGBUILD")

		exists, e := utils.ExistsFile(pkgBuildPath)
		if e != nil {
			err = e
			return
		}

		if !exists {
			continue
		}

		cmd := exec.Command(
			"/bin/sh",
			"-c",
			fmt.Sprintf(
				`sh -c 'source "%s"; echo $pkgver-$pkgrel'`,
				pkgBuildPath,
			),
		)

		output, e := cmd.Output()
		if e != nil {
			err = errortypes.ExecError{
				errors.Wrapf(err, "queue: Failed to get package version"),
			}
			return
		}

		pkgInfo := strings.Split(strings.TrimSpace(string(output)), "-")

		pkg := &Package{
			Name:    pkgName,
			Version: pkgInfo[0],
			Release: pkgInfo[1],
			Repo:    repo,
			Arch:    arch,
		}
		pkgs = append(pkgs, pkg)
	}

	return
}

func scanSource(pth string) (pkgs []*Package, err error) {
	pkgs = []*Package{}

	packages, err := ioutil.ReadDir(pth)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrapf(err, "queue: Failed to read dir %s", pth),
		}
		return
	}

	for _, entry := range packages {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()

		reposPath := path.Join(pth, name, "repos")

		exists, e := utils.ExistsDir(reposPath)
		if e != nil {
			err = e
			return
		}

		if !exists {
			continue
		}

		reposPkgs, e := scanRepos(name, reposPath)
		if e != nil {
			err = e
			return
		}

		pkgs = append(pkgs, reposPkgs...)
	}

	return
}

func getCurPackages() (pkgs []*Package, err error) {
	pth := path.Join(config.Config.RootPath, "sources")
	pkgs = []*Package{}

	exists, err := utils.ExistsDir(pth)
	if err != nil {
		return
	}

	if !exists {
		return
	}

	sources, err := ioutil.ReadDir(pth)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrapf(err, "queue: Failed to read dir %s", pth),
		}
		return
	}

	for _, entry := range sources {
		if !entry.IsDir() {
			continue
		}

		sourcePath := path.Join(pth, entry.Name())

		sourcePkgs, e := scanSource(sourcePath)
		if e != nil {
			err = e
			return
		}

		pkgs = append(pkgs, sourcePkgs...)
	}

	return
}

func Build() (err error) {
	curPkgs, err := getCurPackages()
	if err != nil {
		return
	}

	_ = curPkgs

	return
}

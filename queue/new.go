package queue

import (
	"fmt"
	"github.com/autoabs/autoabs/config"
	"github.com/autoabs/autoabs/errortypes"
	"github.com/autoabs/autoabs/pkg"
	"github.com/autoabs/autoabs/utils"
	"github.com/dropbox/godropbox/errors"
	"io/ioutil"
	"os/exec"
	"path"
	"strings"
)

func scanNewRepos(pkgName, pth string) (pkgs []*pkg.Package, err error) {
	pkgs = []*pkg.Package{}

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

		pk := &pkg.Package{
			Name:    pkgName,
			Version: pkgInfo[0],
			Release: pkgInfo[1],
			Repo:    repo,
			Arch:    arch,
		}
		pkgs = append(pkgs, pk)

		pk.Print()
	}

	return
}

func scanNewSource(pth string) (pkgs []*pkg.Package, err error) {
	pkgs = []*pkg.Package{}

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

		reposPkgs, e := scanNewRepos(name, reposPath)
		if e != nil {
			err = e
			return
		}

		pkgs = append(pkgs, reposPkgs...)
	}

	return
}

func getNewPackages() (pkgs []*pkg.Package, err error) {
	pth := path.Join(config.Config.RootPath, "sources")
	pkgs = []*pkg.Package{}

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

		sourcePkgs, e := scanNewSource(sourcePath)
		if e != nil {
			err = e
			return
		}

		pkgs = append(pkgs, sourcePkgs...)
	}

	return
}
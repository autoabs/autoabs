package queue

import (
	"fmt"
	"github.com/autoabs/autoabs/config"
	"github.com/autoabs/autoabs/errortypes"
	"github.com/autoabs/autoabs/pkg"
	"github.com/autoabs/autoabs/settings"
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

		if !settings.System.HasRepo(repo) || !settings.System.HasArch(arch) {
			continue
		}

		sourcePath := path.Join(pth, name)
		pkgBuildPath := path.Join(sourcePath, "PKGBUILD")

		exists, e := utils.ExistsFile(pkgBuildPath)
		if e != nil {
			err = e
			return
		}

		if !exists {
			continue
		}

		cmd := exec.Command(
			"/bin/bash",
			"-c",
			fmt.Sprintf(
				`source "%s"; echo ${pkgname[*]}:$pkgver:$pkgrel`,
				pkgBuildPath,
			),
		)

		output, e := cmd.Output()
		if e != nil {
			err = errortypes.ExecError{
				errors.Wrapf(e, "queue: Failed to get package version"),
			}
			return
		}

		pkgInfo := strings.Split(strings.TrimSpace(string(output)), ":")

		subNames := strings.Split(pkgInfo[0], " ")

		for _, subName := range subNames {
			pk := &pkg.Package{
				Name:       pkgName,
				SubName:    subName,
				SubNames:   subNames,
				Version:    pkgInfo[1],
				Release:    pkgInfo[2],
				Repo:       repo,
				Arch:       arch,
				SourcePath: sourcePath,
			}
			pkgs = append(pkgs, pk)
		}
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

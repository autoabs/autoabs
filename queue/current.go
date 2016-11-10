package queue

import (
	"github.com/autoabs/autoabs/config"
	"github.com/autoabs/autoabs/constants"
	"github.com/autoabs/autoabs/errortypes"
	"github.com/autoabs/autoabs/pkg"
	"github.com/autoabs/autoabs/utils"
	"github.com/dropbox/godropbox/errors"
	"io/ioutil"
	"path"
	"strings"
)

func scanCurPackages(repo, arch, pth string) (pkgs []*pkg.Package, err error) {
	pkgs = []*pkg.Package{}

	packages, err := ioutil.ReadDir(pth)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrapf(err, "queue: Failed to read dir %s", pth),
		}
		return
	}

	for _, entry := range packages {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		if !strings.HasSuffix(name, constants.PackageExt) {
			continue
		}

		pkgPth := path.Join(pth, name)

		nameSpl := strings.Split(name, "-")
		len := len(nameSpl)

		if len < 4 {
			continue
		}

		if arch != strings.Split(nameSpl[len-1], ".")[0] {
			continue
		}

		pk := &pkg.Package{
			Name:    strings.Join(nameSpl[:len-3], "-"),
			Version: nameSpl[len-3],
			Release: nameSpl[len-2],
			Repo:    repo,
			Arch:    arch,
			Path:    pkgPth,
		}
		pkgs = append(pkgs, pk)

		pk.Print()
	}

	return
}

func scanCurArchs(repo, pth string) (pkgs []*pkg.Package, err error) {
	pkgs = []*pkg.Package{}

	archs, err := ioutil.ReadDir(pth)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrapf(err, "queue: Failed to read dir %s", pth),
		}
		return
	}

	for _, entry := range archs {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()

		packagesPath := path.Join(pth, name)

		pks, e := scanCurPackages(repo, name, packagesPath)
		if e != nil {
			err = e
			return
		}

		pkgs = append(pkgs, pks...)
	}

	return
}

func getCurPackages() (pkgs []*pkg.Package, err error) {
	pth := path.Join(config.Config.RootPath, "repo")
	pkgs = []*pkg.Package{}

	exists, err := utils.ExistsDir(pth)
	if err != nil {
		return
	}

	if !exists {
		return
	}

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

		archsPath := path.Join(pth, name, "os")

		exists, e := utils.ExistsDir(archsPath)
		if e != nil {
			err = e
			return
		}

		if !exists {
			continue
		}

		archsPkgs, e := scanCurArchs(name, archsPath)
		if e != nil {
			err = e
			return
		}

		pkgs = append(pkgs, archsPkgs...)
	}

	return
}

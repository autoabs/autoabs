package pkg

import (
	"github.com/autoabs/autoabs/config"
	"github.com/autoabs/autoabs/constants"
	"github.com/autoabs/autoabs/errortypes"
	"github.com/autoabs/autoabs/utils"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"io/ioutil"
	"path"
	"strings"
)

type scanner struct {
	Packages    map[string]*Package
	OldPackages []*Package
	Keys        set.Set
}

func (s *scanner) scanPackages(repo, arch, pth string) (err error) {
	packages, err := ioutil.ReadDir(pth)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrapf(err, "pkg: Failed to read dir %s", pth),
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
		ln := len(nameSpl)

		if ln < 4 {
			continue
		}

		if arch != strings.Split(nameSpl[ln-1], ".")[0] {
			continue
		}

		pk := &Package{
			Name:    strings.Join(nameSpl[:ln-3], "-"),
			Version: nameSpl[ln-3],
			Release: nameSpl[ln-2],
			Repo:    repo,
			Arch:    arch,
			Path:    pkgPth,
		}
		key := pk.Key()

		pk2 := s.Packages[key]
		if pk2 != nil {
			if utils.VersionNewer(pk.FullVersion(), pk2.FullVersion()) {
				s.OldPackages = append(s.OldPackages, pk2)
				s.Packages[key] = pk
			} else {
				s.OldPackages = append(s.OldPackages, pk)
			}
		} else {
			s.Packages[key] = pk
		}

		s.Keys.Add(key)
	}

	return
}

func (s *scanner) scanArchs(repo, pth string) (err error) {
	archs, err := ioutil.ReadDir(pth)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrapf(err, "pkg: Failed to read dir %s", pth),
		}
		return
	}

	for _, entry := range archs {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()

		packagesPath := path.Join(pth, name)

		err = s.scanPackages(repo, name, packagesPath)
		if err != nil {
			return
		}
	}

	return
}

func (s *scanner) Scan() (err error) {
	pth := path.Join(config.Config.RootPath, "repo")

	s.Packages = map[string]*Package{}
	s.OldPackages = []*Package{}
	s.Keys = set.NewSet()

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
			errors.Wrapf(err, "pkg: Failed to read dir %s", pth),
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

		err = s.scanArchs(name, archsPath)
		if err != nil {
			return
		}
	}

	return
}

func GetAll() (packages map[string]*Package, oldPackages []*Package,
	keys set.Set, err error) {

	scnr := &scanner{}

	err = scnr.Scan()
	if err != nil {
		return
	}
	packages = scnr.Packages
	oldPackages = scnr.OldPackages
	keys = scnr.Keys

	return
}

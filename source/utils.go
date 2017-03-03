package source

import (
	"fmt"
	"github.com/autoabs/autoabs/config"
	"github.com/autoabs/autoabs/errortypes"
	"github.com/autoabs/autoabs/settings"
	"github.com/autoabs/autoabs/utils"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"io/ioutil"
	"os/exec"
	"path"
	"strings"
)

type scanner struct {
	Sources map[string]*Source
	Keys    set.Set
}

func (s *scanner) scanRepos(pkgName, pth string) (err error) {
	repos, err := ioutil.ReadDir(pth)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrapf(err, "source: Failed to read dir %s", pth),
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
				errors.Wrapf(e, "source: Failed to get package version"),
			}
			return
		}

		pkgInfo := strings.Split(strings.TrimSpace(string(output)), ":")

		subNames := strings.Split(pkgInfo[0], " ")

		source := &Source{
			Name:     pkgName,
			SubNames: subNames,
			Version:  pkgInfo[1],
			Release:  pkgInfo[2],
			Repo:     repo,
			Arch:     arch,
			Path:     sourcePath,
		}

		for _, key := range source.Keys() {
			s.Sources[key] = source
			s.Keys.Add(key)
		}
	}

	return
}

func (s *scanner) scanSource(pth string) (err error) {
	targets, err := ioutil.ReadDir(pth)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrapf(err, "source: Failed to read dir %s", pth),
		}
		return
	}

	for _, entry := range targets {
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

		err = s.scanRepos(name, reposPath)
		if err != nil {
			return
		}
	}

	return
}

func (s *scanner) Scan() (err error) {
	pth := path.Join(config.Config.RootPath, "sources")

	s.Sources = map[string]*Source{}
	s.Keys = set.NewSet()

	exists, err := utils.ExistsDir(pth)
	if err != nil {
		return
	}

	if !exists {
		return
	}

	paths, err := ioutil.ReadDir(pth)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrapf(err, "source: Failed to read dir %s", pth),
		}
		return
	}

	for _, entry := range paths {
		if !entry.IsDir() {
			continue
		}

		sourcePath := path.Join(pth, entry.Name())

		err = s.scanSource(sourcePath)
		if err != nil {
			return
		}
	}

	return
}

func GetAll() (sources map[string]*Source, keys set.Set, err error) {
	scnr := &scanner{}

	err = scnr.Scan()
	if err != nil {
		return
	}
	sources = scnr.Sources
	keys = scnr.Keys

	return
}

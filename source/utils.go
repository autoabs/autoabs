package source

import (
	"fmt"
	"github.com/autoabs/autoabs/config"
	"github.com/autoabs/autoabs/constants"
	"github.com/autoabs/autoabs/database"
	"github.com/autoabs/autoabs/errortypes"
	"github.com/autoabs/autoabs/settings"
	"github.com/autoabs/autoabs/utils"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"path"
	"strings"
)

type scanner struct {
	Sources   map[string]*Source
	Keys      set.Set
	Container string
}

func (s *scanner) Start() (err error) {
	output, err := utils.ExecOutput("",
		"docker",
		"run",
		"-it",
		"-d",
		"-v", fmt.Sprintf("%s:%s", config.Config.RootPath,
			config.Config.RootPath),
		constants.BuildImage,
		"/bin/bash",
	)

	s.Container = strings.TrimSpace(output)

	return
}

func (s *scanner) Stop() {
	if s.Container == "" {
		return
	}

	utils.ExecOutput("",
		"docker",
		"rm",
		"-f",
		s.Container,
	)
}

func (s *scanner) run(cmd string, arg ...string) (output string, err error) {
	args := []string{
		"exec",
		s.Container,
		cmd,
	}
	args = append(args, arg...)

	output, err = utils.ExecOutput("",
		"docker",
		args...,
	)

	return
}

func (s *scanner) scanRepos(db *database.Database,
	pkgName, pth string) (err error) {

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

		output, e := s.run(
			"/bin/bash",
			"-c",
			fmt.Sprintf(
				`source "%s"; echo ${pkgname[*]}:$pkgver:$pkgrel`,
				pkgBuildPath,
			),
		)
		if e != nil {
			err = e
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

		err = source.Upsert(db)
		if err != nil {
			return
		}

		for _, key := range source.Keys() {
			s.Sources[key] = source
			s.Keys.Add(key)
		}
	}

	return
}

func (s *scanner) scanSource(db *database.Database, pth string) (err error) {
	curCommit, err := getCommit(db, pth)
	if err != nil {
		return
	}

	newCommit, err := utils.GitCommit(pth)
	if err != nil {
		return
	}

	if curCommit == "" {
		targets, e := ioutil.ReadDir(pth)
		if e != nil {
			err = &errortypes.ReadError{
				errors.Wrapf(e, "source: Failed to read dir %s", pth),
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

			err = s.scanRepos(db, name, reposPath)
			if err != nil {
				return
			}
		}
	} else {
	}

	err = setCommit(db, pth, newCommit)
	if err != nil {
		return
	}

	return
}

func (s *scanner) preload(db *database.Database) (err error) {
	coll := db.Sources()

	cursor := coll.Find(&bson.M{}).Iter()
	if err != nil {
		return
	}

	source := &Source{}
	for cursor.Next(source) {
		for _, key := range source.Keys() {
			s.Sources[key] = source
			s.Keys.Add(key)
		}
		source = &Source{}
	}

	return
}

func (s *scanner) Scan() (err error) {
	pth := path.Join(config.Config.RootPath, "sources")

	db := database.GetDatabase()
	defer db.Close()

	s.Sources = map[string]*Source{}
	s.Keys = set.NewSet()

	exists, err := utils.ExistsDir(pth)
	if err != nil {
		return
	}

	if !exists {
		return
	}

	err = s.preload(db)
	if err != nil {
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

		err = s.scanSource(db, sourcePath)
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

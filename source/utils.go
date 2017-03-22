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
	"os"
	"path"
	"path/filepath"
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

func (s *scanner) scanPkgbuild(db *database.Database, repo, pth string) (
	err error) {

	pthSpl := strings.Split(pth, "/")

	sourcePath, _ := filepath.Split(pth)
	sourcePath = strings.TrimRight(sourcePath, "/")

	if len(pthSpl) > 3 && pthSpl[len(pthSpl)-3] == "repos" {
		split := strings.Split(pthSpl[len(pthSpl)-2], "-")
		if len(split) == 2 {
			repo = split[0]

			if !settings.System.HasArch(split[1]) {
				return
			}
		}
	}

	if !settings.System.HasRepo(repo) {
		return
	}

	output, e := s.run(
		"/bin/bash",
		"-c",
		fmt.Sprintf(
			`source "%s"; echo %%%%-AUTOABS-%%%%${pkgname[*]}:`+
				`${arch[*]}:$pkgver:$pkgrel%%%%-AUTOABS-%%%%`,
			pth,
		),
	)
	if e != nil {
		err = e
		return
	}

	output = strings.TrimSpace(string(output))
	outputSpl := strings.Split(output, "%%-AUTOABS-%%")
	output = outputSpl[1]

	pkgInfo := strings.Split(output, ":")

	subNames := strings.Split(pkgInfo[0], " ")

	for _, arch := range strings.Split(pkgInfo[1], " ") {
		if !settings.System.HasArch(arch) {
			continue
		}

		source := &Source{
			Name:     subNames[0],
			SubNames: subNames,
			Version:  pkgInfo[2],
			Release:  pkgInfo[3],
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

func (s *scanner) remPkgbuild(db *database.Database, repo, pth string) (
	err error) {

	coll := db.Sources()

	pth = strings.TrimRight(pth, "/PKGBUILD")

	source := &Source{}
	err = coll.FindOne(&bson.M{
		"path": pth,
	}, source)
	if err != nil {
		switch err.(type) {
		case *database.NotFoundError:
			err = nil
		default:
			return
		}
	}

	for _, key := range source.Keys() {
		_, ok := s.Sources[key]
		if ok {
			delete(s.Sources, key)
		}
		s.Keys.Remove(key)
	}

	err = source.Remove(db)
	if err != nil {
		return
	}

	return
}

func (s *scanner) scanSource(db *database.Database, repo, pth string) (
	err error) {

	curCommit, err := getCommit(db, pth)
	if err != nil {
		return
	}

	newCommit, err := utils.GitCommit(pth)
	if err != nil {
		return
	}

	if curCommit == "" {
		err = filepath.Walk(pth, func(pth string,
			info os.FileInfo, err error) (e error) {

			if err != nil {
				e = &errortypes.ReadError{
					errors.Wrap(err,
						"source: Failed to read source directory"),
				}
				return
			}

			if info.IsDir() {
				return
			}

			if !strings.HasSuffix(pth, "PKGBUILD") ||
				strings.HasSuffix(pth, "/trunk/PKGBUILD") {

				return
			}

			err = s.scanPkgbuild(db, repo, pth)
			if err != nil {
				return
			}

			return
		})
	} else if curCommit != newCommit {
		changed, e := utils.GitChanged(pth, curCommit, newCommit)
		if e != nil {
			err = e
			return
		}

		for keyInf := range changed.Iter() {
			pkgPth := filepath.Join(pth, keyInf.(string))

			if !strings.HasSuffix(pkgPth, "PKGBUILD") ||
				strings.HasSuffix(pkgPth, "/trunk/PKGBUILD") {

				continue
			}

			exists, e := utils.ExistsFile(pkgPth)
			if e != nil {
				err = e
				return
			}

			if !exists {
				err = s.remPkgbuild(db, repo, pkgPth)
			} else {
				err = s.scanPkgbuild(db, repo, pkgPth)
				if err != nil {
					return
				}
			}
		}
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

		name := entry.Name()

		sourcePath := path.Join(pth, name)

		err = s.scanSource(db, name, sourcePath)
		if err != nil {
			return
		}
	}

	return
}

func GetAll() (sources map[string]*Source, keys set.Set, err error) {
	scnr := &scanner{}

	err = scnr.Start()
	if err != nil {
		return
	}
	defer scnr.Stop()

	err = scnr.Scan()
	if err != nil {
		return
	}
	sources = scnr.Sources
	keys = scnr.Keys

	return
}

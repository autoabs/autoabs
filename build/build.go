package build

import (
	"archive/tar"
	"bufio"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/autoabs/autoabs/config"
	"github.com/autoabs/autoabs/constants"
	"github.com/autoabs/autoabs/database"
	"github.com/autoabs/autoabs/errortypes"
	"github.com/autoabs/autoabs/signing"
	"github.com/autoabs/autoabs/utils"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"gopkg.in/mgo.v2/bson"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"
)

type Build struct {
	Id         bson.ObjectId   `bson:"_id" json:"id"`
	Name       string          `bson:"name" json:"name"`
	SubNames   []string        `bson:"sub_names" json:"sub_names"`
	Builder    string          `bson:"builder" json:"builder"`
	Start      time.Time       `bson:"start,omitempty" json:"start,omitempty"`
	Stop       time.Time       `bson:"stop,omitempty" json:"stop,omitempty"`
	State      string          `bson:"state" json:"state"`
	StateRank  int             `bson:"state_rank" json:"state_rank"`
	Version    string          `bson:"version" json:"version"`
	Release    string          `bson:"release" json:"release"`
	Repo       string          `bson:"repo" json:"repo"`
	RepoState  bson.ObjectId   `bson:"repo_state,omitempty" json:"repo_state"`
	Uploaded   bool            `bson:"uploaded" json:"uploaded"`
	Arch       string          `bson:"arch" json:"arch"`
	Log        []string        `bson:"log,omitempty" json:"log"`
	PkgIds     []bson.ObjectId `bson:"pkg_ids" json:"pkg_ids"`
	PkgBuildId bson.ObjectId   `bson:"pkg_build_id" json:"pkg_build_id"`
}

type BuildLog struct {
	Build     bson.ObjectId `bson:"b"`
	Timestamp time.Time     `bson:"t"`
	Log       string        `bson:"l"`
}

func (b *Build) tmpPath() string {
	return path.Join(config.Config.RootPath, "tmp",
		b.Name+"-"+b.Repo+"-"+b.Arch+"-"+b.Version+"-"+b.Release)
}

func (b *Build) repoPath() string {
	return path.Join(config.Config.RootPath, "repo", b.Repo, "os", b.Arch)
}

func (b *Build) databasePath() string {
	return path.Join(config.Config.RootPath, "repo", b.Repo, "os", b.Arch,
		fmt.Sprintf("%s.db.tar.gz", b.Repo))
}

func (b *Build) extract(db *database.Database) (err error) {
	tmpPath := b.tmpPath()

	gfs := db.PkgBuildGrid()

	gf, err := gfs.OpenId(b.PkgBuildId)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	arc := tar.NewReader(gf)

	for {
		hdr, e := arc.Next()
		if e == io.EOF {
			break
		}
		if e != nil {
			err = &errortypes.ReadError{
				errors.Wrap(e, "build: Failed to read tar header"),
			}
			return
		}

		pth := path.Join(tmpPath, hdr.Name)
		dirPth := path.Dir(pth)

		err = utils.ExistsMkdir(dirPth, 0700)
		if err != nil {
			return
		}

		file, e := os.OpenFile(
			pth,
			os.O_RDWR|os.O_CREATE|os.O_TRUNC,
			os.FileMode(hdr.Mode),
		)
		if e != nil {
			err = &errortypes.WriteError{
				errors.Wrap(e, "build: Failed to open tar file"),
			}
			return
		}

		_, err = io.Copy(file, arc)
		if err != nil {
			return
		}

		err = file.Close()
		if err != nil {
			err = &errortypes.WriteError{
				errors.Wrap(err, "build: Failed to write tar file"),
			}
			return
		}
	}

	return
}

func (b *Build) extractTest(db *database.Database) (err error) {
	tmpPath := b.tmpPath()

	err = utils.ExistsMkdir(tmpPath, 0700)
	if err != nil {
		return
	}

	pkgNames := []string{}
	pkgPackages := []string{}

	for _, name := range b.SubNames {
		pkgNames = append(pkgNames, fmt.Sprintf("'%s'", name))
		pkgPackages = append(pkgPackages,
			fmt.Sprintf(testPkgBuildPackage, name, name))
	}

	pkgNamesStr := strings.Join(pkgNames, " ")
	pkgPackagesStr := strings.Join(pkgPackages, "\n")

	data := fmt.Sprintf(testPkgBuild, pkgNamesStr, b.Version,
		b.Release, b.Arch, pkgPackagesStr)

	err = ioutil.WriteFile(path.Join(tmpPath, "PKGBUILD"), []byte(data), 0644)
	if err != nil {
		return
	}

	return
}

func (b *Build) addFile(db *database.Database, pkgPath string) (err error) {
	gfs := db.PkgGrid()
	coll := db.Builds()

	_, name := path.Split(pkgPath)

	gf, err := gfs.Create(name)
	if err != nil {
		err = database.ParseError(err)
		return
	}
	gfId := gf.Id()

	file, e := os.Open(pkgPath)
	if e != nil {
		e = &errortypes.ReadError{
			errors.Wrap(e, "build: Failed to open build file"),
		}
		return
	}
	defer file.Close()

	_, e = io.Copy(gf, file)
	if e != nil {
		e = &errortypes.WriteError{
			errors.Wrap(e, "build: Failed to read source file"),
		}
		return
	}

	err = gf.Close()
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrap(err, "build: Failed to close grid file"),
		}
		return
	}

	err = coll.UpdateId(b.Id, &bson.M{
		"$push": &bson.M{
			"pkg_ids": gfId,
		},
	})
	if err != nil {
		err = database.ParseError(err)
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("build: Failed to add pkg id")
		return
	}

	return
}

func (b *Build) build(db *database.Database) (err error) {
	coll := db.BuildsLog()
	tmpPath := b.tmpPath()

	err = ClearLog(db, b.Id)
	if err != nil {
		return
	}

	err = utils.ExistsRemove(tmpPath)
	if err != nil {
		return
	}

	defer utils.ExistsRemove(tmpPath)

	err = b.extract(db)
	if err != nil {
		return
	}

	cmd := exec.Command(
		"/usr/bin/docker",
		"run",
		"--rm",
		"-v", tmpPath+":/pkg",
		constants.BuildImage,
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		err = &errortypes.ExecError{
			errors.Wrap(err, "build: Failed to get stdout"),
		}
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		err = &errortypes.ExecError{
			errors.Wrap(err, "build: Failed to get stderr"),
		}
		return
	}

	err = cmd.Start()
	if err != nil {
		err = &errortypes.ExecError{
			errors.Wrap(err, "build: Failed to build"),
		}
		return
	}

	output := make(chan *BuildLog, 100)
	outputWait := sync.WaitGroup{}
	outputWait.Add(1)

	go func() {
		defer stdout.Close()
		defer func() {
			output <- nil
		}()

		out := bufio.NewReader(stdout)
		for {
			line, _, err := out.ReadLine()
			if err != nil {
				if err != io.EOF &&
					!strings.Contains(err.Error(), "file already closed") &&
					!strings.Contains(err.Error(), "bad file descriptor") {

					err = &errortypes.ReadError{
						errors.Wrap(err, "build: Failed to read stdout"),
					}
					logrus.WithFields(logrus.Fields{
						"error": err,
					}).Error("build: Stdout error")
				}

				return
			}

			log := &BuildLog{
				Build:     b.Id,
				Timestamp: time.Now(),
				Log:       string(line),
			}

			output <- log
		}
	}()

	go func() {
		defer stderr.Close()

		out := bufio.NewReader(stderr)
		for {
			line, _, err := out.ReadLine()
			if err != nil {
				if err != io.EOF &&
					!strings.Contains(err.Error(), "file already closed") &&
					!strings.Contains(err.Error(), "bad file descriptor") {

					err = &errortypes.ReadError{
						errors.Wrap(err, "build: Failed to read stderr"),
					}
					logrus.WithFields(logrus.Fields{
						"error": err,
					}).Error("build: Stderr error")
				}

				return
			}

			log := &BuildLog{
				Build:     b.Id,
				Timestamp: time.Now(),
				Log:       string(line),
			}

			output <- log
		}
	}()

	go func() {
		defer outputWait.Done()

		for {
			log := <-output
			if log == nil {
				return
			}

			err = coll.Insert(log)
			if err != nil {
				err = database.ParseError(err)
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("build: Output push error")
			}
		}
	}()

	defer outputWait.Wait()

	err = cmd.Wait()
	if err != nil {
		err = &errortypes.ExecError{
			errors.Wrap(err, "build: Build error"),
		}
		return
	}

	files, err := ioutil.ReadDir(tmpPath)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrapf(err, "build: Failed to read dir %s", tmpPath),
		}
		return
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), constants.PackageExt) {
			pkgPath := path.Join(tmpPath, file.Name())

			err = b.addFile(db, pkgPath)
			if err != nil {
				return
			}

			if config.Config.SigKeyName != "" {
				err = signing.SignPackage(pkgPath)
				if err != nil {
					return
				}

				err = b.addFile(db, pkgPath+".sig")
				if err != nil {
					return
				}
			}
		}
	}

	return
}

func (b *Build) Build(db *database.Database) (err error) {
	coll := db.Builds()

	err = coll.Update(&bson.M{
		"_id":   b.Id,
		"state": "pending",
	}, &bson.M{
		"$set": &bson.M{
			"state":      "building",
			"state_rank": BuildingRank,
			"builder":    config.Config.ServerName,
			"start":      time.Now(),
		},
	})
	if err != nil {
		err = database.ParseError(err)

		switch err.(type) {
		case *database.NotFoundError:
			err = nil
		}

		return
	}
	b.State = "building"
	b.StateRank = BuildingRank
	b.Builder = config.Config.ServerName
	PublishChange(db)

	err = b.build(db)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("build: Build failed")

		b.State = "failed"
		b.StateRank = FailedRank
		b.Stop = time.Now()
		coll.CommitFields(b.Id, b, set.NewSet("state", "state_rank", "stop"))
		PublishChange(db)

		return
	}

	b.State = "completed"
	b.StateRank = CompletedRank
	b.Stop = time.Now()
	coll.CommitFields(b.Id, b, set.NewSet("state", "state_rank", "stop"))
	PublishChange(db)

	return
}

func (b *Build) Archive(db *database.Database) (err error) {
	coll := db.Builds()

	err = b.removePkg(db)
	if err != nil {
		return
	}

	err = coll.Update(&bson.M{
		"_id": b.Id,
	}, &bson.M{
		"$set": &bson.M{
			"state":      "archived",
			"state_rank": ArchivedRank,
			"repo_state": nil,
			"pkg_ids":    []bson.ObjectId{},
		},
	})
	if err != nil {
		err = database.ParseError(err)

		switch err.(type) {
		case *database.NotFoundError:
			err = nil
		}

		return
	}
	b.State = "archived"
	b.StateRank = ArchivedRank
	b.RepoState = ""
	PublishChange(db)

	return
}

func (b *Build) Rebuild(db *database.Database) (err error) {
	coll := db.Builds()
	start := time.Now()

	err = b.removePkg(db)
	if err != nil {
		return
	}

	err = coll.Update(&bson.M{
		"_id": b.Id,
	}, &bson.M{
		"$set": &bson.M{
			"state":      "pending",
			"state_rank": PendingRank,
			"builder":    "",
			"uploaded":   false,
			"log":        []string{},
			"start":      start,
			"pkg_ids":    []bson.ObjectId{},
		},
	})
	if err != nil {
		err = database.ParseError(err)

		switch err.(type) {
		case *database.NotFoundError:
			err = nil
		}

		return
	}
	b.State = "pending"
	b.StateRank = PendingRank
	b.Builder = ""
	b.Log = []string{}
	b.Start = start
	PublishChange(db)

	return
}

func (b *Build) removePkg(db *database.Database) (err error) {
	pkgGfs := db.PkgGrid()

	for _, gfId := range b.PkgIds {
		err = pkgGfs.RemoveId(gfId)
		if err != nil {
			err = database.ParseError(err)

			switch err.(type) {
			case *database.NotFoundError:
				err = nil
			}

			return
		}
	}

	b.PkgIds = []bson.ObjectId{}

	return
}

func (b *Build) removePkgBuild(db *database.Database) (err error) {
	pkgBuildGfs := db.PkgBuildGrid()

	if b.PkgBuildId != "" {
		err = pkgBuildGfs.RemoveId(b.PkgBuildId)
		if err != nil {
			err = database.ParseError(err)

			switch err.(type) {
			case *database.NotFoundError:
				err = nil
			}

			return
		}
	}

	b.PkgBuildId = ""

	return
}

func (b *Build) Remove(db *database.Database) (err error) {
	coll := db.Builds()

	err = ClearLog(db, b.Id)
	if err != nil {
		return
	}

	err = b.removePkg(db)
	if err != nil {
		return
	}

	err = b.removePkgBuild(db)
	if err != nil {
		return
	}

	err = coll.RemoveId(b.Id)
	if err != nil {
		err = database.ParseError(err)

		switch err.(type) {
		case *database.NotFoundError:
			err = nil
		}

		return
	}

	PublishChange(db)

	return
}

func (b *Build) Upload(db *database.Database, force bool) (err error) {
	repoPath := b.repoPath()
	coll := db.Builds()
	gfs := db.PkgGrid()

	if b.Uploaded && !force {
		return
	}

	err = utils.ExistsMkdir(repoPath, 0755)
	if err != nil {
		return
	}

	pkgPaths := []string{}

	for _, pkgId := range b.PkgIds {
		gf, e := gfs.OpenId(pkgId)
		if e != nil {
			err = database.ParseError(e)
			return
		}

		pth := path.Join(repoPath, gf.Name())
		file, e := os.Create(pth)
		if e != nil {
			e = &errortypes.WriteError{
				errors.Wrap(err, "build: Failed to open pkg file"),
			}
			return
		}

		_, err = io.Copy(file, gf)
		if err != nil {
			return
		}

		err = file.Close()
		if err != nil {
			err = &errortypes.WriteError{
				errors.Wrap(err, "build: Failed to write pkg file"),
			}
			return
		}

		if strings.HasSuffix(pth, constants.PackageExt) {
			pkgPaths = append(pkgPaths, pth)
		}
	}

	for _, pkgPath := range pkgPaths {
		cmd := exec.Command(
			"/usr/bin/repo-add",
			b.databasePath(),
			pkgPath,
		)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			err = errortypes.ExecError{
				errors.Wrapf(err, "build: Failed to add package"),
			}
			return
		}
	}

	b.Uploaded = true
	coll.CommitFields(b.Id, b, set.NewSet("uploaded"))
	PublishChange(db)

	return
}

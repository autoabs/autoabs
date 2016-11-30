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
	"time"
)

type Build struct {
	Id         bson.ObjectId   `bson:"_id" json:"id"`
	Name       string          `bson:"name" json:"name"`
	Builder    string          `bson:"builder" json:"builder"`
	Start      time.Time       `bson:"start,omitempty" json:"start,omitempty"`
	Stop       time.Time       `bson:"stop,omitempty" json:"stop,omitempty"`
	State      string          `bson:"state" json:"state"`
	Version    string          `bson:"version" json:"version"`
	Release    string          `bson:"release" json:"release"`
	Repo       string          `bson:"core" json:"repo"`
	Arch       string          `bson:"arch" json:"arch"`
	Log        []string        `bson:"log,omitempty" json:"log"`
	PkgIds     []bson.ObjectId `bson:"pkg_ids" json:"-"`
	PkgBuildId bson.ObjectId   `bson:"pkg_build_id" json:"-"`
}

func (b *Build) tmpPath() string {
	return path.Join(config.Config.RootPath, "tmp",
		b.Name+"-"+b.Repo+"-"+b.Arch+"-"+b.Version+"-"+b.Release)
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
			e = &errortypes.ReadError{
				errors.Wrap(err, "build: Failed to read tar header"),
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
			e = &errortypes.WriteError{
				errors.Wrap(err, "build: Failed to write tar file"),
			}
			return
		}

		_, err = io.Copy(file, arc)
		if err != nil {
			return
		}
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
	coll := db.Builds()
	tmpPath := b.tmpPath()

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

	go func() {
		defer stdout.Close()

		out := bufio.NewReader(stdout)
		for {
			line, _, err := out.ReadLine()
			if err != nil {
				if !strings.Contains(
					err.Error(), "bad file descriptor") && err != io.EOF {

					err = &errortypes.ReadError{
						errors.Wrap(err, "build: Failed to read stdout"),
					}
					logrus.WithFields(logrus.Fields{
						"error": err,
					}).Error("build: Stdout error")
				}

				return
			}

			fmt.Println(string(line))

			err = coll.UpdateId(b.Id, &bson.M{
				"$push": &bson.M{
					"log": string(line),
				},
			})
			if err != nil {
				err = database.ParseError(err)
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("build: Stdout push error")
			}
		}
	}()

	go func() {
		defer stderr.Close()

		out := bufio.NewReader(stderr)
		for {
			line, _, err := out.ReadLine()
			if err != nil {
				if !strings.Contains(
					err.Error(), "bad file descriptor") && err != io.EOF {

					err = &errortypes.ReadError{
						errors.Wrap(err, "build: Failed to read stderr"),
					}
					logrus.WithFields(logrus.Fields{
						"error": err,
					}).Error("build: Stderr error")
				}

				return
			}

			fmt.Println(string(line))

			err = coll.UpdateId(b.Id, &bson.M{
				"$push": &bson.M{
					"log": string(line),
				},
			})
			if err != nil {
				err = database.ParseError(err)
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("build: Stderr push error")
			}
		}
	}()

	err = cmd.Start()
	if err != nil {
		err = &errortypes.ExecError{
			errors.Wrap(err, "build: Failed to build"),
		}
		return
	}

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
			"state":   "building",
			"builder": config.Config.ServerName,
			"start":   time.Now(),
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
	b.Builder = config.Config.ServerName

	err = b.build(db)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("build: Build failed")

		b.State = "failed"
		b.Stop = time.Now()
		coll.CommitFields(b.Id, b, set.NewSet("state", "stop"))

		return
	}

	b.State = "completed"
	b.Stop = time.Now()
	coll.CommitFields(b.Id, b, set.NewSet("state", "stop"))

	return
}

func (b *Build) Inactivate(db *database.Database) (err error) {
	coll := db.Builds()

	err = coll.Update(&bson.M{
		"_id":   b.Id,
		"state": "failed",
	}, &bson.M{
		"$set": &bson.M{
			"state":   "inactive",
			"builder": "",
			"start":   time.Now(),
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
	b.State = "inactive"
	b.Builder = ""

	return
}

func (b *Build) Retry(db *database.Database) (err error) {
	coll := db.Builds()

	err = coll.Update(&bson.M{
		"_id":   b.Id,
		"state": "failed",
	}, &bson.M{
		"$set": &bson.M{
			"state":   "pending",
			"builder": "",
			"start":   time.Now(),
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
	b.Builder = ""

	return
}

func (b *Build) Remove(db *database.Database) (err error) {
	pkgGfs := db.PkgGrid()
	pkgBuildGfs := db.PkgBuildGrid()
	coll := db.Builds()

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

	err = coll.RemoveId(b.Id)
	if err != nil {
		err = database.ParseError(err)

		switch err.(type) {
		case *database.NotFoundError:
			err = nil
		}

		return
	}

	return
}

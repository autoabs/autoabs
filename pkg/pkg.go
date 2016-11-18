package pkg

import (
	"archive/tar"
	"bufio"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/autoabs/autoabs/build"
	"github.com/autoabs/autoabs/config"
	"github.com/autoabs/autoabs/constants"
	"github.com/autoabs/autoabs/database"
	"github.com/autoabs/autoabs/errortypes"
	"github.com/autoabs/autoabs/utils"
	"github.com/dropbox/godropbox/errors"
	"gopkg.in/mgo.v2/bson"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

type Package struct {
	Name       string
	SubName    string
	Version    string
	Release    string
	Repo       string
	Arch       string
	Path       string
	SourcePath string
}

func (p *Package) Key() string {
	return p.SubName + "-" + p.Repo + "-" + p.Arch
}

func (p *Package) IdKey() string {
	return p.Name + "-" + p.Repo + "-" + p.Arch
}

func (p *Package) RepoPath() string {
	return path.Join(config.Config.RootPath, "repo", p.Repo, "os", p.Arch)
}

func (p *Package) DatabasePath() string {
	return path.Join(config.Config.RootPath, "repo", p.Repo, "os", p.Arch,
		fmt.Sprintf("%s.db.tar.gz", p.Repo))
}

func (p *Package) BuildPath() string {
	return path.Join(config.Config.RootPath, "builds",
		p.Name+"-"+p.Repo+"-"+p.Arch+"-"+p.Version+"-"+p.Release)
}

func (p *Package) TmpPath() string {
	return path.Join(config.Config.RootPath, "tmp",
		p.Name+"-"+p.Repo+"-"+p.Arch+"-"+p.Version+"-"+p.Release)
}

func (p *Package) LogPath() string {
	return path.Join(config.Config.RootPath, "logs",
		p.Name+"-"+p.Repo+"-"+p.Arch+"-"+p.Version+"-"+p.Release)
}

func (p *Package) QueueBuild() (err error) {
	db := database.GetDatabase()
	defer db.Close()

	coll := db.Builds()
	gfs := db.PkgBuildGrid()

	gf, err := gfs.Create("pkgbuild.tar")
	if err != nil {
		err = database.ParseError(err)
		return
	}

	gf.SetContentType("application/x-tar")
	gfId := gf.Id().(bson.ObjectId)

	bild := &build.Build{
		Id:         bson.NewObjectId(),
		Name:       p.Name,
		State:      "pending",
		Version:    p.Version,
		Release:    p.Release,
		Repo:       p.Repo,
		Arch:       p.Arch,
		PkgIds:     []bson.ObjectId{},
		PkgBuildId: gfId,
	}

	resp, err := coll.Upsert(&bson.M{
		"name":    p.Name,
		"version": p.Version,
		"release": p.Release,
		"repo":    p.Repo,
		"arch":    p.Arch,
	}, &bson.M{
		"$setOnInsert": bild,
	})
	if err != nil {
		err = database.ParseError(err)
		return
	}

	if resp.Matched != 0 {
		return
	}

	arc := tar.NewWriter(gf)

	ln := len(p.SourcePath) + 1
	err = filepath.Walk(p.SourcePath, func(path string,
		info os.FileInfo, err error) (e error) {

		if info.IsDir() {
			return
		}

		if p.SourcePath+"/" != path[:ln] {
			return
		}

		name := path[ln:]

		hdr := &tar.Header{
			Name: name,
			Mode: int64(info.Mode()),
			Size: info.Size(),
		}

		e = arc.WriteHeader(hdr)
		if e != nil {
			e = &errortypes.WriteError{
				errors.Wrap(e, "pkg: Failed to write tar header"),
			}
			return
		}

		file, e := os.Open(path)
		if e != nil {
			e = &errortypes.ReadError{
				errors.Wrap(e, "pkg: Failed to open source file"),
			}
			return
		}
		defer file.Close()

		_, e = io.Copy(arc, file)
		if e != nil {
			e = &errortypes.WriteError{
				errors.Wrap(e, "pkg: Failed to read source file"),
			}
			return
		}

		println(path)

		return
	})

	err = arc.Close()
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrap(err, "pkg: Failed to close tar file"),
		}
		return
	}

	err = gf.Close()
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrap(err, "pkg: Failed to close grid file"),
		}
		return
	}

	return
}

func (p *Package) Build() (err error) {
	buildPath := p.BuildPath()
	tmpPath := p.TmpPath()

	defer utils.ExistsRemove(tmpPath)

	logrus.WithFields(logrus.Fields{
		"package": p.Name,
	}).Info("profile: Building package")

	err = utils.ExistsRemove(tmpPath)
	if err != nil {
		return
	}

	err = utils.CopyAll(buildPath, tmpPath)
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
			errors.Wrap(err, "pkg: Failed to get stdout"),
		}
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		err = &errortypes.ExecError{
			errors.Wrap(err, "pkg: Failed to get stderr"),
		}
		return
	}

	logPath := p.LogPath()
	utils.ExistsRemove(logPath)

	err = utils.ExistsMkdir(logPath, 0755)
	if err != nil {
		return
	}

	logPath = path.Join(logPath, "build.log")
	logLock := sync.Mutex{}
	logFile, err := os.OpenFile(logPath,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrapf(err, "pkg: Failed to open file %s", logPath),
		}
		return
	}
	defer logFile.Close()

	go func() {
		out := bufio.NewReader(stdout)
		for {
			line, _, err := out.ReadLine()
			if err != nil {
				if !strings.Contains(
					err.Error(), "bad file descriptor") && err != io.EOF {

					err = &errortypes.ReadError{
						errors.Wrap(err, "profile: Failed to read stdout"),
					}
					logrus.WithFields(logrus.Fields{
						"error": err,
					}).Error("profile: Stdout error")
				}

				return
			}

			logLock.Lock()
			_, err = logFile.WriteString(string(line) + "\n")
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error":    err,
					"log_path": logPath,
				}).Error("profile: Failed to write file")
			}
			logLock.Unlock()
		}
	}()

	go func() {
		out := bufio.NewReader(stderr)
		for {
			line, _, err := out.ReadLine()
			if err != nil {
				if !strings.Contains(
					err.Error(), "bad file descriptor") && err != io.EOF {

					err = &errortypes.ReadError{
						errors.Wrap(err, "profile: Failed to read stderr"),
					}
					logrus.WithFields(logrus.Fields{
						"error": err,
					}).Error("profile: Stderr error")
				}

				return
			}

			logLock.Lock()
			_, err = logFile.WriteString(string(line) + "\n")
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error":    err,
					"log_path": logPath,
				}).Error("profile:  Failed to write file")
			}
			logLock.Unlock()
		}
	}()

	err = cmd.Start()
	if err != nil {
		err = &errortypes.ExecError{
			errors.Wrap(err, "pkg: Failed to build"),
		}
		return
	}

	err = cmd.Wait()
	if err != nil {
		err = &errortypes.ExecError{
			errors.Wrap(err, "pkg: Build error"),
		}
		return
	}

	files, err := ioutil.ReadDir(tmpPath)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrapf(err, "pkg: Failed to read dir %s", tmpPath),
		}
		return
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), constants.PackageExt) {
			err = p.Add(path.Join(tmpPath, file.Name()))
			if err != nil {
				return
			}
		}
	}

	utils.ExistsRemove(buildPath)

	return
}

func (p *Package) Add(pkgPath string) (err error) {
	repoPath := p.RepoPath()

	err = utils.ExistsMkdir(repoPath, 0755)
	if err != nil {
		return
	}

	err = utils.Copy(pkgPath, p.RepoPath())
	if err != nil {
		return
	}

	cmd := exec.Command(
		"/usr/bin/repo-add",
		p.DatabasePath(),
		pkgPath,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		err = errortypes.ExecError{
			errors.Wrapf(err, "package: Failed to add package"),
		}
		return
	}

	return
}

func (p *Package) Remove() {
	if p.Path == "" {
		return
	}

	cmd := exec.Command(
		"/usr/bin/repo-remove",
		p.DatabasePath(),
		p.Path,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	os.Remove(p.Path)
	os.Remove(p.Path + ".sig")

	return
}

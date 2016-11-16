package pkg

import (
	"bufio"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/autoabs/autoabs/config"
	"github.com/autoabs/autoabs/constants"
	"github.com/autoabs/autoabs/errortypes"
	"github.com/autoabs/autoabs/utils"
	"github.com/dropbox/godropbox/errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
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

func (p *Package) Print() {
	fmt.Printf("Id: %s\n", p.Name)
	fmt.Printf("  Name: %s\n", p.SubName)
	fmt.Printf("  Version: %s\n", p.Version)
	fmt.Printf("  Release: %s\n", p.Release)
	fmt.Printf("  Repo: %s\n", p.Repo)
	fmt.Printf("  Arch: %s\n", p.Arch)
	fmt.Printf("  Path: %s\n", p.Path)
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
		p.Id+"-"+p.Repo+"-"+p.Arch+"-"+p.Version+"-"+p.Release)
}

func (p *Package) TmpPath() string {
	return path.Join(config.Config.RootPath, "tmp",
		p.Id+"-"+p.Repo+"-"+p.Arch+"-"+p.Version+"-"+p.Release)
}

func (p *Package) LogPath() string {
	return path.Join(config.Config.RootPath, "logs",
		p.Id+"-"+p.Repo+"-"+p.Arch+"-"+p.Version+"-"+p.Release)
}

func (p *Package) QueueBuild() (err error) {
	buildPath := p.BuildPath()

	err = utils.ExistsRemove(buildPath)
	if err != nil {
		return
	}

	err = utils.CopyAll(p.SourcePath, buildPath)
	if err != nil {
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

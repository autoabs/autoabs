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
)

type Package struct {
	Id         string
	Name       string
	Version    string
	Release    string
	Repo       string
	Arch       string
	Path       string
	SourcePath string
}

func (p *Package) Key() string {
	return p.Name + "-" + p.Repo + "-" + p.Arch
}

func (p *Package) IdKey() string {
	return p.Id + "-" + p.Repo + "-" + p.Arch
}

func (p *Package) Print() {
	fmt.Printf("Id: %s\n", p.Id)
	fmt.Printf("  Name: %s\n", p.Name)
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

			fmt.Println(string(line))
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

			fmt.Println(string(line))
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

	return
}

func (p *Package) Add(pkgPath string) (err error) {
	err = utils.Copy(pkgPath, p.RepoPath())
	if err != nil {
		return
	}

	err = exec.Command(
		"/usr/bin/repo-add",
		p.DatabasePath(),
		pkgPath,
	).Run()
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

	exec.Command(
		"/usr/bin/repo-remove",
		p.DatabasePath(),
		p.Path,
	).Run()

	os.Remove(p.Path)
	os.Remove(p.Path + ".sig")

	return
}

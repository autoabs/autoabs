package pkg

import (
	"fmt"
	"github.com/autoabs/autoabs/config"
	"github.com/autoabs/autoabs/errortypes"
	"github.com/autoabs/autoabs/utils"
	"github.com/dropbox/godropbox/errors"
	"os"
	"os/exec"
	"path"
)

type Package struct {
	Name       string
	Version    string
	Release    string
	Repo       string
	Arch       string
	Path       string
	SourcePath string
	Previous   *Package
}

func (p *Package) Key() string {
	return p.Name + "-" + p.Repo + "-" + p.Arch
}

func (p *Package) Print() {
	fmt.Printf("%s %s-%s: %s %s %s\n",
		p.Name, p.Version, p.Release, p.Repo, p.Arch, p.Path)
}

func (p *Package) RepoPath() string {
	return path.Join(config.Config.RootPath, "repo", p.Repo, "os", p.Arch)
}

func (p *Package) DatabasePath() string {
	return path.Join(config.Config.RootPath, "repo", p.Repo, "os", p.Arch,
		fmt.Sprintf("%s.db.tar.gz", p.Repo))
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

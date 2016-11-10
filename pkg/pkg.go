package pkg

import (
	"fmt"
	"github.com/autoabs/autoabs/config"
	"github.com/autoabs/autoabs/errortypes"
	"github.com/dropbox/godropbox/errors"
	"os"
	"os/exec"
	"path"
)

type Package struct {
	Name     string
	Version  string
	Release  string
	Repo     string
	Arch     string
	Path     string
	Previous *Package
}

func (p *Package) Key() string {
	return p.Name + "-" + p.Repo + "-" + p.Arch
}

func (p *Package) Print() {
	fmt.Printf("%s %s-%s: %s %s %s\n",
		p.Name, p.Version, p.Release, p.Repo, p.Arch, p.Path)
}

func (p *Package) RepoPath() string {
	return path.Join(config.Config.RootPath, "repo", p.Repo, "os", p.Arch,
		fmt.Sprintf("%s.db.tar.gz", p.Repo))
}

func (p *Package) Remove() (err error) {
	if p.Path == "" {
		return
	}

	err = exec.Command(
		"repo-remove",
		p.RepoPath(),
		p.Path,
	).Run()
	if err != nil {
		err = errortypes.ExecError{
			errors.Wrapf(err, "package: Failed to remove package"),
		}
		return
	}

	os.Remove(p.Path)
	os.Remove(p.Path + ".sig")

	return
}

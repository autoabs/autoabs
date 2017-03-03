package pkg

import (
	"fmt"
	"github.com/autoabs/autoabs/config"
	"github.com/autoabs/autoabs/database"
	"gopkg.in/mgo.v2/bson"
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
}

func (p *Package) Key() string {
	return p.Name + "-" + p.Repo + "-" + p.Arch
}

func (p *Package) FullVersion() string {
	return p.Version + "-" + p.Release
}

func (p *Package) DatabasePath() string {
	return path.Join(config.Config.RootPath, "repo", p.Repo, "os", p.Arch,
		fmt.Sprintf("%s.db.tar.gz", p.Repo))
}

func (p *Package) SyncState(db *database.Database, stateId bson.ObjectId) (
	err error) {

	coll := db.Builds()

	err = coll.Update(&bson.M{
		"sub_names": p.Name,
		"version":   p.Version,
		"release":   p.Release,
		"repo":      p.Repo,
		"arch":      p.Arch,
	}, &bson.M{
		"$set": &bson.M{
			"repo_state": stateId,
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

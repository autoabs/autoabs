package pkg

import (
	"archive/tar"
	"fmt"
	"github.com/autoabs/autoabs/build"
	"github.com/autoabs/autoabs/config"
	"github.com/autoabs/autoabs/database"
	"github.com/autoabs/autoabs/errortypes"
	"github.com/dropbox/godropbox/errors"
	"gopkg.in/mgo.v2/bson"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

type Package struct {
	Name       string
	SubName    string
	SubNames   []string
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

func (p *Package) LogPath() string {
	return path.Join(config.Config.RootPath, "logs",
		p.Name+"-"+p.Repo+"-"+p.Arch+"-"+p.Version+"-"+p.Release)
}

func (p *Package) QueueBuild(db *database.Database, force bool) (err error) {
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
		SubNames:   p.SubNames,
		State:      "pending",
		Version:    p.Version,
		Release:    p.Release,
		Repo:       p.Repo,
		Arch:       p.Arch,
		PkgIds:     []bson.ObjectId{},
		PkgBuildId: gfId,
	}

	if force {
		err = coll.Insert(bild)
		if err != nil {
			err = database.ParseError(err)
			return
		}
	} else {
		resp, e := coll.Upsert(&bson.M{
			"name":    p.Name,
			"version": p.Version,
			"release": p.Release,
			"repo":    p.Repo,
			"arch":    p.Arch,
		}, &bson.M{
			"$setOnInsert": bild,
		})
		if e != nil {
			err = database.ParseError(e)
			return
		}

		if resp.Matched != 0 {
			return
		}
	}

	arc := tar.NewWriter(gf)

	ln := len(p.SourcePath) + 1
	err = filepath.Walk(p.SourcePath, func(path string,
		info os.FileInfo, err error) (e error) {

		if err != nil {
			e = &errortypes.WriteError{
				errors.Wrap(err, "pkg: Failed to read pkg directory"),
			}
			return
		}

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

		return
	})
	if err != nil {
		return
	}

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

func (p *Package) SyncState(db *database.Database, stateId bson.ObjectId) (
	err error) {

	coll := db.Builds()

	err = coll.Update(&bson.M{
		"sub_names": p.SubName,
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

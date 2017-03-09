package source

import (
	"archive/tar"
	"github.com/autoabs/autoabs/build"
	"github.com/autoabs/autoabs/database"
	"github.com/autoabs/autoabs/errortypes"
	"github.com/dropbox/godropbox/errors"
	"gopkg.in/mgo.v2/bson"
	"io"
	"os"
	"path/filepath"
)

type Source struct {
	Name     string
	SubNames []string
	Version  string
	Release  string
	Repo     string
	Arch     string
	Path     string
}

func (s *Source) Keys() []string {
	keys := []string{}

	for _, subName := range s.SubNames {
		keys = append(keys, subName+"-"+s.Repo+"-"+s.Arch)
	}

	return keys
}

func (s *Source) Queue(db *database.Database, force bool) (
	queued bool, err error) {

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
		Name:       s.Name,
		SubNames:   s.SubNames,
		State:      "pending",
		StateRank:  build.PendingRank,
		Version:    s.Version,
		Release:    s.Release,
		Repo:       s.Repo,
		Arch:       s.Arch,
		PkgIds:     []bson.ObjectId{},
		PkgBuildId: gfId,
	}

	resp, e := coll.Upsert(&bson.M{
		"name":    s.Name,
		"version": s.Version,
		"release": s.Release,
		"repo":    s.Repo,
		"arch":    s.Arch,
	}, &bson.M{
		"$setOnInsert": bild,
	})
	if e != nil {
		err = database.ParseError(e)
		return
	}

	if resp.Matched != 0 {
		if force {
			bildCur, e := build.GetKey(db, s.Name, s.Version,
				s.Release, s.Repo, s.Arch)
			if err != nil {
				switch e.(type) {
				case *database.NotFoundError:
					e = nil
				default:
					err = e
					return
				}
			}

			if bildCur.State == "completed" {
				err = bildCur.Upload(db, true)
				if err != nil {
					return
				}
			}
		}

		return
	}

	queued = true

	arc := tar.NewWriter(gf)

	ln := len(s.Path) + 1
	err = filepath.Walk(s.Path, func(path string,
		info os.FileInfo, err error) (e error) {

		if err != nil {
			e = &errortypes.WriteError{
				errors.Wrap(err, "source: Failed to read pkg directory"),
			}
			return
		}

		if info.IsDir() {
			return
		}

		if s.Path+"/" != path[:ln] {
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
				errors.Wrap(e, "source: Failed to write tar header"),
			}
			return
		}

		file, e := os.Open(path)
		if e != nil {
			e = &errortypes.ReadError{
				errors.Wrap(e, "source: Failed to open source file"),
			}
			return
		}
		defer file.Close()

		_, e = io.Copy(arc, file)
		if e != nil {
			e = &errortypes.WriteError{
				errors.Wrap(e, "source: Failed to read source file"),
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
			errors.Wrap(err, "source: Failed to close tar file"),
		}
		return
	}

	err = gf.Close()
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrap(err, "source: Failed to close grid file"),
		}
		return
	}

	return
}

package build

import (
	"archive/tar"
	"github.com/autoabs/autoabs/config"
	"github.com/autoabs/autoabs/database"
	"github.com/autoabs/autoabs/errortypes"
	"github.com/dropbox/godropbox/errors"
	"github.com/pacur/pacur/utils"
	"gopkg.in/mgo.v2/bson"
	"io"
	"os"
	"path"
	"time"
)

type Build struct {
	Id         bson.ObjectId   `bson:"_id"`
	Name       string          `bson:"name"`
	Start      time.Time       `bson:"start,omitempty"`
	Stop       time.Time       `bson:"stop,omitempty"`
	State      string          `bson:"state"`
	Version    string          `bson:"version"`
	Release    string          `bson:"release"`
	Repo       string          `bson:"core"`
	Arch       string          `bson:"arch"`
	PkgIds     []bson.ObjectId `bson:"pkg_ids"`
	PkgBuildId bson.ObjectId   `bson:"pkg_build_id"`
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

		err = utils.ExistsMakeDir(dirPth)
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

func (b *Build) Build(db *database.Database) (err error) {
	err = b.extract(db)
	if err != nil {
		return
	}

	return
}

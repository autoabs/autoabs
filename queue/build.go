package queue

import (
	"github.com/autoabs/autoabs/config"
	"github.com/autoabs/autoabs/errortypes"
	"github.com/autoabs/autoabs/pkg"
	"github.com/autoabs/autoabs/utils"
	"github.com/dropbox/godropbox/errors"
	"io/ioutil"
	"path"
	"strings"
)

func getBuildPackages() (pkgs []*pkg.Package, err error) {
	pth := path.Join(config.Config.RootPath, "builds")
	pkgs = []*pkg.Package{}

	exists, err := utils.ExistsDir(pth)
	if err != nil {
		return
	}

	if !exists {
		return
	}

	builds, err := ioutil.ReadDir(pth)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrapf(err, "queue: Failed to read dir %s", pth),
		}
		return
	}

	for _, entry := range builds {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		nameSpl := strings.Split(name, "-")
		ln := len(nameSpl)

		if ln < 4 {
			continue
		}

		pk := &pkg.Package{
			Name:    strings.Join(nameSpl[:ln-4], "-"),
			Version: nameSpl[ln-2],
			Release: nameSpl[ln-1],
			Repo:    nameSpl[ln-4],
			Arch:    nameSpl[ln-3],
		}

		pkgs = append(pkgs, pk)
	}

	return
}

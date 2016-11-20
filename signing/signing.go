package signing

import (
	"github.com/autoabs/autoabs/config"
	"github.com/autoabs/autoabs/utils"
)

func SignPackage(pkgPath string) (err error) {
	err = utils.Exec("", "gpg",
		"--detach-sign",
		"-u", config.Config.SigKeyName,
		"--no-armor",
		pkgPath)
	if err != nil {
		return
	}

	return
}

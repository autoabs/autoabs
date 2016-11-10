package utils

import (
	"github.com/autoabs/autoabs/errortypes"
	"github.com/dropbox/godropbox/errors"
	"os"
	"os/exec"
)

func Exists(pth string) (exists bool, err error) {
	_, err = os.Stat(pth)
	if err == nil {
		exists = true
		return
	}

	if os.IsNotExist(err) {
		err = nil
		return
	}

	err = errortypes.ReadError{
		errors.Wrapf(err, "utils: Failed to stat %s", pth),
	}
	return
}

func ExistsDir(pth string) (exists bool, err error) {
	stat, err := os.Stat(pth)
	if err == nil {
		exists = stat.IsDir()
		return
	}

	if os.IsNotExist(err) {
		err = nil
		return
	}

	err = errortypes.ReadError{
		errors.Wrapf(err, "utils: Failed to stat %s", pth),
	}
	return
}

func ExistsFile(pth string) (exists bool, err error) {
	stat, err := os.Stat(pth)
	if err == nil {
		exists = !stat.IsDir()
		return
	}

	if os.IsNotExist(err) {
		err = nil
		return
	}

	err = errortypes.ReadError{
		errors.Wrapf(err, "utils: Failed to stat %s", pth),
	}
	return
}

func ExistsMkdir(pth string, perm os.FileMode) (err error) {
	exists, err := ExistsDir(pth)
	if err != nil {
		return
	}

	if !exists {
		err = os.MkdirAll(pth, perm)
		if err != nil {
			err = &errortypes.WriteError{
				errors.Wrapf(err, "utils: Failed to mkdir %s", pth),
			}
			return
		}
	}

	return
}

func ExistsRemove(pth string) (err error) {
	exists, err := Exists(pth)
	if err != nil {
		return
	}

	if exists {
		err = os.RemoveAll(pth)
		if err != nil {
			err = &errortypes.WriteError{
				errors.Wrapf(err, "utils: Failed to rm %s", pth),
			}
			return
		}
	}

	return
}

func Copy(sourcePath, destPath string) (err error) {
	err = exec.Command(
		"/usr/bin/cp",
		sourcePath,
		destPath,
	).Run()
	if err != nil {
		err = errortypes.ExecError{
			errors.Wrapf(err, "package: Failed to copy %s to %s",
				sourcePath, destPath),
		}
		return
	}

	return
}

func CopyAll(sourcePath, destPath string) (err error) {
	err = exec.Command(
		"/usr/bin/cp",
		"-r",
		sourcePath,
		destPath,
	).Run()
	if err != nil {
		err = errortypes.ExecError{
			errors.Wrapf(err, "package: Failed to copy %s to %s",
				sourcePath, destPath),
		}
		return
	}

	return
}

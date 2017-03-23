package utils

import (
	"github.com/dropbox/godropbox/container/set"
	"strings"
)

func GitCommit(pth string) (commit string, err error) {
	output, err := ExecOutput(pth, "git", "rev-parse", "HEAD")
	if err != nil {
		return
	}

	commit = strings.TrimSpace(output)

	return
}

func GitChanged(pth, commitX, commitY string) (changed set.Set, err error) {
	changed = set.NewSet()

	output, err := ExecOutput(pth, "git", "diff",
		"--name-only", commitX, commitY)
	if err != nil {
		return
	}

	output = strings.TrimSpace(output)

	for _, name := range strings.Split(output, "\n") {
		changed.Add(strings.TrimSpace(name))
	}

	return
}

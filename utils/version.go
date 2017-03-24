package utils

import (
	"os/exec"
	"strconv"
	"strings"
)

func VersionNewer(x, y string) bool {
	cmd := exec.Command("vercmp", x, y)

	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	n, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		panic(err)
	}

	if n >= 0 {
		return true
	}

	return false
}

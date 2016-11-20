package signing

import (
	"fmt"
	"github.com/autoabs/autoabs/errortypes"
	"github.com/autoabs/autoabs/utils"
	"github.com/dropbox/godropbox/errors"
	"strings"
)

func GetId(name, email string) (id string, err error) {
	uidStr := fmt.Sprintf("%s <%s>", name, email)

	output, err := utils.ExecOutput("", "gpg", "--list-keys")
	if err != nil {
		return
	}

	curId := ""

	hasKey := false
	for _, line := range strings.Split(output, "\n") {
		if curId != "" {
			if strings.HasPrefix(line, "uid") {
				if strings.Contains(line, uidStr) {
					break
				} else {
					curId = ""
				}
			}
			continue
		}

		if hasKey {
			curId = strings.TrimSpace(line)
			hasKey = false
			continue
		}

		if !strings.HasPrefix(line, "pub") {
			continue
		}

		if !strings.Contains(line, "/") {
			hasKey = true
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		split := strings.Split(fields[1], "/")
		if len(split) < 2 {
			continue
		}

		curId = split[1]
		break
	}

	if curId == "" {
		err = &errortypes.NotFoundError{
			errors.New("signing: Failed to find gpg id"),
		}
		return
	}

	id = curId

	return
}

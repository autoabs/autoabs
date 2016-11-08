package logger

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/autoabs/autoabs/config"
	"github.com/autoabs/autoabs/errortypes"
	"github.com/dropbox/godropbox/errors"
	"os"
	"path"
	"sync"
)

var fileLock = sync.Mutex{}

type fileSender struct{}

func (s *fileSender) Init() {}

func (s *fileSender) Parse(entry *logrus.Entry) {
	msg := formatPlain(entry)

	fileLock.Lock()
	defer fileLock.Unlock()

	path := path.Join(config.Config.RootPath, "system.log")

	file, err := os.OpenFile(path, os.O_CREATE|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrap(err, "logger: Failed to write entry"),
		}
		fmt.Println(err.Error())
		return
	}
	defer file.Close()

	file.Write(msg)
}

func init() {
	senders = append(senders, &fileSender{})
}

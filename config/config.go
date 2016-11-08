package config

import (
	"encoding/json"
	"github.com/autoabs/autoabs/errortypes"
	"github.com/dropbox/godropbox/errors"
	"io/ioutil"
	"os"
)

const (
	confPath          = "/etc/autoabs.json"
	rootPathDefault   = "/autoabs"
	serverPortDefault = 80
	serverHostDefault = "0.0.0.0"
)

var Config = &ConfigData{}

type ConfigData struct {
	path       string `json:"path"`
	loaded     bool   `json:"-"`
	RootPath   string `json:"root_path"`
	ServerPort int    `json:"server_port"`
	ServerHost string `json:"server_host"`
}

func (c *ConfigData) Load(path string) (err error) {
	c.path = path

	_, err = os.Stat(c.path)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		} else {
			err = &errortypes.ReadError{
				errors.Wrap(err, "config: File stat error"),
			}
		}
		return
	}

	file, err := ioutil.ReadFile(c.path)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "config: File read error"),
		}
		return
	}

	err = json.Unmarshal(file, Config)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "config: File unmarshal error"),
		}
		return
	}

	if c.RootPath == "" {
		c.RootPath = rootPathDefault
	}

	if c.ServerPort == 0 {
		c.ServerPort = serverPortDefault
	}

	if c.ServerHost == "" {
		c.ServerHost = serverHostDefault
	}

	c.loaded = true

	return
}

func (c *ConfigData) Save() (err error) {
	if !c.loaded {
		err = &errortypes.WriteError{
			errors.New("config: Config file has not been loaded"),
		}
		return
	}

	data, err := json.Marshal(c)
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrap(err, "config: File marshal error"),
		}
		return
	}

	err = ioutil.WriteFile(c.path, data, 0600)
	if err != nil {
		err = &errortypes.WriteError{
			errors.Wrap(err, "config: File write error"),
		}
		return
	}

	return
}

func Load() (err error) {
	err = Config.Load(confPath)
	if err != nil {
		return
	}

	return
}

func Save() (err error) {
	err = Config.Save()
	if err != nil {
		return
	}

	return
}

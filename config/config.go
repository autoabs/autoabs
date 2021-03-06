package config

import (
	"encoding/json"
	"github.com/autoabs/autoabs/errortypes"
	"github.com/autoabs/autoabs/requires"
	"github.com/autoabs/autoabs/utils"
	"github.com/dropbox/godropbox/errors"
	"io/ioutil"
	"os"
)

var (
	confPath          = "/etc/autoabs.json"
	rootPathDefault   = "/autoabs"
	mongoUriDefault   = "mongodb://localhost:27017/autoabs"
	serverPortDefault = 9600
	serverHostDefault = "0.0.0.0"
)

var Config = &ConfigData{}

type ConfigData struct {
	path          string `json:"-"`
	loaded        bool   `json:"-"`
	RootPath      string `json:"root_path"`
	MongoUri      string `json:"mongo_uri"`
	ServerName    string `json:"server_name"`
	ServerPort    int    `json:"server_port"`
	ServerHost    string `json:"server_host"`
	SigKeyName    string `json:"sig_key_name"`
	WebNodeId     string `json:"web_node_id"`
	StorageNodeId string `json:"storage_node_id"`
	BuilderNodeId string `json:"builder_node_id"`
}

func (c *ConfigData) Load(path string) (err error) {
	c.path = path

	_, err = os.Stat(c.path)
	if err != nil {
		if os.IsNotExist(err) {
			c.loaded = true
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

	if c.MongoUri == "" {
		c.MongoUri = mongoUriDefault
	}

	if c.ServerName == "" {
		c.ServerName = utils.RandName()
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

func init() {
	module := requires.New("config")

	module.Handler = func() {
		err := Load()
		if err != nil {
			panic(err)
		}

		save := false

		if Config.WebNodeId == "" {
			save = true
			Config.WebNodeId = utils.RandName()
		}

		if Config.StorageNodeId == "" {
			save = true
			Config.StorageNodeId = utils.RandName()
		}

		if Config.BuilderNodeId == "" {
			save = true
			Config.BuilderNodeId = utils.RandName()
		}

		if save {
			err = Save()
			if err != nil {
				panic(err)
			}
		}
	}
}

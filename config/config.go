package config

import (
	"encoding/json"
	"github.com/autoabs/autoabs/errortypes"
	"github.com/autoabs/autoabs/requires"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"io/ioutil"
	"os"
)

var (
	confPath           = "/etc/autoabs.json"
	rootPathDefault    = "/autoabs"
	mongoUriDefault    = "mongodb://localhost:27017/autoabs"
	serverPortDefault  = 9600
	serverHostDefault  = "0.0.0.0"
	targetReposDefault = set.NewSet(
		"community",
		"core",
		"extra",
		"multilib",
	)
	targetArchsDefault = set.NewSet(
		"any",
		"x86_64",
	)
)

var Config = &ConfigData{}

type ConfigData struct {
	path        string   `json:"path"`
	loaded      bool     `json:"-"`
	RootPath    string   `json:"root_path"`
	MongoUri    string   `json:"mongo_uri"`
	ServerPort  int      `json:"server_port"`
	ServerHost  string   `json:"server_host"`
	targetRepos []string `json:"target_repos"`
	TargetRepos set.Set  `json:"-"`
	targetArchs []string `json:"target_archs"`
	TargetArchs set.Set  `json:"-"`
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

	if c.MongoUri == "" {
		c.MongoUri = mongoUriDefault
	}

	if c.ServerPort == 0 {
		c.ServerPort = serverPortDefault
	}

	if c.ServerHost == "" {
		c.ServerHost = serverHostDefault
	}

	if len(c.targetRepos) == 0 {
		c.TargetRepos = targetReposDefault.Copy()
	} else {
		c.TargetRepos = set.NewSet()
		for _, item := range c.targetRepos {
			c.TargetRepos.Add(item)
		}
	}

	if len(c.targetArchs) == 0 {
		c.TargetArchs = targetArchsDefault.Copy()
	} else {
		c.TargetArchs = set.NewSet()
		for _, item := range c.targetArchs {
			c.TargetRepos.Add(item)
		}
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

	c.targetRepos = []string{}
	for item := range c.TargetRepos.Iter() {
		c.targetRepos = append(c.targetRepos, item.(string))
	}

	c.targetArchs = []string{}
	for item := range c.TargetArchs.Iter() {
		c.targetArchs = append(c.targetArchs, item.(string))
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
	}
}

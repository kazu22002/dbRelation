package config

import (
	"errors"
	"github.com/BurntSushi/toml"
	"io/ioutil"
)
func New(flags *Flags) (*Context, error) {
	context := new(Context)
	context.Config = new(Config)
	// read config file
	err := context.initConfigFromFile(flags.CfgFileName)
	if err != nil {
		return nil, err
	}

	return context, err
}

type Context struct {
	Config *Config
	Flags  *Flags
}

type Flags struct {
	CfgFileName string
}

func (ct *Context) initConfigFromFile(cfgFileName string) error {
	// read config
	tomlData, err := ioutil.ReadFile(cfgFileName)
	if err != nil {
		return errors.New("Configuration file read error: " + cfgFileName + "\nError:" + err.Error())
	}
	_, err = toml.Decode(string(tomlData[:]), &ct.Config)
	if err != nil {
		return errors.New("Configuration file decoding error: " + cfgFileName + "\nError:" + err.Error())
	}

	return nil
}

type Config struct {
	Database struct {
		Db  string `toml:"db"`
		Dsn string `toml:"dsn"`
	} `toml:"database"`
}
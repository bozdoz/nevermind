package common

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
)

type config struct {
	Current Version `json:"current"`
}

// config filename
const CONFIG_NAME = "config.json"

// config file permissions
const CONFIG_PERM = 0644

// overridden in tests.
var (
	getNvmDir  = GetNVMDir
	fileOpener = openFile
)

func openFile(name string, flag int, perm fs.FileMode) (io.ReadWriteCloser, error) {
	return os.OpenFile(name, flag, perm)
}

func GetConfig() (cfg config, err error) {
	cfg = config{}
	configFile, err := getNvmDir(CONFIG_NAME)

	if err != nil {
		return
	}

	file, err := fileOpener(configFile, os.O_CREATE|os.O_RDONLY, CONFIG_PERM)

	if err != nil {
		return
	}

	err = json.NewDecoder(file).Decode(&cfg)

	defer file.Close()

	if err != nil && (errors.Is(err, io.EOF) || os.IsNotExist(err)) {
		// can't decode empty file; start fresh
		cfg = config{}
		err = nil
	}

	return
}

func SetConfig(cfg config) error {
	configFile, err := getNvmDir(CONFIG_NAME)

	if err != nil {
		return err
	}

	log.Println("set config", configFile, cfg)

	file, err := fileOpener(configFile, os.O_CREATE|os.O_WRONLY, CONFIG_PERM)

	if err != nil {
		return err
	}

	err = json.NewEncoder(file).Encode(&cfg)

	return err
}

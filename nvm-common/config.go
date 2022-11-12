package common

import (
	"encoding/json"
	"log"
	"os"
)

type config struct {
	Current Version `json:"current"`
}

func GetConfig() (cfg config, err error) {
	var file *os.File

	configFile := GetNVMDir("config.json")

	log.Println("config:", configFile)

	// check if config file exists
	_, err = os.Stat(configFile)

	if err == nil {
		log.Println("config exists")
	} else {
		log.Println("creating config")
	}

	file, err = os.OpenFile(configFile, os.O_CREATE|os.O_RDONLY, 0644)

	if err != nil {
		return
	}

	err = json.NewDecoder(file).Decode(&cfg)

	defer file.Close()

	if err != nil && (err.Error() == "EOF" || os.IsNotExist(err)) {
		// can't decode empty file; start fresh
		cfg = config{}
		err = nil
	}

	log.Println("config", cfg)

	return
}

func SetConfig(cfg config) error {
	configFile := GetNVMDir("config.json")

	log.Println("set config", configFile, cfg)

	file, err := os.OpenFile(configFile, os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	err = json.NewEncoder(file).Encode(&cfg)

	return err
}

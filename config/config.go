package config

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type ConfigInFile struct {
	Zoomeye Zoomeye `toml:"zoomeye"`
}

type Zoomeye struct {
	APIKey string `toml:"apikey"`
}

const (
	configEnvKey = "GTOO_CONFIG"
)

// ReadConfig 读取cfg
func ReadConfig() (*ConfigInFile, error) {
	fn, err := getConfigFilename()
	if err != nil {
		log.Println("Error when get config file.")
		log.Println("You can define config file with env GTOO_CONFIG or just write to ~/.gtoo_config.toml")
		return nil, err
	}
	log.Printf("Use config file: %s", fn)
	var cfg ConfigInFile
	_, err = toml.DecodeFile(fn, &cfg)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &cfg, nil
}

func getConfigFilename() (string, error) {
	fp := os.Getenv(configEnvKey)
	if len(fp) > 0 {
		return fp, nil
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fp = fmt.Sprintf("%s/.gtoo_config.toml", homeDir)
	_, err = os.Stat(fp)
	if err != nil {
		// TODO
		// 这里将文件内容写到~/.gtoo_config.toml里面
		return "", err
	}
	return fp, nil
}

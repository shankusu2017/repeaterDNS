package config

import (
	"io/ioutil"
	"os"
)

func InitConfig(path string) map[string]bool {
	if path == "" {
		path = "./config/localDomain.conf"
	}

	os.Open(path)
	ioutil
	ioutil.ReadFile(path)
}

package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	RepeaterSrv  string `json:"RepeaterSrv"` // repeaterSrv
	RepeaterPort int    `json:"RepeaterPort"`
}

func Init(cfg *ServerConfig) {
	srvCfg := selectConfigFile()
	file, err := ioutil.ReadFile(srvCfg)
	if err != nil {
		log.Fatalf("f6918311 read %s fail, err(%s)\n", srvCfg, err)
	}
	err = json.Unmarshal(file[:], cfg)
	if err != nil {
		log.Fatalf("f6918311 Unmarshal [%v] fail, err(%s)\n", string(file[:]), err)
	}
}

func GetPublicDNS() []string {
	dns := []string{"8.8.8.8", "8.8.4.4", "1.1.1.1", "199.85.126.10",
		"199.85.127.10", "208.67.222.222", "208.67.220.220", "84.200.69.80",
		"84.200.70.40", "8.26.56.26", "8.20.247.20", "64.6.64.6",
		"64.6.65.6", "192.95.54.3", "192.95.54.1", " 81.218.119.11",
		"209.88.198.133"}
	return dns
}

func (srv *ServerConfig) IsConnectMode() bool {
	return srv.RepeaterPort == 53
}

func (srv *ServerConfig) IsProxy() bool {
	return !srv.IsConnectMode()
}

func (srv *ServerConfig) GetRepeaterPort() int {
	return srv.RepeaterPort
}

func (srv *ServerConfig) GetRepeaterSrv() string {
	return srv.RepeaterSrv
}

func selectConfigFile() string {
	paths := []string{"config/server.cfg.first", "config/server.cfg", "config/server.cfg.bak"}
	for _, path := range paths {
		_, err := os.Stat(path) //os.Stat获取文件信息
		if err == nil {
			log.Printf("058de794 config file: %s\n", path)
			return path
		} else {
			if os.IsNotExist(err) {
				continue
			}
		}
	}

	log.Fatalf("FATAL 2806d49c config file is nil")
	return ""
}

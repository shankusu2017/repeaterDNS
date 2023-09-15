package config

import (
	"encoding/hex"
	"encoding/json"
	"github.com/shankusu2017/utils"
	"io/ioutil"
	"log"
	"os"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	ServerMode   string `json:"ServerMode"`
	RepeaterSrv  string `json:"RepeaterSrv"` // repeaterSrv
	RepeaterPort int    `json:"RepeaterPort"`
	Cache        bool   `json:"Cache"`
}

var (
	mainCfg *ServerConfig
)

const (
	DebugFlag = false
	DebugPAC  = true
)

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
	mainCfg = cfg
}

func GetRepeaterSrvAddr() (string, int) {
	return mainCfg.GetRepeaterSrv(), mainCfg.GetRepeaterPort()
}

func GetRandomPublicDNS() string {
	return utils.SliceRandOne(publicDNS())
}

func IsCache() bool {
	return mainCfg.Cache
}

func publicDNS() []string {
	dns := []string{"8.8.8.8", "8.8.4.4", "1.1.1.1", "199.85.126.10",
		"199.85.127.10", "208.67.222.222", "208.67.220.220", "84.200.69.80",
		"84.200.70.40", "8.26.56.26", "8.20.247.20", "64.6.64.6",
		"64.6.65.6", "192.95.54.3", "192.95.54.1", "81.218.119.11",
		"209.88.198.133"}
	return dns
}

func GetIV16() []byte {
	iv, _ := hex.DecodeString("daca11ed1f3fc59b2d233ec67cc6f028")
	return iv[:]
}

func GetKey16() []byte {
	key, _ := hex.DecodeString("ac5299e1424c188fdb618ee0ee5481f7")
	return key[:]
}

func (srv *ServerConfig) IsConnectMode() bool {
	return srv.ServerMode != "proxy"
}

func (srv *ServerConfig) IsProxy() bool {
	return !srv.IsConnectMode()
}

func (srv *ServerConfig) GetRepeaterSrv() string {
	return srv.RepeaterSrv
}

func (srv *ServerConfig) GetRepeaterPort() int {
	return srv.RepeaterPort
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

package config

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	ServerMode   string `json:"ServerMode"`  // connect 还是 proxy
	RepeaterSrv  string `json:"RepeaterSrv"` // repeaterSrv
	RepeaterPort int    `json:"RepeaterPort"`
}

func InitNet(config *ServerConfig) {
	srvCfg := selectConfigFile()
	file, err := ioutil.ReadFile(srvCfg)
	if err != nil {
		log.Fatalf("f6918311 read %s fail, err(%s)\n", srvCfg, err)
	}
	err = json.Unmarshal(file[:], config)
	if err != nil {
		log.Fatalf("f6918311 Unmarshal [%v] fail, err(%s)\n", string(file[:]), err)
	}
}

func InitConfig(domainCfgPath string) map[string]string {
	if domainCfgPath == "" {
		domainCfgPath = "./config/localDomain.conf"
	}

	rand.Seed(time.Now().UnixNano())

	domainDnsMap := make(map[string]string, 65536)

	fi, err := os.Open(domainCfgPath)
	if err != nil {
		log.Fatalf("FATAL bc360dee open config file(%s) err(%s)", domainCfgPath, err.Error())
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		line := string(a)
		//server=/cn/223.5.5.5
		//server=/acm.org/223.5.5.5
		if strings.HasPrefix(line, "server=/") == false {
			log.Fatalf("c772598d WARN invalid line:%s\n", line)
		}
		// 去掉 server=/
		idx := len("server=/")
		line = line[idx:]

		mapInfo := strings.Split(line, "/")
		if len(mapInfo) == 2 {
			name := mapInfo[0]
			dns := mapInfo[1]
			domainDnsMap[name] = dns
		} else {
			log.Fatalf("FATAL d807d85f valid domain cfg:%s\n", mapInfo)
		}
	}

	log.Printf("INFO a4347775 read %d domain info\n", len(domainDnsMap))

	return domainDnsMap
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
	return srv.ServerMode == "connect"
}

func (srv *ServerConfig) IsProxy() bool {
	return srv.ServerMode == "proxy"
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

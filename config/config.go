package config

import (
	"bufio"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	domainDnsMgr *domainDnsMgrT
)

type domainDnsMgrT struct {
	domainDnsMap   map[string]string
	pubDNSServerIP []string
	mtx            sync.RWMutex
}

func InitConfig(path string) {
	if path == "" {
		path = "./config/localDomain.conf"
	}

	rand.Seed(time.Now().UnixNano())

	domainDnsMgr = new(domainDnsMgrT)
	domainDnsMgr.domainDnsMap = make(map[string]string, 65536)
	domainDnsMgr.pubDNSServerIP = []string{"8.8.8.8", "8.8.4.4", "1.1.1.1", "199.85.126.10",
		"199.85.127.10", "208.67.222.222", "208.67.220.220", "84.200.69.80",
		"84.200.70.40", "8.26.56.26", "8.20.247.20", "64.6.64.6",
		"64.6.65.6", "192.95.54.3", "192.95.54.1", " 81.218.119.11",
		"209.88.198.133"}

	fi, err := os.Open(path)
	if err != nil {
		log.Fatalf("FATAL bc360dee open config file(%s) err(%s)", path, err.Error())
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
			domainDnsMgr.domainDnsMap[name] = dns
		} else {
			log.Fatalf("FATAL d807d85f valid domain cfg:%s\n", mapInfo)
		}
	}

	log.Printf("INFO a4347775 read %d domain info\n", len(domainDnsMgr.domainDnsMap))
}

func fastFind(domain string) (string, bool) {
	domainDnsMgr.mtx.Lock()
	defer domainDnsMgr.mtx.Unlock()

	ip, ok := domainDnsMgr.domainDnsMap[domain]
	return ip, ok
}

func slowFind(domain string) string {
	domainDnsMgr.mtx.Lock()
	defer domainDnsMgr.mtx.Unlock()

	// www.video.baidu.com
	mapInfo := strings.Split(domain, ".")
	for i := 1; i < len(mapInfo); i++ {
		subDomain := strings.Join(mapInfo[i:], ".")
		ip, ok := domainDnsMgr.domainDnsMap[subDomain]
		if ok {
			return ip
		}
	}

	idx := (int)(rand.Uint32()) % len(domainDnsMgr.pubDNSServerIP)
	dns := domainDnsMgr.pubDNSServerIP[idx]
	domainDnsMgr.domainDnsMap[domain] = dns

	return dns
}

func GetDNSServerIP(domain string) string {
	dns, ok := fastFind(domain)
	if ok {
		return dns
	}

	return slowFind(domain)
}

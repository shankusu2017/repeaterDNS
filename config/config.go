package config

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func InitConfig(path string) map[string]string {
	if path == "" {
		path = "./config/localDomain.conf"
	}

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
			log.Printf("c772598d WARN invalid line:%s\n", line)
			continue
		}
		// 去掉 server=/
		idx := len("server=/")
		line = line[idx:]

		addrInfo := strings.Split(line, "/")
		if len(addrInfo) == 2 {
			ip := addrInfo[0]
			netmask := addrInfo[1]
			subMask, err := strconv.Atoi(netmask)
			if err != nil {
				log.Fatalf("FATAL 62b5190e netmask invalid(%s)\n", netmask)
			} else {
				if subMask >= 0 && subMask < len(p.localIPSubnet) {
					p.localIPSubnet[subMask][ip] = true
					//fmt.Printf("%s %s\n", ip, netmask)
				} else {
					fmt.Printf("WARN 6b4d85ca netmask invalid(%s)\n", netmask)
				}
			}
		} else {
			log.Fatal(fmt.Sprintf("valid gateWay addr:%s", addrInfo))
		}
	}
}

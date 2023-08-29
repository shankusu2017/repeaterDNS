package main

import (
	"github.com/shankusu2017/repeaterDNS/listener"
	"github.com/shankusu2017/repeaterDNS/resolver"
	"log"
	"time"
)

var records map[string]string

func main() {
	records = map[string]string{
		"baidu.com":  "223.143.166.121",
		"github.com": "79.52.123.201",
	}

	resolver.Init()

	listener.Init()
	log.Printf("INFO f21ab893 init done\n")
	listener.Start(resolver.Resolvev2)

	time.Sleep(time.Hour * 65536)
}

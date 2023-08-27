package main

import (
	"github.com/shankusu2017/repeaterDNS/listener"
	"github.com/shankusu2017/repeaterDNS/resolver"
	"time"
)

var records map[string]string

func main() {
	records = map[string]string{
		"baidu.com":  "223.143.166.121",
		"github.com": "79.52.123.201",
	}

	resolver.Init()

	u := listener.Init()
	listener.Start(u, resolver.Resolve)

	time.Sleep(time.Hour * 65536)
}

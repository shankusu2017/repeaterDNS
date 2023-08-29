package main

import (
	"github.com/shankusu2017/repeaterDNS/config"
	"github.com/shankusu2017/repeaterDNS/listener"
	"github.com/shankusu2017/repeaterDNS/lookup"
	"time"
)

var (
	cfg config.ServerConfig
)

func main() {
	config.Init(&cfg)

	lookup.Init()

	listener.Init(&cfg)
	listener.StartLoop(lookup.Resolve)

	time.Sleep(time.Hour * 65536)
}

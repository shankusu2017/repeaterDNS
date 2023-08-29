package main

import (
	"github.com/shankusu2017/repeaterDNS/config"
	"github.com/shankusu2017/repeaterDNS/listener"
	"github.com/shankusu2017/repeaterDNS/lookup"
	"github.com/shankusu2017/repeaterDNS/repeater"
	"time"
)

var (
	cfg config.ServerConfig
)

func main() {
	config.Init(&cfg)

	lookup.Init()

	if cfg.IsConnectMode() {
		listener.Init()
	} else {
		repeater.Init(&cfg)
	}

	if cfg.IsConnectMode() {
		listener.StartLoop(lookup.Resolve)
	} else {
		repeater.StartLoop()
	}

	time.Sleep(time.Hour * 65536)
}

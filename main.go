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
	config.InitNet(&cfg)

	lookup.Init()

	listener.Init()
	listener.Start(lookup.Resolve)

	time.Sleep(time.Hour * 65536)
}

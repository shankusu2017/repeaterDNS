package main

import (
	"github.com/shankusu2017/repeaterDNS/config"
	"github.com/shankusu2017/repeaterDNS/listener"
	"github.com/shankusu2017/repeaterDNS/resolver"
	"time"
)

var (
	cfg config.ServerConfig
)

func main() {
	config.InitNet(&cfg)

	resolver.Init()

	listener.Init()
	listener.Start(resolver.Resolve)

	time.Sleep(time.Hour * 65536)
}

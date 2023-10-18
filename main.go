package main

import (
	"github.com/shankusu2017/repeaterDNS/config"
	"github.com/shankusu2017/repeaterDNS/listener"
	"github.com/shankusu2017/repeaterDNS/lookup"
	"github.com/shankusu2017/repeaterDNS/repeater"
	"log"
	"net/http"
	_ "net/http/pprof"
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
		listener.StartLoopResolve(lookup.Resolve)
	} else {
		repeater.StartLoop()
	}

	lookup.StartLoopDeadlineCheck()

	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	// 定期重启，用于暂时解决 OOM 问题
	time.Sleep(time.Hour * 1)
}

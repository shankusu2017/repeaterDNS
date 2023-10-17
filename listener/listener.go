package listener

import (
	"github.com/shankusu2017/constant"
	"github.com/shankusu2017/repeaterDNS/lookup"
	"log"
	"net"
	"sync"
	"time"
)

var (
	listener *net.UDPConn
	instance sync.Once
)

func Init() {
	// Listen on UDP Port
	addr := net.UDPAddr{
		Port: 53,
		IP:   net.ParseIP("0.0.0.0"),
	}

	var err error
	listener, err = net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatalf("f0f008d0 FATAL listen(%s) fail, error(%s)\n", addr.String(), err)
	}
}

func StartLoopResolve(f func(net.Addr, []byte)) {
	go loopRcv(f)
}

func loopRcv(f func(net.Addr, []byte)) {
	log.Printf("DEBUG 9d5ff164 listen start\n")
	for {
		buf := make([]byte, constant.Size256K)
		n, clientAddr, err := listener.ReadFrom(buf)
		if err != nil {
			log.Printf("84ded1d7 ERROR err:%s, add:%s\n", err.Error(), clientAddr.String())
			continue
		}
		//log.Printf("DEBUG 89aa4cad rcv request from:%s\n", clientAddr.String())
		buf = buf[:n]
		go f(clientAddr, buf)
	}
}

func Send(addr net.Addr, b []byte) {
	listener.WriteTo(b, addr)
}

func StartLoopDeadlineCheck() {
	ticker1 := time.NewTicker(5 * time.Second)
	// 一定要调用Stop()，回收资源
	defer ticker1.Stop()
	go func(t *time.Ticker) {
		for {
			// 每5秒中从chan t.C 中读取一次
			<-t.C
			lookup.DeadlineCheck()
		}
	}(ticker1)
}

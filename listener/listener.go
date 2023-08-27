package listener

import (
	"log"
	"net"
	"sync"
)

var (
	listener *net.UDPConn
	instance sync.Once
)

func Init() {
	instance.Do(initDo)
}

func initDo() {
	// Listen on UDP Port
	addr := net.UDPAddr{
		Port: 1053,
		IP:   net.ParseIP("0.0.0.0"),
	}

	var err error
	listener, err = net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatalf("f0f008d0 FATAL listen(%s) fail, error(%s)\n", addr.String(), err)
	}
}

func Start(f func(net.Addr, []byte)) {
	// Wait to get request on that port
	for {
		buf := make([]byte, 4096)
		n, clientAddr, err := listener.ReadFrom(buf)
		if err != nil {
			log.Printf("84ded1d7 ERROR err:%s, add:%s\n", err.Error(), clientAddr.String())
			continue
		}
		buf = buf[:n]
		f(clientAddr, buf)
	}
}

func Send(addr net.Addr, b []byte) {
	listener.WriteTo(b, addr)
}

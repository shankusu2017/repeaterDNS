package listener

import (
	"log"
	"net"
)

func Init() *net.UDPConn {
	// Listen on UDP Port
	addr := net.UDPAddr{
		Port: 53,
		IP:   net.ParseIP("0.0.0.0"),
	}

	u, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatalf("f0f008d0 FATAL listen(%s) fail, error(%s)\n", addr.String(), err)
	}

	return u
}

func Start(u *net.UDPConn, f func(net.Addr, []byte)) {
	// Wait to get request on that port
	for {
		buf := make([]byte, 4096)
		n, clientAddr, err := u.ReadFrom(buf)
		if err != nil {
			log.Printf("84ded1d7 ERROR err:%s, add:%s\n", err.Error(), clientAddr.String())
			continue
		}
		buf = buf[:n]
		f(clientAddr, buf)
	}
}

func Send(addr net.Addr, b []byte) {

}

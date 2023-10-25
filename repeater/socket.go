package repeater

import (
	"github.com/shankusu2017/constant"
	"log"
	"net"
)

func SendAndRcv(dns string, port int, b []byte) []byte {
	buf := make([]byte, constant.Size256K)

	addr := &net.UDPAddr{IP: net.ParseIP(dns), Port: port}
	udp, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Printf("ERROR deef8c7d error:%s\n", err.Error())
		return []byte{}
	}
	defer udp.Close()

	n, err := udp.Write(b)
	if err != nil {
		log.Printf("ERROR 69224d63 error:%s\n", err.Error())
		return []byte{}
	}

	n, err = udp.Read(buf[:])
	if err != nil {
		log.Printf("ERROR d098e4a5 error:%s, dns:%s\n", err.Error(), dns)
		return []byte{}
	}
	if n <= 0 {
		log.Printf("ERROR d0984d63 error:%d\n", n)
		return []byte{}
	}

	return buf[:n]
}

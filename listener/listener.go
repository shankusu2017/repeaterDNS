package listener

import (
	"github.com/shankusu2017/constant"
	"log"
	"net"
)

var (
	listener        *net.UDPConn
	listener1270053 *net.UDPConn
)

// NOTE: 无法同时绑定 0.0.0.0 和 127.0.0.53（本机使用）
// 修改配置文件，将 127.0.0.53 改为 127.0.0.1 ( https://unix.stackexchange.com/questions/612416/why-does-etc-resolv-conf-point-at-127-0-0-53 )
func Init() {
	// Listen on UDP Port
	addr := net.UDPAddr{
		Port: 53,
		// 一个服务器单独做 DNS 服务时，直接填写 0.0.0.0
		// huawei的一台服务器，两张网卡，两个公网，
		// 填写 作为 dns 服务器ip的对应的网卡的 local.ip.addr，直接填 0.0.0.0 ，但第二块网卡才是 DNS 服务器时
		// 会导致回写的 destion.netcard 为 第一块网卡，导致 dns 解析结果无法发回客户端
		IP: net.ParseIP("192.168.1.132"),
	}
	var err error
	listener, err = net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatalf("f0f008d0 FATAL listen(%s) fail, error(%s)\n", addr.String(), err)
	} else {
		log.Printf("0x731e41a8 listen(%s) OK\n", addr.String())
	}

	// for ping etc
	//addr = net.UDPAddr{
	//	Port: 53,
	//	IP:   net.ParseIP("127.0.0.1"),
	//}
	//listener1270053, err = net.ListenUDP("udp", &addr)
	//if err != nil {
	//	log.Fatalf("0x3a3ccce0 FATAL listen(%s) fail, error(%s)\n", addr.String(), err)
	//} else {
	//	log.Printf("0x4fc134c1 listen(%s) OK\n", addr.String())
	//}
}

func StartLoopResolve(f func(*net.UDPConn, net.Addr, []byte)) {
	if listener != nil {
		go loopRcv(listener, f)
	}
	if listener1270053 != nil {
		go loopRcv(listener1270053, f)
	}
}

func loopRcv(conn *net.UDPConn, f func(*net.UDPConn, net.Addr, []byte)) {
	log.Printf("DEBUG 9d5ff164 listen start\n")
	for {
		buf := make([]byte, constant.Size256K)
		n, clientAddr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Printf("84ded1d7 ERROR err:%s, add:%s\n", err.Error(), clientAddr.String())
			continue
		}
		//log.Printf("DEBUG 89aa4cad rcv request from:%s\n", clientAddr.String())
		buf = buf[:n]
		go f(conn, clientAddr, buf)
	}
}

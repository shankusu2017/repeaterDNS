package repeater

import (
	"github.com/shankusu2017/constant"
	"github.com/shankusu2017/repeaterDNS/config"
	"github.com/shankusu2017/utils"
	"log"
	"net"
	"strconv"
)

var (
	srvConn net.PacketConn
)

func Init(cfg *config.ServerConfig) {
	l, err := net.ListenPacket("udp", ":"+strconv.Itoa(cfg.GetRepeaterPort()))
	if err != nil {
		log.Fatalf("FATAL c378234f, listen:%d err:%v\n", cfg.GetRepeaterPort(), err)
	}

	srvConn = l
}

func StartLoop() {
	go loopRcv()
}

func loopRcv() {
	log.Printf("DEBUG 95514e81 repeaterSRV listen start\n")
	for {
		buf := make([]byte, constant.Size256K)
		n, clientAddr, err := srvConn.ReadFrom(buf)
		if err != nil {
			log.Printf("84ded1d7 ERROR err:%s, add:%s\n", err.Error(), clientAddr.String())
			continue
		}
		log.Printf("DEBUG 89aa4cad rcv request from:%s\n", clientAddr.String())
		buf = buf[:n]
		go send2PublicDNSAndRepater2Cli(clientAddr, buf)
	}
}

func send2PublicDNSAndRepater2Cli(clientAddr net.Addr, b []byte) {
	// 解密来自客户端的消息
	plainText, err := utils.AESDeCrypt(b, config.GetIV16(), config.GetKey16())
	if err != nil {
		log.Printf("ERROR 342adfe9 invalid.b.len:%d\n", len(b))
		return
	}

	// 转发给公共服务器
	rsp := SendReq2LocalAndRcvRsp(config.GetRandomPublicDNS(), plainText)

	// 先加密，再讲结果发给repeaterCli
	text := utils.AESCrypt(rsp, config.GetIV16(), config.GetKey16())
	srvConn.WriteTo(text, clientAddr)
}

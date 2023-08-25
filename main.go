package main

import (
	"fmt"
	"github.com/google/gopacket"
	layers "github.com/google/gopacket/layers"
	"log"
	"net"
)

var records map[string]string

func main() {
	records = map[string]string{
		"baidu.com":  "223.143.166.121",
		"github.com": "79.52.123.201",
	}

	//Listen on UDP Port
	addr := net.UDPAddr{
		Port: 53,
		IP:   net.ParseIP("0.0.0.0"),
	}
	u, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatalf("f0f008d0 FATAL error(%s)\n", err)
	}

	// Wait to get request on that port
	for {
		buf := make([]byte, 4096)
		n, clientAddr, err := u.ReadFrom(buf)
		if err != nil {
			log.Printf("84ded1d7 ERROR err:%s, add:%s\n", err.Error(), clientAddr.String())
			continue
		}
		buf = buf[:n]

		packet := gopacket.NewPacket(buf, layers.LayerTypeDNS, gopacket.Default)
		dnsPacket := packet.Layer(layers.LayerTypeDNS)
		tcp, _ := dnsPacket.(*layers.DNS)
		go serveDNS(u, clientAddr, tcp)
	}
}

func serveDNS(u *net.UDPConn, clientAddr net.Addr, request *layers.DNS) {
	replyMess := request
	var dnsAnswer layers.DNSResourceRecord
	dnsAnswer.Type = layers.DNSTypeA
	var ip string
	var err error
	var ok bool
	// TODO 根据域名来决定走墙内解析还是墙外解析
	ip, ok = records[string(request.Questions[0].Name)]
	if !ok {
		//Todo: Log no data present for the IP and handle:todo
	}
	a, _, _ := net.ParseCIDR(ip + "/24")
	dnsAnswer.Type = layers.DNSTypeA
	dnsAnswer.IP = a
	dnsAnswer.Name = []byte(request.Questions[0].Name)
	fmt.Println(request.Questions[0].Name)
	dnsAnswer.Class = layers.DNSClassIN
	replyMess.QR = true
	replyMess.ANCount = 1
	replyMess.OpCode = layers.DNSOpCodeNotify
	replyMess.AA = true
	replyMess.Answers = append(replyMess.Answers, dnsAnswer)
	replyMess.ResponseCode = layers.DNSResponseCodeNoErr
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{} // See SerializeOptions for more details.
	err = replyMess.SerializeTo(buf, opts)
	if err != nil {
		panic(err)
	}
	u.WriteTo(buf.Bytes(), clientAddr)
}

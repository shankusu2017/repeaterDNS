package resolver

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/shankusu2017/repeaterDNS/listener"
	"github.com/shankusu2017/repeaterDNS/proto"
	"github.com/shankusu2017/repeaterDNS/resolver/cache"
	"github.com/shankusu2017/repeaterDNS/resolver/lookup"
	"github.com/shankusu2017/utils"
	"log"
	"net"
	"sync"
)

var (
	instance sync.Once
)

func Init() {
	instance.Do(InitDo)
}

func InitDo() {
	cache.Init()
	lookup.Init()
}

func ResolveDone(clientAddr net.Addr, request *layers.DNS, ip string) {
	fmt.Println(request.Questions[0].Name)

	var dnsAnswer layers.DNSResourceRecord
	dnsAnswer.Type = layers.DNSTypeA
	a, _, _ := net.ParseCIDR(ip + "/24")
	dnsAnswer.Type = layers.DNSTypeA
	dnsAnswer.IP = a
	dnsAnswer.Name = []byte(request.Questions[0].Name)
	dnsAnswer.Class = layers.DNSClassIN

	replyMess := request
	replyMess.QR = true
	replyMess.ANCount = 1
	replyMess.OpCode = layers.DNSOpCodeQuery
	replyMess.AA = true
	replyMess.Answers = append(replyMess.Answers, dnsAnswer)
	replyMess.ResponseCode = layers.DNSResponseCodeNoErr
	buffer := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{} // See SerializeOptions for more details.
	err := replyMess.SerializeTo(buffer, opts)
	if err != nil {
		panic(err)
	}

	listener.Send(clientAddr, buffer.Bytes())
}

func Resolve(clientAddr net.Addr, buf []byte) {
	request := proto.Buf2DNSReq(buf)
	domain := request.Questions[0].Name
	recode := cache.GetIP(string(domain))
	if recode != nil {
		ResolveDone(clientAddr, request, recode.GetRandIP())
		log.Printf("INFO c7a8a141 resolved by cache, domain:%s, ip:%v\n", domain, recode.GetAllIP())
		return
	} else {
		//
		ips := lookup.LookUP(string(domain))
		if len(ips) > 0 {
			cache.SetIP(string(domain), ips)
		}
		ResolveDone(clientAddr, request, utils.SliceRandOne(ips))
		log.Printf("INFO ff7b7bcc by upServer, domain: %s, ip: %v\n", domain, ips)
		return
	}

	//
	//
	//replyMess := request
	//var dnsAnswer layers.DNSResourceRecord
	//dnsAnswer.Type = layers.DNSTypeA
	//var ip string
	//var err error
	//
	//ip = config.GetDNSServerIP(string(request.Questions[0].Name))
	//
	//a, _, _ := net.ParseCIDR(ip + "/24")
	//dnsAnswer.Type = layers.DNSTypeA
	//dnsAnswer.IP = a
	//dnsAnswer.Name = []byte(request.Questions[0].Name)
	//fmt.Println(request.Questions[0].Name)
	//dnsAnswer.Class = layers.DNSClassIN
	//replyMess.QR = true
	//replyMess.ANCount = 1
	//replyMess.OpCode = layers.DNSOpCodeNotify
	//replyMess.AA = true
	//replyMess.Answers = append(replyMess.Answers, dnsAnswer)
	//replyMess.ResponseCode = layers.DNSResponseCodeNoErr
	//buffer := gopacket.NewSerializeBuffer()
	//opts := gopacket.SerializeOptions{} // See SerializeOptions for more details.
	//err = replyMess.SerializeTo(buffer, opts)
	//if err != nil {
	//	panic(err)
	//}
	//
	//listener.Send(clientAddr, buffer.Bytes())
}

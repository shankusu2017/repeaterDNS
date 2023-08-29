package resolver

import (
	"github.com/shankusu2017/repeaterDNS/listener"
	"github.com/shankusu2017/repeaterDNS/proto"
	"github.com/shankusu2017/repeaterDNS/resolver/cache"
	"github.com/shankusu2017/repeaterDNS/resolver/lookup"
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

func ResolveDone(clientAddr net.Addr, rsp []byte) {
	listener.Send(clientAddr, rsp)
}

func Resolve(clientAddr net.Addr, b []byte) {
	request := proto.Buf2DNSReq(b)
	domain := string(request.Questions[0].Name)
	recode := cache.GetRecord(domain)
	if recode != nil {
		ResolveDone(clientAddr, recode.GetRsp())
		log.Printf("INFO c7a8a141 resolved by cache, domain:%s, rsp:%v\n", domain, recode.GetRsp())
		return
	} else {
		rsp := lookup.Lookup(b, domain)
		if len(rsp) > 0 {
			cache.SetRecord(domain, rsp)
		}
		ResolveDone(clientAddr, rsp)
		log.Printf("INFO ff7b7bcc by upServer, domain: %s, rsp: %v\n", domain, rsp)
		return
	}
}

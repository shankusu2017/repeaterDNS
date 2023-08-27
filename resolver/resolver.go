package resolver

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/shankusu2017/repeaterDNS/config"
	"github.com/shankusu2017/repeaterDNS/proto"
	"github.com/shankusu2017/repeaterDNS/resolver/cache"
	"math/rand"
	"net"
	"strings"
	"sync"
)

type resolverMgrT struct {
	domainDnsMap   map[string]string         // 域名 ----> 负责解析此域名的 dns 服务器地址
	pubDNSServerIP []string                  // 全球公开的 dns 地址
	domainIPMap    map[string]*cache.RecordT // 域名 ----> 解析结果
	mtx            sync.RWMutex
}

var (
	instance    sync.Once
	resolverMgr *resolverMgrT
)

func Init() {
	instance.Do(InitDo)
}

func InitDo() {
	resolverMgr = new(resolverMgrT)
	resolverMgr.domainDnsMap, resolverMgr.pubDNSServerIP = config.InitConfig("./config/localDomain.conf")
	resolverMgr.domainIPMap = make(map[string]*cache.RecordT, constant.Size1K)
}

func Resolve(clientAddr net.Addr, buf []byte) {
	request := proto.Buf2DNSReq(buf)

	replyMess := request
	var dnsAnswer layers.DNSResourceRecord
	dnsAnswer.Type = layers.DNSTypeA
	var ip string
	var err error

	ip = config.GetDNSServerIP(string(request.Questions[0].Name))

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

func fastFind(domain string) (string, bool) {
	domainDnsMgr.mtx.Lock()
	defer domainDnsMgr.mtx.Unlock()

	ip, ok := domainDnsMgr.domainDnsMap[domain]
	return ip, ok
}

func slowFind(domain string) string {
	domainDnsMgr.mtx.Lock()
	defer domainDnsMgr.mtx.Unlock()

	// www.video.baidu.com
	mapInfo := strings.Split(domain, ".")
	for i := 1; i < len(mapInfo); i++ {
		subDomain := strings.Join(mapInfo[i:], ".")
		ip, ok := domainDnsMgr.domainDnsMap[subDomain]
		if ok {
			return ip
		}
	}

	idx := (int)(rand.Uint32()) % len(domainDnsMgr.pubDNSServerIP)
	dns := domainDnsMgr.pubDNSServerIP[idx]
	domainDnsMgr.domainDnsMap[domain] = dns

	return dns
}

func GetDNSServerIP(domain string) string {
	dns, ok := fastFind(domain)
	if ok {
		return dns
	}

	return slowFind(domain)
}

package lookup

import (
	"context"
	"fmt"
	"github.com/shankusu2017/constant"
	"github.com/shankusu2017/repeaterDNS/config"
	"github.com/shankusu2017/utils"
	"log"
	"net"
	"strings"
	"sync"
)

type lookupMgrT struct {
	domainDnsMap   map[string]string // 域名 ----> 负责解析此域名的 dns 服务器地址
	pubDNSServerIP []string          // 全球公开的 dns 地址
	mtx            sync.RWMutex
}

var (
	lookupMgr *lookupMgrT
)

func Init() {
	lookupMgr = new(lookupMgrT)
	lookupMgr.domainDnsMap, lookupMgr.pubDNSServerIP = config.InitConfig("./config/localDomain.conf")
}

func LookUP(domain string) []string {
	dns := findDNS(domain)
	ip := lookupHost(domain, dns)
	return ip
}

func findDNS(domain string) string {
	dns, ok := fastFind(domain)
	if ok {
		return dns
	}

	return slowFind(domain)
}

func fastFind(domain string) (string, bool) {
	lookupMgr.mtx.Lock()
	defer lookupMgr.mtx.Unlock()

	ip, ok := lookupMgr.domainDnsMap[domain]
	return ip, ok
}

func slowFind(domain string) string {
	lookupMgr.mtx.Lock()
	defer lookupMgr.mtx.Unlock()

	// www.video.baidu.com
	mapInfo := strings.Split(domain, ".")
	for i := 1; i < len(mapInfo); i++ {
		subDomain := strings.Join(mapInfo[i:], ".")
		ip, ok := lookupMgr.domainDnsMap[subDomain]
		if ok {
			return ip
		}
	}

	dns := utils.SliceRandOne(lookupMgr.pubDNSServerIP)

	lookupMgr.domainDnsMap[domain] = dns

	return dns
}

// domain == www.google.com
func lookupHost(domain, dns string) []string {
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "udp", fmt.Sprintf("%s:53", dns))
		},
	}
	ips, err := resolver.LookupHost(context.Background(), domain)
	if err != nil {
		log.Printf("ERROR f69d570b lookupHost fail(%s)\n", err.Error())
		return []string{}
	}

	fmt.Printf("INFO dceb27f4 domain:%s, dns:%v\n", domain, ips)
	return ips
}

func Lookupv2(req []byte, dns string) []byte {
	addr := &net.UDPAddr{IP: net.ParseIP(dns), Port: 53}
	udp, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Printf("ERROR deef8c7d error:%s\n", err.Error())
		return []byte{}
	}

	n, err := udp.Write(req)
	if err != nil {
		log.Printf("ERROR 69224d63 error:%s\n", err.Error())
		return []byte{}
	}

	buf := make([]byte, constant.Size256K)
	n, err = udp.Read(buf[:])
	if err != nil {
		log.Printf("ERROR d098e4a5 error:%s\n", err.Error())
		return []byte{}
	}

	buf = buf[:n]
	return buf
}

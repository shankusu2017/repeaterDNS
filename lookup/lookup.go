package lookup

import (
	"bufio"
	"github.com/shankusu2017/constant"
	"github.com/shankusu2017/repeaterDNS/config"
	"github.com/shankusu2017/repeaterDNS/listener"
	"github.com/shankusu2017/repeaterDNS/proto"
	"github.com/shankusu2017/repeaterDNS/repeater"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type RecordT struct {
	recode []byte    // 域名对应的 dns 信息
	t      time.Time // 记录生成的时间(记录均有有效期)
}

func (r *RecordT) GetRsp() []byte {
	return r.recode
}

func (r *RecordT) IsExpired() bool {
	return r.t.After(time.Now().Add(constant.Time30Min))
}

type domainT struct {
	dns     string
	isLocal bool
}

func (r *domainT) GetDns() string {
	return r.dns
}

func (r *domainT) isLocalDomain() bool {
	return r.isLocal
}

type lookupMgrT struct {
	domain2DnsMap    map[string]*domainT // 域名 ----> 负责解析此域名的 dns 服务器地址
	domain2RecodeMap map[string]*RecordT
	mtx              sync.RWMutex
}

var (
	lookupMgr *lookupMgrT
)

func Init() {
	lookupMgr = new(lookupMgrT)
	lookupMgr.domain2DnsMap = initDomain2Dns("./config/localDomain.conf")
	lookupMgr.domain2RecodeMap = make(map[string]*RecordT, constant.Size1M)
}

func initDomain2Dns(domainCfgPath string) map[string]*domainT {
	rand.Seed(time.Now().UnixNano())
	if domainCfgPath == "" {
		domainCfgPath = "./config/localDomain.conf"
	}

	domainDnsMap := make(map[string]*domainT, 65536)

	fi, err := os.Open(domainCfgPath)
	if err != nil {
		log.Fatalf("FATAL bc360dee open config file(%s) err(%s)", domainCfgPath, err.Error())
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		line := string(a)
		//server=/cn/223.5.5.5
		//server=/acm.org/223.5.5.5
		if strings.HasPrefix(line, "server=/") == false {
			log.Fatalf("c772598d WARN invalid line:%s\n", line)
		}
		// 去掉 server=/
		idx := len("server=/")
		line = line[idx:]

		mapInfo := strings.Split(line, "/")
		if len(mapInfo) == 2 {
			name := mapInfo[0]
			dns := mapInfo[1]
			domainDnsMap[name] = &domainT{dns: dns, isLocal: true}
		} else {
			log.Fatalf("FATAL d807d85f valid domain cfg:%s\n", mapInfo)
		}
	}

	log.Printf("INFO a4347775 read %d domain info\n", len(domainDnsMap))

	return domainDnsMap
}

func findDNS(domain string) (string, bool) {
	dns, ok := fastFindDNS(domain)
	if ok {
		return dns.GetDns(), dns.isLocalDomain()
	}

	dns = slowFindDNS(domain)
	return dns.GetDns(), dns.isLocalDomain()
}

func fastFindDNS(domain string) (*domainT, bool) {
	lookupMgr.mtx.Lock()
	defer lookupMgr.mtx.Unlock()

	dns, ok := lookupMgr.domain2DnsMap[domain]
	return dns, ok
}

func slowFindDNS(domain string) *domainT {
	lookupMgr.mtx.Lock()
	defer lookupMgr.mtx.Unlock()

	// www.video.baidu.com
	mapInfo := strings.Split(domain, ".")
	for i := 1; i < len(mapInfo); i++ {
		subDomain := strings.Join(mapInfo[i:], ".")
		dns, ok := lookupMgr.domain2DnsMap[subDomain]
		if ok {
			return dns
		}
	}

	dns, _ := config.GetRepeaterSrvAddr()
	lookupMgr.domain2DnsMap[domain] = &domainT{dns: dns, isLocal: false}

	return lookupMgr.domain2DnsMap[domain]
}

func findRecord(domain string) *RecordT {
	lookupMgr.mtx.RLock()
	defer lookupMgr.mtx.RUnlock()

	record := lookupMgr.domain2RecodeMap[domain]
	if record == nil {
		return nil
	}

	if record.IsExpired() {
		delete(lookupMgr.domain2RecodeMap, domain)
		return nil
	}

	return record
}

func setRecord(domain string, rsp []byte) {
	if len(rsp) <= 1 {
		log.Printf("WARN 36430aed domain:%s, invalid.len:%d\n", domain, rsp)
		return
	}
	if config.IsCache() == false {
		return
	}

	record := new(RecordT)
	record.t = time.Now()
	record.recode = make([]byte, len(rsp))
	copy(record.recode[:], rsp)

	lookupMgr.mtx.Lock()
	defer lookupMgr.mtx.Unlock()
	lookupMgr.domain2RecodeMap[domain] = record
}

func lookHost(req []byte, domain string) []byte {
	record := findRecord(domain)
	if record != nil {
		return record.GetRsp()
	}

	var buf []byte
	dns, isLocalDns := findDNS(domain)
	if isLocalDns {
		buf = repeater.SendReq2LocalAndRcvRsp(dns, req)
	} else {
		buf = repeater.SendReq2OutsideAndRcvRsp(req)
	}

	setRecord(domain, buf)

	return buf
}

func Resolve(clientAddr net.Addr, b []byte) {
	request := proto.Buf2DNSReq(b)
	domain := string(request.Questions[0].Name)
	rsp := lookHost(b, domain)
	if len(rsp) > 0 {
		listener.Send(clientAddr, rsp)
	}
	log.Printf("INFO c7a8a141 resolved by cache, domain:%s, rsp:%s\n", domain, string(rsp))
}

package lookup

import (
	"bufio"
	"github.com/shankusu2017/constant"
	"github.com/shankusu2017/repeaterDNS/config"
	"github.com/shankusu2017/repeaterDNS/listener"
	"github.com/shankusu2017/repeaterDNS/proto"
	"github.com/shankusu2017/utils"
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

type lookupMgrT struct {
	domain2DnsMap    map[string]string // 域名 ----> 负责解析此域名的 dns 服务器地址
	pubDNSServerIP   []string          // 全球公开的 dns 地址
	domain2RecodeMap map[string]*RecordT
	mtx              sync.RWMutex
}

var (
	lookupMgr *lookupMgrT
	instance  sync.Once
)

func Init() {
	instance.Do(initDo)
}

func initDo() {
	lookupMgr = new(lookupMgrT)
	lookupMgr.domain2DnsMap = initDomain2Dns("./config/localDomain.conf")
	lookupMgr.pubDNSServerIP = config.GetPublicDNS()
	lookupMgr.domain2RecodeMap = make(map[string]*RecordT, constant.Size256K)
}

func resolveDone(clientAddr net.Addr, rsp []byte) {
	listener.Send(clientAddr, rsp)
}

func Resolve(clientAddr net.Addr, b []byte) {
	request := proto.Buf2DNSReq(b)
	domain := string(request.Questions[0].Name)
	rsp := lookHost(b, domain)
	if len(rsp) > 0 {
		resolveDone(clientAddr, rsp)
	}
	log.Printf("INFO c7a8a141 resolved by cache, domain:%s, rsp:%s\n", domain, string(rsp))
}

func initDomain2Dns(domainCfgPath string) map[string]string {
	rand.Seed(time.Now().UnixNano())
	if domainCfgPath == "" {
		domainCfgPath = "./config/localDomain.conf"
	}

	domainDnsMap := make(map[string]string, 65536)

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
			domainDnsMap[name] = dns
		} else {
			log.Fatalf("FATAL d807d85f valid domain cfg:%s\n", mapInfo)
		}
	}

	log.Printf("INFO a4347775 read %d domain info\n", len(domainDnsMap))

	return domainDnsMap
}

func findDNS(domain string) string {
	dns, ok := fastFindDNS(domain)
	if ok {
		return dns
	}

	return slowFindDNS(domain)
}

func fastFindDNS(domain string) (string, bool) {
	lookupMgr.mtx.Lock()
	defer lookupMgr.mtx.Unlock()

	ip, ok := lookupMgr.domain2DnsMap[domain]
	return ip, ok
}

func slowFindDNS(domain string) string {
	lookupMgr.mtx.Lock()
	defer lookupMgr.mtx.Unlock()

	// www.video.baidu.com
	mapInfo := strings.Split(domain, ".")
	for i := 1; i < len(mapInfo); i++ {
		subDomain := strings.Join(mapInfo[i:], ".")
		ip, ok := lookupMgr.domain2DnsMap[subDomain]
		if ok {
			return ip
		}
	}

	dns := utils.SliceRandOne(lookupMgr.pubDNSServerIP)
	lookupMgr.domain2DnsMap[domain] = dns

	return dns
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

	dns := findDNS(domain)
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
	if n <= 0 {
		log.Printf("ERROR d0984d63 error:%d\n", n)
		return []byte{}
	}

	buf = buf[:n]

	setRecord(domain, buf)

	return buf
}

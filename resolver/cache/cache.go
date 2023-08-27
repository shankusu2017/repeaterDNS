package cache

import (
	"github.com/shankusu2017/constant"
	"math/rand"
	"sync"
	"time"
)

type RecordT struct {
	IP []string  // 域名对应的 IP
	T  time.Time // 记录生成的时间(记录均有有效期)
}

func (r *RecordT) GetRandIP() string {
	if len(r.IP) == 0 {
		return ""
	}

	ret := rand.Uint32()
	i := (int)(ret) % len(r.IP)
	return r.IP[i]
}

type cacheMgrT struct {
	domain2recodeMap map[string]*RecordT
	mtx              sync.RWMutex
}

var (
	cacheMgr *cacheMgrT
)

func Init() {
	cacheMgr = new(cacheMgrT)
	cacheMgr.domain2recodeMap = make(map[string]*RecordT, constant.Size256K)
}

func GetIP(domain string) *RecordT {
	cacheMgr.mtx.RLock()
	defer cacheMgr.mtx.RUnlock()

	recode := cacheMgr.domain2recodeMap[domain]
	return recode
}

func SetIP(domain string, ip []string) {
	record := new(RecordT)
	record.T = time.Now()
	record.IP = make([]string, len(ip))
	copy(record.IP[:], ip)

	cacheMgr.mtx.Lock()
	defer cacheMgr.mtx.Unlock()
	cacheMgr.domain2recodeMap[domain] = record
}

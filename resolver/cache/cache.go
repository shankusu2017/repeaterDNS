package cache

import (
	"sync"
	"time"
)

type RecordT struct {
	ip []string  // 域名对应的 ip
	t  time.Time // 记录生成的时间(记录均有有效期)
}

type cacheMgrT struct {
	domain2recodeMap map[string]*RecordT
	mtx              sync.Mutex
}

var (
	cacheMgr     *cacheMgrT
	instanceOnce sync.Once
)

func onlyOneTime() {
	cacheMgr = new(cacheMgrT)
	cacheMgr.domain2recodeMap = make(map[string]*RecordT, 65536)
}

func Init() {
	instanceOnce.Do(onlyOneTime)
}

func GetIP(domain string) ([]string, bool) {
	cacheMgr.mtx.Lock()
	defer cacheMgr.mtx.Unlock()

	recode, ok := cacheMgr.domain2recodeMap[domain]
	if ok != true {
		return []string{}, false
	}

}

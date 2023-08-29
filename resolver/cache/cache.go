package cache

import (
	"github.com/shankusu2017/constant"
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

func GetRecord(domain string) *RecordT {
	cacheMgr.mtx.RLock()
	defer cacheMgr.mtx.RUnlock()

	record := cacheMgr.domain2recodeMap[domain]
	return record
}

func SetRecord(domain string, rsp []byte) {
	record := new(RecordT)
	record.t = time.Now()
	record.recode = make([]byte, len(rsp))
	copy(record.recode[:], rsp)

	cacheMgr.mtx.Lock()
	defer cacheMgr.mtx.Unlock()
	cacheMgr.domain2recodeMap[domain] = record
}

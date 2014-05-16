package main

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"time"
)

type MemcacheRequestInfo struct {
	*BaseRequestInfo
	hostAndPort string
	mc          *memcache.Client
}

var persistentMemcacheConnect *memcache.Client

func NewMemcacheRequestInfo(baseRequest *BaseRequestInfo) *MemcacheRequestInfo {
	r := &MemcacheRequestInfo{BaseRequestInfo: baseRequest}
	r.hostAndPort = *storageAddress
	return r
}

func (this *MemcacheRequestInfo) initConnect() bool {
	//connection pull is hidden by memcache client implementation...
	// init connect time will be always zero
	if *usePersistentConnections {
		if persistentMemcacheConnect == nil {
			persistentMemcacheConnect = memcache.New(this.hostAndPort)
		}
		this.mc = persistentMemcacheConnect
	} else {
		this.mc = memcache.New(this.hostAndPort)
	}
	return true
}

func (this *MemcacheRequestInfo) addLock() bool {
	start := time.Now()
	err := this.mc.Add(&memcache.Item{Key: "lock_" + this.key, Value: []byte("")})

	if err == nil || err == memcache.ErrNotStored {
		end := time.Now()
		this.addLockDuration = end.Sub(start)
	} else {
		fmt.Printf("Add Lock Request failed: %#v\n", err)
	}
	return (err == nil || err == memcache.ErrNotStored)
}

func (this *MemcacheRequestInfo) getData() bool {
	start := time.Now()
	_, err := this.mc.Get(this.key)

	if err == nil || err == memcache.ErrCacheMiss {
		end := time.Now()
		this.getDuration = end.Sub(start)
	} else {
		fmt.Printf("Get request failed: %#v\n", err)
	}
	return (err == nil || err == memcache.ErrCacheMiss)
}

func (this *MemcacheRequestInfo) setData() bool {
	start := time.Now()
	err := this.mc.Set(&memcache.Item{Key: this.key, Value: []byte(SESSION_CONTENT)})

	if err == nil {
		end := time.Now()
		this.setDuration = end.Sub(start)
	} else {
		fmt.Printf("Set Request failed: %#v\n", err)
	}
	return (err == nil)
}
func (this *MemcacheRequestInfo) deleteLock() bool {
	start := time.Now()
	err := this.mc.Delete("lock_" + this.key)

	if err == nil {
		end := time.Now()
		this.deleteLockDuration = end.Sub(start)
	} else {
		fmt.Printf("Delete Request failed: %#v\n", err)
	}
	return (err == nil)
}

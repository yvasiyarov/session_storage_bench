package session_storage_bench

import (
	"math/rand"
	"strconv"
	"time"
)

type RequestInfo interface {
	MakeRequest()
	IsFailed() bool
	GetDuration() time.Duration
}

type BaseRequestInfo struct {
	key                  string
	duration             time.Duration
	initConnectDuration  time.Duration
	addLockDuration      time.Duration
	getDuration          time.Duration
	setDuration          time.Duration
	deleteLockDuration   time.Duration
	closeConnectDuration time.Duration
	isFailed             bool
}

func NewBaseRequestInfo() *BaseRequestInfo {
	r := &BaseRequestInfo{}
	r.key = strconv.FormatInt(rand.Int63(), 10)
	return r
}

func (this *BaseRequestInfo) addLock() bool {
	return false
}
func (this *BaseRequestInfo) deleteLock() bool {
	return false
}
func (this *BaseRequestInfo) getData() bool {
	return false
}
func (this *BaseRequestInfo) setData() bool {
	return false
}
func (this *BaseRequestInfo) initConnect() bool {
	return false
}
func (this *BaseRequestInfo) closeConnect() bool {
	return false
}
func (this *BaseRequestInfo) IsFailed() bool {
	return this.isFailed
}
func (this *BaseRequestInfo) GetDuration() time.Duration {
	return this.getDuration
}

func (this *BaseRequestInfo) MakeRequest() {
	this.isFailed = true

	this.initConnect()
	lockResult := this.addLock()
	if !lockResult {
		return
	}

	getResult := this.getData()
	if !getResult {
		return
	}

	time.Sleep(200 * time.Millisecond)

	setResult := this.setData()
	if !setResult {
		return
	}

	unlockResult := this.deleteLock()
	if !unlockResult {
		return
	}
	this.closeConnect()
	this.isFailed = false

	this.duration = this.initConnectDuration + this.addLockDuration + this.getDuration + this.setDuration + this.deleteLockDuration + this.closeConnectDuration
}

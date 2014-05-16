package main

import (
	"math/rand"
	"strconv"
	"time"
)

type RequestInfo interface {
	//	MakeRequest()
	IsFailed() bool
	SetIsFailed(bool)
	GetDuration() time.Duration
	addLock() bool
	deleteLock() bool
	getData() bool
	setData() bool
	initConnect() bool
	closeConnect() bool
}

type BaseRequestInfo struct {
	key                  string
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
	return true
}

func (this *BaseRequestInfo) deleteLock() bool {
	return true
}
func (this *BaseRequestInfo) getData() bool {
	return true
}
func (this *BaseRequestInfo) setData() bool {
	return true
}
func (this *BaseRequestInfo) initConnect() bool {
	return true
}
func (this *BaseRequestInfo) closeConnect() bool {
	return true
}
func (this *BaseRequestInfo) IsFailed() bool {
	return this.isFailed
}
func (this *BaseRequestInfo) SetIsFailed(isFailed bool) {
	this.isFailed = isFailed
}
func (this *BaseRequestInfo) GetDuration() time.Duration {
	return this.initConnectDuration + this.addLockDuration + this.getDuration + this.setDuration + this.deleteLockDuration + this.closeConnectDuration
}

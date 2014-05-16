package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"sync"
	"time"
)

type RedisRequestInfo struct {
	*BaseRequestInfo
	hostAndPort string
	redis       redis.Conn
}

var redisConnectionPull *sync.Pool

func NewRedisRequestInfo(baseRequest *BaseRequestInfo) *RedisRequestInfo {
	r := &RedisRequestInfo{BaseRequestInfo: baseRequest}
	r.hostAndPort = *storageAddress
	return r
}

func getNewRedisConnection() redis.Conn {
	rc, err := redis.Dial("tcp", *storageAddress)
	if err != nil {
		fmt.Printf("Creating Connection failed: %#v" + err.Error() + "\n")
		panic("Creating Connection failed: %#v" + err.Error() + "\n")
	}
	return rc
}

func (this *RedisRequestInfo) initConnect() bool {
	start := time.Now()

	if *usePersistentConnections {
		if redisConnectionPull == nil {
			redisConnectionPull = &sync.Pool{
				New: func() interface{} {
					return getNewRedisConnection()
				},
			}
		}
		redis, ok := redisConnectionPull.Get().(redis.Conn)
		if !ok {
			return false
		}
		this.redis = redis
	} else {
		this.redis = getNewRedisConnection()
	}

	end := time.Now()
	this.initConnectDuration = end.Sub(start)

	return true
}

func (this *RedisRequestInfo) closeConnect() bool {
	start := time.Now()

	if *usePersistentConnections {
		redisConnectionPull.Put(this.redis)
	} else {
		this.redis.Close()
	}

	end := time.Now()
	this.closeConnectDuration = end.Sub(start)

	return true
}

func (this *RedisRequestInfo) addLock() bool {
	start := time.Now()
	_, err := this.redis.Do("SET", "lock_"+this.key, []byte(""))

	if err == nil {
		end := time.Now()
		this.addLockDuration = end.Sub(start)
	} else {
		fmt.Printf("Add Lock Request failed: %#v\n", err.Error())
	}
	return err == nil
}

func (this *RedisRequestInfo) getData() bool {
	start := time.Now()
	_, err := this.redis.Do("GET", this.key)

	if err == nil {
		end := time.Now()
		this.getDuration = end.Sub(start)
	} else {
		fmt.Printf("Get request failed: %#v\n", err.Error())
	}
	return err == nil
}

func (this *RedisRequestInfo) setData() bool {
	start := time.Now()
	_, err := this.redis.Do("SET", this.key, []byte(SESSION_CONTENT))

	if err == nil {
		end := time.Now()
		this.setDuration = end.Sub(start)
	} else {
		fmt.Printf("Set Request failed: %#v\n", err.Error())
	}
	return (err == nil)
}

func (this *RedisRequestInfo) deleteLock() bool {
	start := time.Now()
	_, err := this.redis.Do("DEL", "lock_"+this.key)

	if err == nil {
		end := time.Now()
		this.deleteLockDuration = end.Sub(start)
	} else {
		fmt.Printf("Delete Request failed: %#v\n", err.Error())
	}
	return (err == nil)
}

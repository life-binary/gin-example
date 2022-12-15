package main

import (
	"fmt"
	"sync"
	"time"
)

import (
	"golang.org/x/time/rate"
)

type limitGroup struct {
	limitMap  *sync.Map
	limitLock *sync.RWMutex
}

var limitInstance *limitGroup

func InitLimit() error {
	limitInstance = &limitGroup{
		limitMap:  &sync.Map{},
		limitLock: &sync.RWMutex{},
	}
	return nil
}

/**
 * key: rate limit by key
 * bucketNum: 桶大小
 * qps: 每秒允许的qps
 */
func GetToken(key string, bucketNum int, interval time.Duration) bool {
	limit, err := limitInstance.getLimit(key, bucketNum, interval)
	if err != nil {
		fmt.Printf("GetToken, key:", key, ", error:", err)
		return true
	}
	return limit.Allow()
}

func (l *limitGroup) getLimit(key string, bucketNum int, interval time.Duration) (*rate.Limiter, error) {
	if data, ok := l.limitMap.Load(key); ok {
		return data.(*rate.Limiter), nil
	}

	l.limitLock.Lock()
	defer l.limitLock.Unlock()
	if data, ok := l.limitMap.Load(key); ok {
		return data.(*rate.Limiter), nil
	}

	limiter := rate.NewLimiter(rate.Every(interval), bucketNum)

	data, _ := l.limitMap.LoadOrStore(key, limiter)
	return data.(*rate.Limiter), nil
}

package ttlmap

/*
originally composed by OneOfOne on StackOverflow. Edited to allow custom ttl per key -> value pair

Aaron Pitman
mm/dd/yyyy
10/02/2018


*/

import (
	"sync"
	"time"
)

type item struct {
	value      string
	lastAccess int64
	ttl        int
}

type TTLMap struct {
	wrappedMap map[string]*item
	mutex      sync.Mutex
}

func New(size int, maxTTL int) (ttlMap *TTLMap) {
	ttlMap = &TTLMap{wrappedMap: make(map[string]*item, size)}
	go func() {
		for now := range time.Tick(time.Second) {
			ttlMap.mutex.Lock()
			for k, v := range ttlMap.wrappedMap {
				if now.Unix()-v.lastAccess > int64(v.ttl) {
					delete(ttlMap.wrappedMap, k)
				}
			}
			ttlMap.mutex.Unlock()
		}
	}()
	return
}

func (ttlMap *TTLMap) Len() int {
	return len(ttlMap.wrappedMap)
}

func (ttlMap *TTLMap) Put(k, v string, ttl int) {
	ttlMap.mutex.Lock()
	it, ok := ttlMap.wrappedMap[k]

	if ttl <= 0 {
		ttl = 10
	}
	if !ok {
		it = &item{value: v, ttl: ttl}
		ttlMap.wrappedMap[k] = it
	}
	it.lastAccess = time.Now().Unix()
	ttlMap.mutex.Unlock()
}

func (ttlMap *TTLMap) Get(k string) (v string) {
	ttlMap.mutex.Lock()
	if it, ok := ttlMap.wrappedMap[k]; ok {
		v = it.value
		it.lastAccess = time.Now().Unix()
	}
	ttlMap.mutex.Unlock()
	return

}

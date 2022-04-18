package ttlmap

import (
	"runtime"
	"sync"
	"time"
)

type TTLMapV2 struct {
	items map[string]item
	mux   *sync.RWMutex
	stop  chan bool
}

func (tm *TTLMapV2) Get(key string) (interface{}, bool) {
	tm.mux.RLock()
	defer tm.mux.RUnlock()

	if item, ok := tm.items[key]; !ok || time.Now().UnixNano() > item.expire {
		return nil, false
	} else {
		return item.value, true
	}
}

func (tm *TTLMapV2) Set(key string, val interface{}, ex time.Duration) {
	tm.mux.Lock()
	tm.items[key] = item{value: val, expire: time.Now().Add(ex).UnixNano()}
	tm.mux.Unlock()
}

func (tm *TTLMapV2) Len() int {
	tm.mux.RLock()
	defer tm.mux.RUnlock()
	return len(tm.items)
}

func (tm *TTLMapV2) clean(tick time.Duration) {
	for range time.Tick(tick) {
		tm.mux.Lock()
		for key, item := range tm.items {
			if time.Now().UnixNano() >= item.expire {
				delete(tm.items, key)
			}
		}
		tm.mux.Unlock()
	}
}

func NewTTLMapV2(cleanTick time.Duration) *MapWrapper {
	tm := &TTLMapV2{items: map[string]item{}, mux: new(sync.RWMutex), stop: make(chan bool)}
	go tm.clean(cleanTick)
	tmw := &MapWrapper{TTLMap: tm}
	runtime.SetFinalizer(tmw, func(_ interface{}) { tm.stop <- true })
	return tmw
}

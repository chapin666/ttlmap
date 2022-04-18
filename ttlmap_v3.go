package ttlmap

import (
	"runtime"
	"sync"
	"time"
)

type TTLMapV3 struct {
	items map[string]item
	mux   *sync.RWMutex
	stop  chan bool
	now   int64
}

func (tm *TTLMapV3) Get(key string) (interface{}, bool) {
	tm.mux.RLock()
	defer tm.mux.RUnlock()
	if item, ok := tm.items[key]; !ok || tm.now > item.expire {
		return nil, false
	} else {
		return item.value, true
	}
}

func (tm *TTLMapV3) Set(key string, val interface{}, ex time.Duration) {
	tm.mux.Lock()
	tm.items[key] = item{value: val, expire: tm.now + ex.Nanoseconds()}
	tm.mux.Unlock()
}

func (tm *TTLMapV3) Len() int {
	tm.mux.RLock()
	defer tm.mux.RUnlock()
	return len(tm.items)
}

func (tm *TTLMapV3) clean(tick time.Duration) {
	for range time.Tick(tick) {
		tm.mux.Lock()
		for key, item := range tm.items {
			if tm.now >= item.expire {
				delete(tm.items, key)
			}
		}
		tm.mux.Unlock()
	}
}

func (tm *TTLMapV3) updateNow(tick time.Duration) {
	for range time.Tick(tick) {
		tm.mux.Lock()
		tm.now = time.Now().UnixNano()
		tm.mux.Unlock()
	}
}

func NewTTLMapV3(cleanTick time.Duration) *MapWrapper {
	tm := &TTLMapV3{items: map[string]item{}, mux: new(sync.RWMutex), stop: make(chan bool), now: time.Now().UnixNano()}
	go tm.clean(cleanTick)
	go tm.updateNow(time.Second)
	tmw := &MapWrapper{TTLMap: tm}
	runtime.SetFinalizer(tmw, func(_ interface{}) { tm.stop <- true })
	return tmw
}

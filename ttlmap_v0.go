package ttlmap

import (
	"sync"
	"time"
)

type TTLMapV0 struct {
	items map[string]item
	mux   *sync.Mutex
}

func (tm *TTLMapV0) Get(key string) (interface{}, bool) {
	tm.mux.Lock()
	defer tm.mux.Unlock()

	if item, ok := tm.items[key]; !ok || time.Now().UnixNano() > item.expire {
		return nil, false
	} else {
		return item.value, true
	}
}

func (tm *TTLMapV0) Set(key string, val interface{}, ex time.Duration) {
	tm.mux.Lock()
	tm.items[key] = item{ value: val, expire: time.Now().Add(ex).UnixNano()}
	tm.mux.Unlock()
}

func (tm *TTLMapV0) Len() int {
	tm.mux.Lock()
	defer tm.mux.Unlock()
	return len(tm.items)
}

func (tm *TTLMapV0) clean(tick time.Duration) {
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

func NewTTLMapV0(cleanTick time.Duration) *TTLMapV0 {
	tm := &TTLMapV0{ items: map[string]item{}, mux: new(sync.Mutex) }
	go tm.clean(cleanTick)
	return tm
}
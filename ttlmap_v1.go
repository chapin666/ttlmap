package ttlmap

import (
	"runtime"
	"sync"
	"time"
)

type TTLMapV1 struct {
	items map[string]item
	mux   *sync.Mutex
	stop  chan bool
}

func (tm *TTLMapV1) Get(key string) (interface{}, bool) {
	tm.mux.Lock()
	defer tm.mux.Unlock()

	if item, ok := tm.items[key]; !ok || time.Now().UnixNano() > item.expire {
		return nil, false
	} else {
		return item.value, true
	}
}

func (tm *TTLMapV1) Set(key string, val interface{}, ex time.Duration) {
	tm.mux.Lock()
	tm.items[key] = item{ value: val, expire: time.Now().Add(ex).UnixNano()}
	tm.mux.Unlock()
}

func (tm *TTLMapV1) Len() int {
	tm.mux.Lock()
	defer tm.mux.Unlock()
	return len(tm.items)
}

func (tm *TTLMapV1) clean(tick time.Duration) {
	for range time.Tick(tick) {
		select {
		case <-tm.stop:
			return
		default:
			tm.mux.Lock()
			for key, item := range tm.items {
				if time.Now().UnixNano() >= item.expire {
					delete(tm.items, key)
				}
			}
			tm.mux.Unlock()
		}
	}
}


func NewTTLMapV1(cleanTick time.Duration) *MapWrapper {
	tm := &TTLMapV1{items: map[string]item{}, mux: new(sync.Mutex), stop: make(chan bool)}
	go tm.clean(cleanTick)
	tmw := &MapWrapper{TTLMap: tm}
	runtime.SetFinalizer(tmw, func(_ interface{}) { tm.stop <- true })
	return tmw
}

package ttlmap

import "time"

type MapWrapper struct {
	TTLMap
}

type TTLMap interface {
	Get(key string) (interface{}, bool)
	Set(key string, val interface{}, ex time.Duration)
	Len() int
	clean(tick time.Duration)
}

type item struct {
	value  interface{}
	expire int64
}


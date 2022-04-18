package ttlmap

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkTTLMap(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	key := RandStringBytesRmndr(20)
	tasks := []struct {
		name string
		ttl  TTLMap
	}{
		{"TTLMapV0", NewTTLMapV0(time.Minute)},
		{"TTLMapV1", NewTTLMapV1(time.Minute)},
		{"TTLMapV2", NewTTLMapV2(time.Minute)},
		{"TTLMapV3", NewTTLMapV3(time.Minute)},
	}
	for _, task := range tasks {
		b.Run(task.name, func(sb *testing.B) {
			sb.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					task.ttl.Set(key, key, time.Minute)
					task.ttl.Get(key)
				}
			})
		})
	}
}

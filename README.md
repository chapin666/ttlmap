# ttlmap

a Go package that provides an in-memory key-value cache for storing TTL-based expirable items.

## benchmark
```
goos: darwin

goarch: amd64

pkg: ttlmap

cpu: Intel(R) Core(TM) i7-4770HQ CPU @ 2.20GHz

BenchmarkTTLMap/TTLMapV0-8               1739791               694.5 ns/op

BenchmarkTTLMap/TTLMapV1-8               1692664               683.0 ns/op

BenchmarkTTLMap/TTLMapV2-8               3029880               402.4 ns/op

BenchmarkTTLMap/TTLMapV3-8               5717553               214.1 ns/op

PASS

ok      ttlmap  7.504s
```

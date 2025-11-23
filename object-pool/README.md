# Object Pooling (Reducing Alloc/GC Pressure)

- Goal: show how reusing buffers via `sync.Pool` cuts allocations for hot paths like WAL/log pages (16KB here) compared to always `make`ing new slices.
- Why it matters in high-load systems: frequent allocations trigger GC and cache churn; pooling stabilizes latency and CPU when handling 100K+ requests or log records per second.
- What to look at: `bench_object_pool_test.go` contrasts pooled vs non-pooled, single-threaded and parallel.
- Try it: `go test -bench . -benchmem`.

# Test results
```
goos: darwin
goarch: arm64
pkg: github.com/creotiv/go-hiload/object-pool
cpu: Apple M4 Max
BenchmarkNoPoolReal-16                   1165755              1033 ns/op           16384 B/op          1 allocs/op
BenchmarkPoolReal-16                    76841186                15.12 ns/op           24 B/op          1 allocs/op
BenchmarkNoPoolParallelReal-16            821582              2072 ns/op           16384 B/op          1 allocs/op
BenchmarkPoolParallelReal-16            39674558                28.66 ns/op           24 B/op          1 allocs/op
```
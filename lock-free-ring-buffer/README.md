# Lock-Free Ring Buffer (SPSC)

- Goal: contrast Go channels with a single-producer/single-consumer ring that uses atomics and cache-line padding instead of locks.
- Why it matters in high-load systems: channels add scheduling, contention, and allocations that show up under 100K+ msg/sec. A bounded lock-free ring gives deterministic capacity, lower latency, and fewer GC triggers.
- What to look at: `bench_channel_vs_spsc_ring_buffer_test.go` benchmarks channel vs ring throughput.
- Try it: `go test -bench . -benchmem`.

# Test Results
```
goos: darwin
goarch: arm64
pkg: github.com/creotiv/go-hiload/lock-free-ring-buffer
cpu: Apple M4 Max
BenchmarkChannel-16             44552839                27.06 ns/op            0 B/op          0 allocs/op
BenchmarkRingBuffer-16          270782535                4.776 ns/op           0 B/op          0 allocs/op
```
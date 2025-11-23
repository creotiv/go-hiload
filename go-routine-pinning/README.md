# Goroutine Pinning (Avoiding CPU Migration)

- Goal: illustrate how `runtime.LockOSThread` keeps a hot loop on one core instead of migrating across CPUs under scheduler pressure.
- Why it matters in high-load systems: migration trashes CPU caches and TLBs, inflating latency for tight compute loops (crypto, compression, per-connection state machines). Pinning stabilizes p99s when work is cache-sensitive.
- What to look at: `bench_go_routine_pinning_test.go` compares pinned vs unpinned under an artificial scheduler hammer.
- Try it: `go test -bench . -benchmem`.

# Test results
```
goos: darwin
goarch: arm64
pkg: github.com/creotiv/go-hiload/go-routine-pinning
cpu: Apple M4 Max
BenchmarkNoPinning-16              28322             56146 ns/op               0 B/op          0 allocs/op
BenchmarkPinned-16                 10000            122112 ns/op               0 B/op          0 allocs/op
```
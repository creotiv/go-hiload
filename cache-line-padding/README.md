# Cache-Line Padding (False Sharing Demo)

- Goal: show how sharing cache lines between producer/consumer counters causes heavy cache-coherency traffic. Padding isolates `head`/`tail` so they live on separate cache lines.
- Why it matters in high-load systems: false sharing turns a cheap atomic into a cross-core ping-pong that tanks throughput and inflates tail latency when rings or queues are polled at millions of ops/sec.
- What to look at: `bench_spsc_ring_test.go` benchmarks a single-producer/single-consumer ring with and without padding.
- Try it: from this folder run `go test -bench . -benchmem`.

# Test results
```
goos: darwin
goarch: arm64
pkg: github.com/creotiv/go-hiload/cache-line-padding
cpu: Apple M4 Max
BenchmarkRingNoPad-16           127307410                8.814 ns/op           0 B/op          0 allocs/op
BenchmarkRingPad-16             408688200                2.955 ns/op           0 B/op          0 allocs/op
```
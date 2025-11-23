# Zero-Allocation JSON Parsing

- Goal: parse batches of structured logs without heap churn by avoiding reflection and allocations; keeps hot paths GC-light.
- Why it matters in high-load systems: at 100K–1M events/sec, per-object allocations and map decoding explode GC pause times and CPU. Zero/low-allocation parsers keep latency predictable and memory stable.
- What to look at: `bench_npjson_parser_test.go` compares stdlib decoding vs hand-rolled, `fastjson`, and `jsoniter`.
- Try it: `go test -bench . -benchmem`.

# Test Results
```
goos: darwin
goarch: arm64
pkg: github.com/creotiv/go-hiload/zero-allocation-parsing
cpu: Apple M4 Max
BenchmarkStdJSON-16                         7724            156253 ns/op           72113 B/op       1230 allocs/op
BenchmarkZeroAllocationJSON-16             71961             16353 ns/op           35792 B/op        609 allocs/op
BenchmarkFastZeroAllocation-16             23970             50224 ns/op          245676 B/op       1222 allocs/op
BenchmarkJsonIter-16                       61417             19601 ns/op           35839 B/op        610 allocs/op
```
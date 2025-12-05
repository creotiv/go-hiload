# Wire vs Container vs Dig

- Goal: compare compile-time DI (Wire), a hand-rolled container, and reflection-based DI (Uber Dig).
- Why it matters in high-load systems: reflection-heavy containers add hundreds of allocations and microseconds per resolve; compile-time wiring keeps hot-path latency and GC pressure near zero.
- What to look at: `bench_di_test.go` benchmarks Wire vs container singletons vs new instances vs Dig; `di.go` has the manual container; `wire.go`/`wire_gen.go` show Wire’s generated path.
- Try it: `go test -bench . -benchmem`.

# Test Results
```
goos: darwin
goarch: arm64
pkg: wirebnech
cpu: Apple M4 Max
BenchmarkWireBuild-16                   434179777                2.670 ns/op           0 B/op          0 allocs/op
BenchmarkContainerSingleton-16          562769953                2.116 ns/op           0 B/op          0 allocs/op
BenchmarkContainerNewInstances-16        51522139               21.28 ns/op           48 B/op          2 allocs/op
BenchmarkDigSingleton-16                  1257747              953.7 ns/op         1600 B/op         50 allocs/op
BenchmarkDigNewContainer-16                 88172            13368 ns/op        21319 B/op        272 allocs/op
```

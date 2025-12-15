# Boxing `interface{}` copies data

- Goal: show that passing large values to `interface{}` boxes a copy and can force heap allocations, while passing pointers avoids both.
- Why it matters: boxing value-types inflates CPU (32 KB memcpy here) and per-iteration allocations that quickly pressure GC; stick to pointers for large payloads in hot paths.
- What to look at: `bench_interface_value_copy_test.go` benchmarks a pointer-only call, the same pointer wrapped in `interface{}`, and a value passed to `interface{}` (copy + heap escape).
- Try it: `GOCACHE=$(pwd)/.gocache go test -bench . -benchmem`.

# Test results
```
goos: darwin
goarch: arm64
pkg: github.com/creotiv/go-hiload/interface-value-copy
cpu: Apple M4 Max
BenchmarkPointerNoInterface-16    	1000000000	         0.3494 ns/op	       0 B/op	       0 allocs/op
BenchmarkInterfacePointer-16      	1000000000	         0.4705 ns/op	       0 B/op	       0 allocs/op
BenchmarkInterfaceValueCopy-16    	  424690	      2438 ns/op	   32768 B/op	       1 allocs/op
```

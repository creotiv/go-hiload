# CPU Cache Layout (Struct of Arrays)

- Goal: show how struct-of-arrays keeps hot fields tightly packed so compute kernels pull more useful data per cache line than array-of-structs.
- Why it matters: physics/ML-style loops often touch one component at a time; AoS drags 32 bytes (`X,Y,Z,Mass`) into L1 for each particle even if only one value is needed, wasting bandwidth and cache slots.
- What to look at: `bench_cpu_l_cache_test.go` runs three passes (X/Y/Z) over 2M particles comparing AoS vs SoA layout.
- Try it: `go test -bench . -benchmem`.

# Test results
```
goos: darwin
goarch: arm64
pkg: github.com/creotiv/go-hiload/cpu-l-cache
cpu: Apple M4 Max
BenchmarkArrayOfStructs-16    	     302	   3993548 ns/op	       0 B/op	       0 allocs/op
BenchmarkStructOfArrays-16    	     310	   3855332 ns/op	       0 B/op	       0 allocs/op
```

# Panic/Defer/Recover Costs

- Goal: show overheads of `defer` in hot loops, panics vs error returns, and how simply having a `recover` present slows the fast path; plus why reusing objects beats per-iteration allocation.
- Why it matters in high-load systems: panic/defer/recover are great for ergonomics but add nanosecond-scale tax that compounds in tight loops; panics also trash branch prediction and unwind stacks. Avoiding needless `defer`/`recover` and reusing objects keeps latency flat.
- What to look at: `bench_hotpath_test.go` covers (1) defer inside a loop vs direct calls, (2) allocate per-iteration vs reuse, (3) panic+recover vs plain error return, (4) function that just contains recover vs one without.
- Try it: `go test -bench . -benchmem`.

# Test Results
```
goos: darwin
goarch: arm64
pkg: github.com/creotiv/go-hiload/panic-deffer-recover
cpu: Apple M4 Max
BenchmarkDeferInLoop-16           	  104577	     10282 ns/op	       0 B/op	       0 allocs/op
BenchmarkNoDeferInLoop-16         	 5098729	       240.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkAllocEachIteration-16    	20755260	        57.54 ns/op	     288 B/op	       1 allocs/op
BenchmarkReuseObject-16           	482392506	         2.460 ns/op	       0 B/op	       0 allocs/op
BenchmarkErrorReturn-16           	1000000000	         0.2285 ns/op	       0 B/op	       0 allocs/op
BenchmarkPanic-16                 	29578990	        40.60 ns/op	       0 B/op	       0 allocs/op
BenchmarkNoRecover-16             	1000000000	         0.2276 ns/op	       0 B/op	       0 allocs/op
BenchmarkRecover-16               	636868164	         1.979 ns/op	       0 B/op	       0 allocs/op
```

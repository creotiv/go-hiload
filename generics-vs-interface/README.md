# Generics vs `interface{}` and Type Assertions

- Goal: contrast monomorphized generics with `interface{}`/`any` code paths that need dynamic type checks.
- Why it matters: empty interfaces force values through fat pointers and RTTI lookups; each assertion or type switch adds bounds checks and branches that slow tight loops even without extra allocations.
- What to look at: `bench_generics_vs_interface_test.go` benchmarks a generic sum, iterating over `[]any` with a type assertion, a type-switch version, and a loop that boxes/unboxes each int.
- Try it: `GOCACHE=$(pwd)/.gocache go test -bench . -benchmem`.

# Why generics are faster here
- Generic functions are compiled per concrete type (`int` in this case), so the loop is just integer math with no dynamic dispatch.
- With `[]any`, each iteration must read the interface header (type + data pointer) and execute an implicit type check; the type switch adds more branching and prevents some compiler hoisting.

# Test results
```
goos: darwin
goarch: arm64
pkg: github.com/creotiv/go-hiload/generics-vs-interface
cpu: Apple M4 Max
BenchmarkGenericSum-16                     	    4496	    267355 ns/op	       0 B/op	       0 allocs/op
BenchmarkInterfaceTypeAssertion-16         	    3513	    342046 ns/op	       0 B/op	       0 allocs/op
BenchmarkInterfaceTypeSwitch-16            	    1783	    672841 ns/op	       0 B/op	       0 allocs/op
```
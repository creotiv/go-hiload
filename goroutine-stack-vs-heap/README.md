# Goroutine Stack vs Heap (Releasing Scratch Memory)

- Goal: show that goroutine-local scratch buffers stay on the stack and disappear once the goroutine ends, while heap buffers retained after goroutine completion keep the Go heap inflated.
- Why it matters in high-load systems: stacks are recycled as goroutines finish, keeping RSS flat even under bursts; keeping heap-backed buffers alive means that memory is stuck in the heap and rarely returned to the OS, leading to creeping resident usage.
- What to look at: `bench_stack_vs_heap_test.go` contrasts stack-local 32KB buffers vs the same buffers deliberately kept on the heap after each goroutine exits. The `heap_inuse_bytes` metric shows how much memory survives the goroutine lifetime.
- Try it: `go test -bench . -benchmem`.

# Test results
```
goos: darwin
goarch: arm64
pkg: github.com/creotiv/go-hiload/goroutine-stack-vs-heap
cpu: Apple M4 Max
BenchmarkGoroutineHeapBuffersRetained-16            1756            658926 ns/op           8388608 heap_inuse_bytes       23479560 process_sys_bytes             0 stack_inuse_bytes     8401687 B/op        520 allocs/op
BenchmarkGoroutineStackBuffers-16                   4273            268725 ns/op                 0 heap_inuse_bytes     >>23479560<<SAME process_sys_bytes             0 stack_inuse_bytes        4135 B/op        257 allocs/op
```

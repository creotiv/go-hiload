# Go High-Load Patterns (Bench Collection)

Small, focused code that demonstrate common high-load engineering patterns and why they matter. Each directory has a README with details and a `go test -bench . -benchmem` runner.

- [cache-line-padding](cache-line-padding/README.md) — avoids false sharing between producer/consumer counters in a ring buffer to cut cache-coherency ping-pong.
- [go-routine-pinning](go-routine-pinning/README.md) — uses `runtime.LockOSThread` to stop goroutine migrations that wreck cache/TLB locality in tight loops.
- [lock-free-ring-buffer](lock-free-ring-buffer/README.md) — single-producer/single-consumer ring using atomics instead of channels for predictable, low-latency queues.
- [o-direct](o-direct/README.md) — shows why buffered disk I/O is risky for WAL/logs and how O_DIRECT + fsync stabilizes durability and latency.
- [zero-allocation-parsing](zero-allocation-parsing/README.md) — zero/low-allocation JSON parsing to keep GC and CPU stable at high event rates.
- [object-pool](object-pool/README.md) — reuses fixed-size buffers via `sync.Pool` to cut allocations and GC churn on hot paths (e.g., WAL/log pages).

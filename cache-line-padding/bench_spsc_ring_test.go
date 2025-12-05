package cachelinepadding

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

// --- Section: No padding ---

type RingNoPad struct {
	head atomic.Uint64
	tail atomic.Uint64
	buf  [1024]uint64
}

func (r *RingNoPad) Enqueue(v uint64) {
	h := r.head.Load()
	r.buf[h%1024] = v
	r.head.Store(h + 1)
}

func (r *RingNoPad) Dequeue() uint64 {
	t := r.tail.Load()
	v := r.buf[t%1024]
	r.tail.Store(t + 1)
	return v
}

// --- Section: With padding ---

type RingPad struct {
	_    [64]byte
	head atomic.Uint64
	_    [56]byte
	tail atomic.Uint64
	_    [56]byte
	buf  [1024]uint64
}

func (r *RingPad) Enqueue(v uint64) {
	h := r.head.Load()
	r.buf[h%1024] = v
	r.head.Store(h + 1)
}

func (r *RingPad) Dequeue() uint64 {
	t := r.tail.Load()
	v := r.buf[t%1024]
	r.tail.Store(t + 1)
	return v
}

// --- Section: Bench core ---

type ringIface interface {
	Enqueue(uint64)
	Dequeue() uint64
}

func benchRing(b *testing.B, ring ringIface) {
	// Let Go use multiple OS threads
	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.LockOSThread()

	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)

	// Producer
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < b.N; i++ {
			ring.Enqueue(uint64(i))
		}
	}()

	// Consumer
	go func() {
		defer wg.Done()
		<-start
		var sum uint64
		for i := 0; i < b.N; i++ {
			sum += ring.Dequeue()
		}
		// Prevent compiler from optimizing the loop away
		runtime.KeepAlive(sum)
	}()

	b.ResetTimer()
	close(start) // start both goroutines "at once"
	wg.Wait()    // wait until they finish
	b.StopTimer()

	// Prevent ring itself from being optimized away
	runtime.KeepAlive(ring)
}

// --- Section: Benchmarks ---

func BenchmarkRingNoPad(b *testing.B) {
	r := &RingNoPad{}
	benchRing(b, r)
}

func BenchmarkRingPad(b *testing.B) {
	r := &RingPad{}
	benchRing(b, r)
}

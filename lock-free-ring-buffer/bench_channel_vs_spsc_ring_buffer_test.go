package lockfreeringbuffer

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

// --- Section: Lock-free SPSC padded ring buffer ---

type Ring struct {
	_    [64]byte
	head atomic.Uint64
	_    [56]byte
	tail atomic.Uint64
	_    [56]byte
	buf  []uint64
	mask uint64
}

func NewRing(size int) *Ring {
	// must be pow2
	buf := make([]uint64, size)
	return &Ring{
		buf:  buf,
		mask: uint64(size - 1),
	}
}

func (r *Ring) Enqueue(v uint64) {
	h := r.head.Load()
	r.buf[h&r.mask] = v
	r.head.Store(h + 1)
}

func (r *Ring) Dequeue() uint64 {
	t := r.tail.Load()
	v := r.buf[t&r.mask]
	r.tail.Store(t + 1)
	return v
}

// --- Section: Channel bench ---

func benchChannel(b *testing.B, ch chan uint64) {
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)

	// Producer
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < b.N; i++ {
			ch <- uint64(i)
		}
	}()

	// Consumer
	go func() {
		defer wg.Done()
		<-start
		var sum uint64
		for i := 0; i < b.N; i++ {
			sum += <-ch
		}
		runtime.KeepAlive(sum)
	}()

	b.ResetTimer()
	close(start)
	wg.Wait()
}

func BenchmarkChannel(b *testing.B) {
	// Bounded buffered channel
	ch := make(chan uint64, 1024)
	benchChannel(b, ch)
}

// --- Section: Ring buffer bench ---

func benchRing(b *testing.B, r *Ring) {
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)

	// Producer
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < b.N; i++ {
			r.Enqueue(uint64(i))
		}
	}()

	// Consumer
	go func() {
		defer wg.Done()
		<-start
		var sum uint64
		for i := 0; i < b.N; i++ {
			sum += r.Dequeue()
		}
		runtime.KeepAlive(sum)
	}()

	b.ResetTimer()
	close(start)
	wg.Wait()
}

func BenchmarkRingBuffer(b *testing.B) {
	r := NewRing(1024)
	benchRing(b, r)
}

package panicdefferrecover

import (
	"errors"
	"testing"
)

// ------------------------------------------------------------
// 1. defer inside loop vs direct call
// ------------------------------------------------------------

func dummy() {}

func BenchmarkDeferInLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		func() {
			for j := 0; j < 1000; j++ {
				defer dummy()
			}
		}()
	}
}

func BenchmarkNoDeferInLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			dummy()
		}
	}
}

// ------------------------------------------------------------
// 2. allocate per-iteration vs reuse object
// ------------------------------------------------------------

var Sink *Obj // forces escape

type Obj struct {
	A int
	B int
	C [256]byte
}

// BAD: allocate every iteration
func BenchmarkAllocEachIteration(b *testing.B) {
	for i := 0; i < b.N; i++ {
		o := &Obj{}
		// hack so Go GC do not optimize and remove it
		// making test not usable
		Sink = o
	}
}

// GOOD: reuse same object without allocation
func BenchmarkReuseObject(b *testing.B) {
	o := &Obj{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// reset the object instead of allocating
		o.A = 0
		o.B = 0
		for j := range o.C {
			o.C[j] = 0
		}
		Sink = o
	}
}

// ------------------------------------------------------------
// 3. panic vs error-return
// ------------------------------------------------------------

func errFunc() error {
	return errors.New("fail")
}

func panicFunc() {
	panic("fail")
}

func BenchmarkErrorReturn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = errFunc()
	}
}

func BenchmarkPanic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		func() {
			defer func() {
				_ = recover()
			}()
			panicFunc()
		}()
	}
}

// ------------------------------------------------------------
// 4. function *containing recover* vs same function without recover
// ------------------------------------------------------------

func doNothing() {}

func doNothingRecover() {
	defer func() {
		_ = recover()
	}()
}

func BenchmarkNoRecover(b *testing.B) {
	for i := 0; i < b.N; i++ {
		doNothing()
	}
}

func BenchmarkRecover(b *testing.B) {
	for i := 0; i < b.N; i++ {
		doNothingRecover()
	}
}

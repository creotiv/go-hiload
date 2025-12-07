package cpulcache

import (
	"runtime"
	"testing"
)

// Particle is the classic array-of-structs layout.
type Particle struct {
	X    float64
	Y    float64
	Z    float64
	Mass float64
}

// ParticleSoA stores each field in its own slice (struct-of-arrays).
type ParticleSoA struct {
	Xs   []float64
	Ys   []float64
	Zs   []float64
	Mass []float64
}

const particleCount = 1 << 21 // 2,097,152 elements -> ~64 MB AoS vs ~16 MB SoA per component stream

func makeAoS() []Particle {
	particles := make([]Particle, particleCount)
	for i := range particles {
		v := float64(i%1024) + 0.1
		particles[i] = Particle{
			X:    v,
			Y:    v * 2,
			Z:    v * 3,
			Mass: 1.0,
		}
	}
	return particles
}

func makeSoA() ParticleSoA {
	xs := make([]float64, particleCount)
	ys := make([]float64, particleCount)
	zs := make([]float64, particleCount)
	ms := make([]float64, particleCount)

	for i := 0; i < particleCount; i++ {
		v := float64(i%1024) + 0.1
		xs[i] = v
		ys[i] = v * 2
		zs[i] = v * 3
		ms[i] = 1.0
	}

	return ParticleSoA{Xs: xs, Ys: ys, Zs: zs, Mass: ms}
}

func BenchmarkArrayOfStructs(b *testing.B) {
	runtime.GC()
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	_ = makeAoS()
	runtime.ReadMemStats(&after)
	b.Logf("bytes allocated: %d", after.TotalAlloc-before.TotalAlloc)
	b.ReportAllocs()
	b.ResetTimer()
	particles := makeAoS()

	for i := 0; i < b.N; i++ {
		for j := 0; j < len(particles); j++ {
			particles[j].X += 1
		}
	}
	runtime.KeepAlive(particles)
}

func BenchmarkStructOfArrays(b *testing.B) {
	runtime.GC()
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	_ = makeSoA()
	runtime.ReadMemStats(&after)
	b.Logf("bytes allocated: %d", after.TotalAlloc-before.TotalAlloc)
	b.ReportAllocs()
	b.ResetTimer()
	particles := makeSoA()

	for i := 0; i < b.N; i++ {
		for j := 0; j < len(particles.Xs); j++ {
			particles.Xs[j] += 1
		}

	}
	runtime.KeepAlive(particles)
}

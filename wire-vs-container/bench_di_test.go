package wirebnech

import (
	"runtime"
	"testing"

	"go.uber.org/dig"
)

var benchConfig = &Config{
	DBDSN:     "postgres://user:pass@localhost:5432/app",
	RedisAddr: "redis://localhost:6379/0",
}

func BenchmarkWireBuild(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		app := InitializeApp(benchConfig)
		if app.DB == nil || app.Redis == nil {
			b.Fatal("wire build returned nil services")
		}
		runtime.KeepAlive(app)
	}
}

func BenchmarkContainerSingleton(b *testing.B) {
	di := NewContainer(benchConfig)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		db := di.GetDB(false)
		redis := di.GetRedis(false)
		if db == nil || redis == nil {
			b.Fatal("container returned nil services")
		}
	}

	runtime.KeepAlive(di)
}

func BenchmarkContainerNewInstances(b *testing.B) {
	di := NewContainer(benchConfig)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		db := di.GetDB(true)
		redis := di.GetRedis(true)
		if db == nil || redis == nil {
			b.Fatal("container returned nil services")
		}
		runtime.KeepAlive(db)
		runtime.KeepAlive(redis)
	}

	runtime.KeepAlive(di)
}

func BenchmarkDigSingleton(b *testing.B) {
	container := newDigContainer(b)

	// Warm up constructors so the timer excludes first-time creation.
	if err := container.Invoke(func(*DB, *Redis) {}); err != nil {
		b.Fatalf("dig warmup: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var app *App

		if err := container.Invoke(func(db *DB, redis *Redis) {
			app = &App{DB: db, Redis: redis}
		}); err != nil {
			b.Fatalf("dig invoke: %v", err)
		}

		runtime.KeepAlive(app)
	}
}

func BenchmarkDigNewContainer(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		container := newDigContainer(b)

		var app *App

		if err := container.Invoke(func(db *DB, redis *Redis) {
			app = &App{DB: db, Redis: redis}
		}); err != nil {
			b.Fatalf("dig invoke: %v", err)
		}

		runtime.KeepAlive(app)
	}
}

func newDigContainer(b *testing.B) *dig.Container {
	container := dig.New()

	if err := container.Provide(func() *Config { return benchConfig }); err != nil {
		b.Fatalf("dig provide config: %v", err)
	}
	if err := container.Provide(NewDB); err != nil {
		b.Fatalf("dig provide db: %v", err)
	}
	if err := container.Provide(NewRedis); err != nil {
		b.Fatalf("dig provide redis: %v", err)
	}

	return container
}

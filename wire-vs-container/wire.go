//go:build wireinject
// +build wireinject

package wirebnech

import "github.com/google/wire"

// InitializeApp shows how Wire stitches dependencies together at compile time.
func InitializeApp(cfg *Config) *App {
	wire.Build(NewDB, NewRedis, wire.Struct(new(App), "*"))
	return &App{}
}

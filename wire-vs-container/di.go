package wirebnech

import (
	"sync"
	"sync/atomic"
)

// Config represents the inputs shared by all services.
type Config struct {
	DBDSN     string
	RedisAddr string
}

var (
	dbCounter    atomic.Int64
	redisCounter atomic.Int64
)

type DB struct {
	DSN string
	ID  int64
}

func NewDB(cfg *Config) *DB {
	return &DB{
		DSN: cfg.DBDSN,
		ID:  dbCounter.Add(1),
	}
}

type Redis struct {
	Addr string
	ID   int64
}

func NewRedis(cfg *Config) *Redis {
	return &Redis{
		Addr: cfg.RedisAddr,
		ID:   redisCounter.Add(1),
	}
}

// App glues services together the way Wire would.
type App struct {
	DB    *DB
	Redis *Redis
}

// Container lazily constructs services and can return singletons or new instances.
type Container struct {
	cfg       *Config
	dbOnce    sync.Once
	db        *DB
	redisOnce sync.Once
	redis     *Redis
}

func NewContainer(cfg *Config) *Container {
	return &Container{cfg: cfg}
}

func (c *Container) GetDB(newInstance bool) *DB {
	if newInstance {
		return NewDB(c.cfg)
	}

	c.dbOnce.Do(func() {
		c.db = NewDB(c.cfg)
	})

	return c.db
}

func (c *Container) GetRedis(newInstance bool) *Redis {
	if newInstance {
		return NewRedis(c.cfg)
	}

	c.redisOnce.Do(func() {
		c.redis = NewRedis(c.cfg)
	})

	return c.redis
}

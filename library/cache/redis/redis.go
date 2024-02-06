package redis

import (
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr         string
	ClientName   string
	Username     string
	Password     string
	Db           int
	MaxRetries   int
	MinIdleConns int
	MaxIdleConns int
}

func NewRedisClient(c *Config) *redis.Client {
	options := &redis.Options{
		Addr:         c.Addr,
		ClientName:   c.ClientName,
		Username:     c.Username,
		Password:     c.Password,
		DB:           c.Db,
		MaxRetries:   c.MaxRetries,
		MinIdleConns: c.MinIdleConns,
		MaxIdleConns: c.MaxIdleConns,
	}

	client := redis.NewClient(options)
	return client
}

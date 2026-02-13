package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/adityayuga/go-sdk/errors"

	redis "github.com/redis/go-redis/v9"
)

type (
	Cache struct {
		rdb *redis.Client
	}

	CacheMethod interface {
		Set(ctx context.Context, key string, val string, ttlInSecond int) error
		Get(ctx context.Context, key string) (string, error)
		Del(ctx context.Context, key string) error
		Exist(ctx context.Context, key string) (bool, error)
		Incr(ctx context.Context, key string) (int64, error)

		Close(ctx context.Context) error
	}

	Config struct {
		URL      string `json:"url"`      // host or ip address
		Port     string `json:"port"`     // default 6379
		Password string `json:"password"` // default not using password
		DB       int    `json:"db"`       // 0 = use default DB
	}
)

func New(ctx context.Context, cfg Config) (CacheMethod, error) {
	if cfg.Port == "" {
		cfg.Port = "6379"
	}

	if cfg.URL == "" {
		return nil, errors.New("[pkg cache] Empty url")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.URL, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB, // 0 = use default DB
	})

	return Cache{
		rdb: rdb,
	}, nil
}

// Set set key value with ttl in second
func (c Cache) Set(ctx context.Context, key string, val string, ttlInSecond int) error {
	return c.rdb.Set(ctx, key, val, time.Duration(ttlInSecond)*time.Second).Err()
}

// Get get value by key
func (c Cache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return val, nil
}

// Del delete key
func (c Cache) Del(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, key).Err()
}

// Exist check if key exist
func (c Cache) Exist(ctx context.Context, key string) (bool, error) {
	result, err := c.rdb.Exists(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return false, err
	}

	if result != 1 {
		return false, nil
	}

	return true, nil
}

// Incr increment the value for key by 1, return new value and error
func (c Cache) Incr(ctx context.Context, key string) (int64, error) {
	return c.rdb.Incr(ctx, key).Result()
}

func (c Cache) Close(ctx context.Context) error {
	return c.rdb.Close()
}

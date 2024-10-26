package cache

import (
	"context"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	cache *cache.Cache
}

func NewRedisCache(client *redis.Client) *RedisCache {
	cache := cache.New(&cache.Options{
		Redis:      client,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	return &RedisCache{
		cache: cache,
	}
}

func (c *RedisCache) Get(key string, wanted interface{}, ctx context.Context) error {
	err := c.cache.Get(ctx, key, &wanted)
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisCache) Put(
	key string,
	value interface{},
	expiresAfter time.Duration,
	ctx context.Context,
) error {
	return c.cache.Set(&cache.Item{
		Key:   key,
		Value: value,
		TTL:   expiresAfter,
		Ctx:   ctx,
	})
}

func (c *RedisCache) Delete(key string, ctx context.Context) error {
	return c.cache.Delete(ctx, key)
}

type TypedRedisCache[T any] struct {
	cache *cache.Cache
}

func NewTypedRedisCache[T any](
	client *redis.Client,
) *TypedRedisCache[T] {
	cache := cache.New(&cache.Options{
		Redis:      client,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	return &TypedRedisCache[T]{
		cache: cache,
	}
}

func (c *TypedRedisCache[T]) Get(key string, wanted *T, ctx context.Context) error {
	err := c.cache.Get(ctx, key, &wanted)
	if err != nil {
		return err
	}

	return nil
}

func (c *TypedRedisCache[T]) Put(
	key string,
	value *T,
	expiresAfter time.Duration,
	ctx context.Context,
) error {
	return c.cache.Set(&cache.Item{
		Key:   key,
		Value: value,
		TTL:   expiresAfter,
		Ctx:   ctx,
	})
}

func (c *TypedRedisCache[T]) Delete(key string, ctx context.Context) error {
	return c.cache.Delete(ctx, key)
}

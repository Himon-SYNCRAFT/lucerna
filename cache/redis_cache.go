package cache

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	cache       *cache.Cache
	redisClient *redis.Client
	poolName    string
	prefix      string
}

type TypedRedisCache[T any] struct {
	cache       *cache.Cache
	redisClient *redis.Client
	poolName    string
	prefix      string
}

type CacheOptions struct {
	prefix string
}

func NewRedisCache(
	client *redis.Client,
	poolName string,
	options *CacheOptions,
) (Cache, error) {
	if client == nil {
		return nil, errors.New("redis client is nil")
	}

	cache := cache.New(&cache.Options{
		Redis:      client,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	redisCache := &RedisCache{
		cache:       cache,
		redisClient: client,
		poolName:    poolName,
	}

	if options != nil {
		redisCache.prefix = options.prefix
	}

	return redisCache, nil
}

func (c *RedisCache) getKey(key string) string {
	result := key

	if c.poolName != "" {
		result = c.poolName + "." + key
	}

	if c.prefix != "" {
		result = c.prefix + "." + result
	}

	return result
}

func (c *RedisCache) Clear(ctx context.Context) error {
	key := "*"

	if c.poolName != "" {
		key = c.poolName + "." + key
	}

	if c.prefix != "" {
		key = c.prefix + "." + key
	}

	log.Print("Clearing cache: " + key)

	var cursor uint64

	keys, cursor, err := c.redisClient.Scan(ctx, cursor, key, 0).Result()
	if err != nil {
		return err
	}

	result, err := c.redisClient.Del(ctx, keys...).Result()

	for _, k := range keys {
		c.cache.DeleteFromLocalCache(k)
	}

	log.Print(fmt.Sprintf("Cleared cache keys: %s, %d", keys, result))
	return err
}

func (c *RedisCache) Get(key string, wanted interface{}, ctx context.Context) error {
	err := c.cache.Get(ctx, c.getKey(key), &wanted)
	if err != nil {
		log.Print("Cache Get (miss): " + c.getKey(key))
		return err
	}

	log.Print("Cache Get (hit): " + c.getKey(key))
	return nil
}

func (c *RedisCache) Put(
	key string,
	value interface{},
	expiresAfter time.Duration,
	ctx context.Context,
) error {
	log.Print("Cache Put: " + c.getKey(key))
	return c.cache.Set(&cache.Item{
		Key:   c.getKey(key),
		Value: value,
		TTL:   expiresAfter,
		Ctx:   ctx,
	})
}

func (c *RedisCache) Delete(key string, ctx context.Context) error {
	log.Print("Cache Delete: " + c.getKey(key))
	return c.cache.Delete(ctx, c.getKey(key))
}

func (c *TypedRedisCache[T]) getKey(key string) string {
	result := key

	if c.poolName != "" {
		result = c.poolName + "." + key
	}

	if c.prefix != "" {
		result = c.prefix + "." + result
	}

	return result
}

func NewTypedRedisCache[T any](
	client *redis.Client,
	poolName string,
	options *CacheOptions,
) (TypedCache[T], error) {
	if client == nil {
		return nil, errors.New("redis client is nil")
	}

	cache := cache.New(&cache.Options{
		Redis:      client,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	redisCache := &TypedRedisCache[T]{
		cache:       cache,
		redisClient: client,
		poolName:    poolName,
	}

	if options != nil {
		redisCache.prefix = options.prefix
	}

	return redisCache, nil
}

func (c *TypedRedisCache[T]) Get(key string, wanted *T, ctx context.Context) error {
	err := c.cache.Get(ctx, c.getKey(key), &wanted)
	if err != nil {
		log.Print("Cache Get (miss): " + c.getKey(key))
		return err
	}

	log.Print("Cache Get (hit): " + c.getKey(key))
	return nil
}

func (c *TypedRedisCache[T]) Put(
	key string,
	value *T,
	expiresAfter time.Duration,
	ctx context.Context,
) error {
	log.Print("Cache Put: " + c.getKey(key))
	return c.cache.Set(&cache.Item{
		Key:   c.getKey(key),
		Value: value,
		TTL:   expiresAfter,
		Ctx:   ctx,
	})
}

func (c *TypedRedisCache[T]) Delete(key string, ctx context.Context) error {
	log.Print("Cache Delete: " + c.getKey(key))
	return c.cache.Delete(ctx, c.getKey(key))
}

func (c *TypedRedisCache[T]) Clear(ctx context.Context) error {
	key := "*"

	if c.poolName != "" {
		key = c.poolName + "." + key
	}

	if c.prefix != "" {
		key = c.prefix + "." + key
	}

	log.Print("Clearing cache: " + key)

	var cursor uint64

	keys, cursor, err := c.redisClient.Scan(ctx, cursor, key, 0).Result()
	if err != nil {
		return err
	}

	_, err = c.redisClient.Del(ctx, keys...).Result()

	for _, k := range keys {
		c.cache.DeleteFromLocalCache(k)
	}

	log.Print("Cleared cache: " + key)
	return err
}

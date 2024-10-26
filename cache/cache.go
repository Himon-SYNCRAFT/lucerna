package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(key string, wanted interface{}, ctx context.Context) error
	Put(key string, value interface{}, expiresAfter time.Duration, ctx context.Context) error
	Delete(key string, ctx context.Context) error
}

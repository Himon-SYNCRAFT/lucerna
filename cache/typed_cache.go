package cache

import (
	"context"
	"time"
)

type TypedCache[T any] interface {
	Get(key string, wanted *T, ctx context.Context) error
	Put(key string, value *T, expiresAfter time.Duration, ctx context.Context) error
	Delete(key string, ctx context.Context) error
	Clear(ctx context.Context) error
}

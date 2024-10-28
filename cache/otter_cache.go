package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/maypok86/otter"
)

type OtterCache struct {
	cache otter.CacheWithVariableTTL[string, []byte]
}

func NewOtterCache() (Cache, error) {
	cache, err := otter.MustBuilder[string, []byte](10000).
		CollectStats().
		Cost(func(key string, value []byte) uint32 {
			return 1
		}).
		WithVariableTTL().
		Build()
	if err != nil {
		return nil, err
	}

	otter := &OtterCache{
		cache: cache,
	}

	return otter, nil
}

func (c *OtterCache) Get(key string, wanted interface{}, ctx context.Context) error {
	result, inCache := c.cache.Get(key)

	if !inCache {
		return errors.New("cache miss")
	}

	json.Unmarshal(result, &wanted)

	fmt.Println(fmt.Sprintf("get wanted %v", wanted))

	return nil
}

func (c *OtterCache) Put(
	key string,
	value interface{},
	expiresAfter time.Duration,
	ctx context.Context,
) error {
	var bytes []byte
	var err error

	bytes, err = json.Marshal(value)
	if err != nil {
		return err
	}

	ok := c.cache.Set(key, bytes, expiresAfter)

	if !ok {
		return errors.New("Cache error. Cannot set cache.")
	}

	return nil
}

func (c *OtterCache) Delete(key string, ctx context.Context) error {
	c.cache.Delete(key)
	return nil
}

func (c *OtterCache) Clear(ctx context.Context) error {
	c.cache.Clear()
	return nil
}

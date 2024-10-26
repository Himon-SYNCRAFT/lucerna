package main

import (
	"context"
	"fmt"
	"time"

	lucerna_cache "lucerna/cache"
	"lucerna/dotenv"
)

type Object struct {
	Name string
}

func main() {
	env := dotenv.LoadEnv()

	var cache lucerna_cache.Cache

	cache = lucerna_cache.NewRedisCache(
		env.Redis.Host,
		env.Redis.Port,
		env.Redis.User,
		env.Redis.Password,
	)

	value := Object{
		Name: "test",
	}

	ctx := context.Background()

	key := "test"

	err := cache.Put(key, value, time.Hour, ctx)
	if err != nil {
		panic(err)
	}

	var value2 Object

	err = cache.Get(key, &value2, ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(value.Name)
	fmt.Println(value2.Name)

	err = cache.Delete(key, ctx)
	if err != nil {
		panic(err)
	}

	var typedCache lucerna_cache.TypedCache[Object]
	typedCache = lucerna_cache.NewTypedRedisCache[Object](
		env.Redis.Host,
		env.Redis.Port,
		env.Redis.User,
		env.Redis.Password,
	)

	err = typedCache.Put(key, value, time.Hour, ctx)
	if err != nil {
		panic(err)
	}

	var value3 Object

	err = typedCache.Get(key, &value3, ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(value.Name)
	fmt.Println(value3.Name)

	err = typedCache.Delete(key, ctx)
	if err != nil {
		panic(err)
	}
}

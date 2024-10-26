package cache

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/redis/go-redis/v9"

	"lucerna/dotenv"
)

var redisConf = dotenv.Redis{
	Host:     "localhost",
	Port:     "6379",
	User:     "",
	Password: "",
}

const redisHostPort = "16379"

type testObject struct {
	Key  string
	Name string
}

func newTestObject(key string, name string) *testObject {
	return &testObject{
		Key:  key,
		Name: name,
	}
}

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	redisOpts := dockertest.RunOptions{
		Repository: "redis",
		Tag:        "alpine",
		Env: []string{
			"REDIS_HOST=" + redisConf.Host,
			"REDIS_PORT=" + redisConf.Port,
			"REDIS_USER=" + redisConf.User,
			"REDIS_PASSWORD=" + redisConf.Password,
		},
		ExposedPorts: []string{"6379"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			docker.Port(redisConf.Port): {
				{HostIP: redisConf.Host, HostPort: redisHostPort},
			},
		},
	}

	resource, err := pool.RunWithOptions(&redisOpts)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisConf.Host, redisConf.Port),
		Username: redisConf.User,
		Password: redisConf.Password,
	})

	if err = pool.Retry(func() error {
		return redisClient.Ping(context.Background()).Err()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()

	err = pool.Purge(resource)
	if err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	os.Exit(code)
}

func TestRedisCache(t *testing.T) {
	var cache Cache
	// var typedCache TypedCache[testObject]
	var err error
	obj := newTestObject("key", "test")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisConf.Host, redisHostPort),
		Username: redisConf.User,
		Password: redisConf.Password,
	})

	ctx := context.Background()

	cache = NewRedisCache(redisClient)
	// typedCache = NewTypedRedisCache[testObject](redisClient)

	err = cache.Put(obj.Key, obj, time.Minute, ctx)
	if err != nil {
		t.Error("Error putting object to cache: ", err)
	}

	var wanted testObject

	err = cache.Get(obj.Key, &wanted, ctx)
	if err != nil {
		t.Error("Error getting object from cache: ", err)
	}

	if wanted.Name != obj.Name {
		t.Error("Expected name to be ", obj.Name, " but got ", wanted.Name)
	}

	if wanted.Name != obj.Name {
		t.Error("Expected name to be ", obj.Name, " but got ", wanted.Name)
	}

	err = cache.Delete(obj.Key, ctx)
	if err != nil {
		t.Error("Error getting object from cache: ", err)
	}

	var wanted2 testObject

	err = cache.Get(obj.Key, &wanted2, ctx)
	if err == nil {
		t.Error("Error getting deleted object from cache")
	}

	if wanted2.Name != "" {
		t.Error("Expected name to be empty but got ", wanted.Name)
	}

	if wanted2.Name != "" {
		t.Error("Expected name to be empty but got ", wanted.Name)
	}
}

func TestTypedRedisCache(t *testing.T) {
	var typedCache TypedCache[testObject]
	var err error
	obj := newTestObject("key2", "test2")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisConf.Host, redisHostPort),
		Username: redisConf.User,
		Password: redisConf.Password,
	})

	ctx := context.Background()

	typedCache = NewTypedRedisCache[testObject](redisClient)

	err = typedCache.Put(obj.Key, obj, time.Minute, ctx)
	if err != nil {
		t.Error("Error putting object to cache: ", err)
	}

	var wanted testObject

	err = typedCache.Get(obj.Key, &wanted, ctx)
	if err != nil {
		t.Error("Error getting object from cache: ", err)
	}

	if wanted.Name != obj.Name {
		t.Error("Expected name to be ", obj.Name, " but got ", wanted.Name)
	}

	if wanted.Name != obj.Name {
		t.Error("Expected name to be ", obj.Name, " but got ", wanted.Name)
	}

	err = typedCache.Delete(obj.Key, ctx)
	if err != nil {
		t.Error("Error getting object from cache: ", err)
	}

	var wanted2 testObject

	err = typedCache.Get(obj.Key, &wanted2, ctx)
	if err == nil {
		t.Error("Error getting deleted object from cache")
	}

	if wanted2.Name != "" {
		t.Error("Expected name to be empty but got ", wanted.Name)
	}

	if wanted2.Name != "" {
		t.Error("Expected name to be empty but got ", wanted.Name)
	}
}

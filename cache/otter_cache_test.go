package cache

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestOtterCache(t *testing.T) {
	var cache Cache
	var err error
	obj := newTestObject("key", "test")

	ctx := context.Background()

	cache, err = NewOtterCache()
	if err != nil {
		t.Error("Error creating cache: ", err)
	}

	err = cache.Put(obj.Key, obj, time.Minute, ctx)
	if err != nil {
		t.Error("Error putting object to cache: ", err)
	}

	var wanted testObject

	err = cache.Get(obj.Key, &wanted, ctx)
	if err != nil {
		t.Error("Error getting object from cache: ", err)
	}

	fmt.Println(fmt.Sprintf("obj %v", obj))
	fmt.Println(fmt.Sprintf("wanted %v", wanted))

	if wanted.Name != obj.Name {
		t.Error("Expected name to be ", obj.Name, " but got ", wanted.Name)
	}

	if wanted.Key != obj.Key {
		t.Error("Expected key to be ", obj.Key, " but got ", wanted.Key)
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
		t.Error("Expected name to be empty but got ", wanted2.Name)
	}

	if wanted2.Key != "" {
		t.Error("Expected key to be empty but got ", wanted2.Key)
	}
}

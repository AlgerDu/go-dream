package memorycache

import (
	"fmt"
	"sync"

	"github.com/AlgerDu/go-dream/src/dinfra"
)

type (
	MemoryCacher struct {
		values   map[string]any
		lock     sync.Mutex
		keyLocks map[string]*sync.Mutex
	}
)

func NewMemoryCacher() *MemoryCacher {
	return &MemoryCacher{
		values:   map[string]any{},
		lock:     sync.Mutex{},
		keyLocks: map[string]*sync.Mutex{},
	}
}

func (cacher *MemoryCacher) Set(key string, value any, expire int64) error {
	cacher.lock.Lock()
	defer cacher.lock.Unlock()

	cacher.values[key] = value

	return nil
}

func (cacher *MemoryCacher) Get(key string) (dinfra.CacheValue, error) {
	value, exist := cacher.values[key]
	if !exist {
		return nil, dinfra.ErrNilCache
	}

	return &MemoryCacheValue{value: value}, nil
}

func (cacher *MemoryCacher) SetH(key string, values map[string]any, expire int64) error {

	cache, exist := cacher.values[key]
	if exist {
		hashValues, ok := cache.(map[string]any)
		if !ok {
			return fmt.Errorf("%s is not hash type", key)
		}
		for k, v := range values {
			hashValues[k] = v
		}
		return nil
	}

	cacher.lock.Lock()
	defer cacher.lock.Unlock()

	cacher.values[key] = values

	return nil
}

func (cacher *MemoryCacher) GetH(key string, index string) (dinfra.CacheValue, error) {
	cache, exist := cacher.values[key]
	if !exist {
		return nil, dinfra.ErrNilCache
	}

	if index == "" {
		return &MemoryCacheValue{value: cache}, nil
	}

	hashValue, ok := cache.(map[string]any)
	if !ok {
		return nil, dinfra.ErrNilCache
	}

	value, exist := hashValue[index]
	if !exist {
		return nil, dinfra.ErrNilCache
	}
	return &MemoryCacheValue{value: value}, nil
}

func (cacher *MemoryCacher) Lock(key string) {
	keyLock, exist := cacher.keyLocks[key]
	if !exist {
		cacher.lock.Lock()
		defer cacher.lock.Unlock()
		keyLock = &sync.Mutex{}
		cacher.keyLocks[key] = keyLock
	}

	keyLock.Lock()
}

func (cacher *MemoryCacher) Unlock(key string) {
	keyLock, exist := cacher.keyLocks[key]
	if !exist {
		panic(fmt.Sprintf("%s key lock not exsit", key))
	}

	keyLock.Unlock()
}

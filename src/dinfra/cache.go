package dinfra

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrNilCache = errors.New("nil cache")
)

type (
	CacheValue interface {
		To(dstType reflect.Type) (any, error)
	}

	Cacher interface {
		Set(key string, value any, expire int64) error
		Get(key string) (CacheValue, error)

		SetH(key string, values map[string]any, expire int64) error // TODO 这里不应该直接定义为 map，可能是 struct 使用 hash 存储；其实需要处理
		GetH(key string, index string) (CacheValue, error)

		Lock(key string)
		Unlock(key string)
	}

	CacheUpdateOptions struct {
		Expire int64
	}

	CacheNilHandler[ValutType any]  func() (ValutType, error)
	CacheValueGetter[ValueType any] func(nilHandler CacheNilHandler[ValueType]) (ValueType, error)

	CacheValueUpdateHandler[ValueType any] func(value ValueType) (ValueType, *CacheUpdateOptions, error)
	CacheValueUpdater[ValueType any]       func(hanler CacheValueUpdateHandler[ValueType]) error

	CacheHValueUpdateHandler[ValueType any] func(value ValueType) (map[string]any, *CacheUpdateOptions, error)
	CacheHValueUpdater[ValueType any]       func(handler CacheHValueUpdateHandler[ValueType]) error

	CacheValueUpdater2[ValueType any] func(value *ValueType, expire int64, options ...func(value ValueType) ValueType) error
)

func valueToType[ValueType any](value CacheValue) (ValueType, error) {
	var zero ValueType
	vt := reflect.TypeOf(zero)
	v, err := value.To(vt)
	if err != nil {
		return zero, err
	}

	tv, ok := v.(ValueType)
	if ok {
		return tv, nil
	}
	return zero, fmt.Errorf("value to %s error", vt.Name())
}

func UseCache[ValueType any](cacher Cacher, key string) (CacheValueGetter[ValueType], CacheValueUpdater[ValueType]) {
	getter := func(nilHandler CacheNilHandler[ValueType]) (ValueType, error) {
		cacheValue, err := cacher.Get(key)
		if err == nil {
			return valueToType[ValueType](cacheValue)
		}

		if err == ErrNilCache && nilHandler != nil {
			value, err := nilHandler()
			if err == nil {
				cacher.Set(key, value, 0)
			}
			return value, err
		}

		var zero ValueType
		return zero, err
	}

	updater := func(hanler CacheValueUpdateHandler[ValueType]) error {
		value, err := getter(nil)
		if err != nil && err != ErrNilCache {
			return err
		}

		cacher.Lock(key)
		defer cacher.Unlock(key)

		value, options, err := hanler(value)
		if err != nil {
			return err
		}

		expire := int64(0)
		if options != nil {
			expire = options.Expire
		}

		return cacher.Set(key, value, expire)
	}

	return getter, updater
}

func UseCache2[ValueType any](cacher Cacher, key string) (CacheValueGetter[ValueType], CacheValueUpdater2[ValueType]) {

	getter := func(nilHandler CacheNilHandler[ValueType]) (ValueType, error) {
		cacheValue, err := cacher.Get(key)
		if err == nil {
			return valueToType[ValueType](cacheValue)
		}

		if err == ErrNilCache && nilHandler != nil {
			value, err := nilHandler()
			if err == nil {
				cacher.Set(key, value, 0)
			}
			return value, err
		}

		var zero ValueType
		return zero, err
	}

	updater := func(value *ValueType, expire int64, options ...func(value ValueType) ValueType) error {
		cacher.Lock(key)
		defer cacher.Unlock(key)

		var toUpdate ValueType

		if value == nil {
			cacheValue, err := getter(nil)
			if err != nil {
				return err
			}
			toUpdate = cacheValue
		} else {
			toUpdate = *value
		}

		if len(options) > 1 {
			for _, handler := range options {
				toUpdate = handler(toUpdate)
			}
		}

		return cacher.Set(key, value, expire)
	}

	return getter, updater
}

func UseHCache[ValueType any](cacher Cacher, key string) (CacheValueGetter[ValueType], CacheHValueUpdater[ValueType]) {
	getter := func(nilHandler CacheNilHandler[ValueType]) (ValueType, error) {
		cacheValue, err := cacher.Get(key)
		if err == nil {
			return valueToType[ValueType](cacheValue)
		}

		if err == ErrNilCache && nilHandler != nil {
			value, err := nilHandler()
			if err == nil {
				cacher.Set(key, value, 0)
			}
			return value, err
		}

		var zero ValueType
		return zero, err
	}

	updater := func(hanler CacheHValueUpdateHandler[ValueType]) error {
		value, err := getter(nil)
		if err != nil && err != ErrNilCache {
			return err
		}

		cacher.Lock(key)
		defer cacher.Unlock(key)

		indexValues, options, err := hanler(value)
		if err != nil {
			return err
		}

		expire := int64(0)
		if options != nil {
			expire = options.Expire
		}

		return cacher.SetH(key, indexValues, expire)
	}

	return getter, updater
}

func UseHCache2[ValueType any](cacher Cacher, key string) (CacheValueGetter[ValueType], CacheValueUpdater2[ValueType]) {

	getter := func(nilHandler CacheNilHandler[ValueType]) (ValueType, error) {
		cacheValue, err := cacher.Get(key)
		if err == nil {
			return valueToType[ValueType](cacheValue)
		}

		if err == ErrNilCache && nilHandler != nil {
			value, err := nilHandler()
			if err == nil {
				cacher.Set(key, value, 0)
			}
			return value, err
		}

		var zero ValueType
		return zero, err
	}

	updater := func(value *ValueType, expire int64, options ...func(value ValueType) ValueType) error {
		cacher.Lock(key)
		defer cacher.Unlock(key)

		var toUpdate ValueType

		if value == nil {
			cacheValue, err := getter(nil)
			if err != nil {
				return err
			}
			toUpdate = cacheValue
		} else {
			toUpdate = *value
		}

		if len(options) > 1 {
			for _, handler := range options {
				toUpdate = handler(toUpdate)
			}
		}

		tmp, ok := any(toUpdate).(map[string]any)
		if !ok {
			panic("it is not map, HACK 后续修改")
		}

		return cacher.SetH(key, tmp, expire)
	}

	return getter, updater
}

func UseHCacheIndex[ValueType any](cacher Cacher, key string, index string) (CacheValueGetter[ValueType], CacheValueUpdater[ValueType]) {

	getter := func(nilHandler CacheNilHandler[ValueType]) (ValueType, error) {
		cacheValue, err := cacher.GetH(key, index)
		if err == nil {
			return valueToType[ValueType](cacheValue)
		}

		if err == ErrNilCache && nilHandler != nil {
			value, err := nilHandler()
			if err == nil {
				cacher.SetH(key, map[string]any{index: value}, 0)
			}
			return value, err
		}

		var zero ValueType
		return zero, err
	}

	updater := func(hanler CacheValueUpdateHandler[ValueType]) error {
		value, err := getter(nil)
		if err != nil && err != ErrNilCache {
			return err
		}

		cacher.Lock(key)
		defer cacher.Unlock(key)

		value, options, err := hanler(value)
		if err != nil {
			return err
		}

		expire := int64(0)
		if options != nil {
			expire = options.Expire
		}

		return cacher.SetH(key, map[string]any{index: value}, expire)
	}

	return getter, updater
}

func UseHCacheIndex2[ValueType any](cacher Cacher, key string, index string) (CacheValueGetter[ValueType], CacheValueUpdater2[ValueType]) {

	getter := func(nilHandler CacheNilHandler[ValueType]) (ValueType, error) {
		cacheValue, err := cacher.GetH(key, index)
		if err == nil {
			return valueToType[ValueType](cacheValue)
		}

		if err == ErrNilCache && nilHandler != nil {
			value, err := nilHandler()
			if err == nil {
				cacher.SetH(key, map[string]any{index: value}, 0)
			}
			return value, err
		}

		var zero ValueType
		return zero, err
	}

	updater := func(value *ValueType, expire int64, options ...func(value ValueType) ValueType) error {
		cacher.Lock(key)
		defer cacher.Unlock(key)

		var toUpdate ValueType

		if value == nil {
			cacheValue, err := getter(nil)
			if err != nil {
				return err
			}
			toUpdate = cacheValue
		} else {
			toUpdate = *value
		}

		if len(options) > 1 {
			for _, handler := range options {
				toUpdate = handler(toUpdate)
			}
		}

		return cacher.SetH(key, map[string]any{index: toUpdate}, expire)
	}

	return getter, updater
}

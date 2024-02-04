package dinfra

import "errors"

var (
	ErrNilCache = errors.New("nil cache")
)

type (
	CacheValue interface {
		BindTo(dst any) error
	}

	Cacher interface {
		Set(key string, value any, expire int64) error
		Get(key string) (CacheValue, error)

		SetH(key string, values map[string]any, expire int64) error
		GetH(key string, index string) (CacheValue, error)

		Lock(key string)
		Unlock(key string)
	}

	CacheValueUpdateHandler[ValueType any] func(value *ValueType) (*ValueType, *CacheUpdateOptions, error)

	CacheValueGetter[ValueType any]  func() (*ValueType, error)
	CacheValueUpdater[ValueType any] func(hanler CacheValueUpdateHandler[ValueType]) error

	CacheHashValueUpdateHandler[ValueType any] func(value *ValueType) (map[string]any, *CacheUpdateOptions, error)
	CacheHashValueUpdater[ValueType any]       func(handler CacheHashValueUpdateHandler[ValueType]) error

	CacheUpdateOptions struct {
		Expire int64
	}
)

func UseCache[ValueType any](cacher Cacher, key string) (CacheValueGetter[ValueType], CacheValueUpdater[ValueType]) {

	getter := func() (*ValueType, error) {
		value, err := cacher.Get(key)
		if err != nil {
			return nil, err
		}
		rst := new(ValueType)
		err = value.BindTo(rst)
		return rst, err
	}

	updater := func(hanler CacheValueUpdateHandler[ValueType]) error {
		value, err := getter()
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

func UseHashCache[ValueType any](cacher Cacher, key string) (CacheValueGetter[ValueType], CacheHashValueUpdater[ValueType]) {

	getter := func() (*ValueType, error) {
		value, err := cacher.Get(key)
		if err != nil {
			return nil, err
		}
		rst := new(ValueType)
		err = value.BindTo(rst)
		return rst, err
	}

	updater := func(hanler CacheHashValueUpdateHandler[ValueType]) error {
		value, err := getter()
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

func UseHashCacheIndex[ValueType any](cacher Cacher, key string, index string) (CacheValueGetter[ValueType], CacheValueUpdater[ValueType]) {

	getter := func() (*ValueType, error) {
		value, err := cacher.GetH(key, index)
		if err != nil {
			return nil, err
		}
		rst := new(ValueType)
		err = value.BindTo(rst)
		return rst, err
	}

	updater := func(hanler CacheValueUpdateHandler[ValueType]) error {
		value, err := getter()
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

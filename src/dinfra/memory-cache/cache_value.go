package memorycache

import "reflect"

type (
	MemoryCacheValue struct {
		value any
	}
)

func (value *MemoryCacheValue) BindTo(dst any) error {
	return nil
}

func (value *MemoryCacheValue) To(dstType reflect.Type) (any, error) {
	return nil, nil
}

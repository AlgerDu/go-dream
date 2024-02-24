package memorycache

import "reflect"

type (
	MemoryCacheValue struct {
		value any
	}
)

func (value *MemoryCacheValue) To(dstType reflect.Type) (any, error) {
	return value.value, nil
}

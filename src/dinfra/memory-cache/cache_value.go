package memorycache

type (
	MemoryCacheValue struct {
		value any
	}
)

func (value *MemoryCacheValue) BindTo(dst any) error {
	return nil
}

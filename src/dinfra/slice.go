package dinfra

func Slice_FindIndex[DataType any](s []DataType, check func(DataType) bool) int {
	for i, v := range s {
		if check(v) {
			return i
		}
	}
	return -1
}

func Slice_AddOrUpdate[DataType any](s []DataType, data DataType, check func(DataType) bool) []DataType {
	findIndex := Slice_FindIndex(s, check)
	if findIndex > -1 {
		s[findIndex] = data
		return s
	}
	return append(s, data)
}

func Slice_Delete[DataType any](s []DataType, check func(DataType) bool) []DataType {
	tmp := []DataType{}
	for _, v := range s {
		if !check(v) {
			tmp = append(tmp, v)
		}
	}
	return tmp
}

func Slice_ConvertTo[FromType any, DstType any](
	froms []FromType,
	converter func(FromType) (DstType, error),
) ([]DstType, error) {
	rst := []DstType{}
	for _, from := range froms {
		dst, err := converter(from)
		if err != nil {
			return rst, err
		}
		rst = append(rst, dst)
	}
	return rst, nil
}

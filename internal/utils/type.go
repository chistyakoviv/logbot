package utils

import "reflect"

func IsEmpty(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)

	// typed nils: *T, []T, map[T]T, chan, func, interface
	switch rv.Kind() {
	case reflect.Pointer, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface:
		return rv.IsNil()
	}

	switch rv.Kind() {
	case reflect.String:
		return rv.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rv.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return rv.Float() == 0
	}

	return false
}

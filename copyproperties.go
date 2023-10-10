package clone

import "reflect"

/**
 *
 * @Author AiTao
 * @Date 2023/6/18 11:48
 * @Url
 **/

// CopyProperties only all exported fields in the struct type are copied.
func CopyProperties(src any) (dst any) {
	if src == nil {
		return nil
	}
	srcValue := reflect.ValueOf(src)
	if !srcValue.IsValid() {
		return nil
	}
	if !isStruct(srcValue) {
		panic("source data must be a struct type")
	}
	dstElem := rnew(srcValue).Elem()
	ctx.flags = OnlyPublicField
	deepCopy(srcValue, dstElem)
	return dstElem.Interface()
}

func isStruct(value reflect.Value) bool {
	for value.Kind() == reflect.Pointer {
		value = value.Elem()
	}
	switch value.Kind() {
	case reflect.Struct:
		return true
	case reflect.Interface:
		return value.Elem().Kind() == reflect.Struct
	}
	return false
}

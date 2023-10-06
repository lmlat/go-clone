package cloner

import (
	"reflect"
)

/**
 * @Author AiTao
 * @Date 2023/6/15 11:03
 **/

// Deep returns the copy value after a deep-cloned source object.
func Deep(src any) (dst any) {
	if src == nil {
		return nil
	}
	srcValue := reflect.ValueOf(src)
	if !srcValue.IsValid() {
		return nil
	}
	if isBasicType(srcValue.Kind()) {
		return src
	}
	dstElem := rnew(srcValue).Elem()
	deepcopy(srcValue, dstElem)
	return dstElem.Interface()
}

// Shallow returns the copy value after a shallow-cloned source object.
func Shallow(src any) (dst any) {
	if src == nil {
		return nil
	}
	srcValue := reflect.ValueOf(src)
	if !srcValue.IsValid() {
		return nil
	}
	if isBasicType(srcValue.Kind()) {
		return src
	}
	dstElem := rnew(srcValue).Elem()
	dstElem.Set(srcValue)
	return dstElem.Interface()
}

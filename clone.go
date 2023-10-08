package cloner

import (
	"reflect"
)

/**
 * @Author AiTao
 * @Date 2023/6/15 11:03
 **/

// Deep returns the copy value after a deep-cloned source object.
func Deep(src any, options ...Options) (dst any) {
	if src == nil {
		return nil
	}
	srcValue := reflect.ValueOf(src)
	if !srcValue.IsValid() {
		return nil
	}
	dstElem := rnew(srcValue).Elem()
	for _, option := range options {
		option(ctx)
	}
	deepCopy(srcValue, dstElem)
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
	dstValue := rnew(srcValue)
	shallowCopy(srcValue, dstValue.UnsafePointer())
	return dstValue.Elem().Interface()
}

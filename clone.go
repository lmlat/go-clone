package clone

import (
	"reflect"
	"time"
	"unsafe"
)

/**
 * @Author AiTao
 * @Date 2023/8/17 11:03
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

func deepcopy(src, dst reflect.Value) {
	if src.CanInterface() {
		if cloneable, ok := src.Interface().(Cloneable); ok {
			dst.Set(reflect.ValueOf(cloneable.DeepClone()))
			return
		}
	}
	switch src.Kind() {
	case reflect.Pointer:
		copyPointer(src, dst)
	case reflect.Interface:
		copyInterface(src, dst)
	case reflect.Struct:
		copyStruct(src, dst)
	case reflect.Slice:
		copySlice(src, dst)
	case reflect.Map:
		copyMap(src, dst)
	case reflect.Func:
		copyFunc(src, dst)
	default:
		dst.Set(src)
	}
}

// copyFunc copy the value of a function type.
func copyFunc(src, dst reflect.Value) {
	if src.IsNil() {
		return
	}
	dst.Set(reflect.MakeFunc(src.Type(), func(args []reflect.Value) []reflect.Value { return src.Call(args) }))
}

// copyMap copy the value of a map type.
func copyMap(src, dst reflect.Value) {
	if src.IsNil() {
		return
	}
	dst.Set(reflect.MakeMap(src.Type()))
	for itr := src.MapRange(); itr.Next(); {
		srcValue := itr.Value()
		dstValue := rnew(srcValue).Elem()
		deepcopy(srcValue, dstValue)
		dst.SetMapIndex(itr.Key(), dstValue)
	}
}

// copySlice copy the value of a slice type.
func copySlice(src, dst reflect.Value) {
	if src.IsNil() {
		return
	}
	srcLen := src.Len()
	dst.Set(reflect.MakeSlice(src.Type(), srcLen, src.Cap()))
	for i := 0; i < srcLen; i++ {
		deepcopy(src.Index(i), dst.Index(i))
	}
}

// copyStruct copy the value of a struct type, containing the fields in the struct that are not exported.
func copyStruct(src, dst reflect.Value) {
	srcValuePtr, dstValuePtr := toPointer(src), toPointer(dst)
	copyStructPointer(srcValuePtr, dstValuePtr)
	dst.Set(dstValuePtr.Elem())
}

// copyStructPointer copy the value of a struct pointer type, containing the fields in the struct that are not exported.
func copyStructPointer(src, dst reflect.Value) {
	srcElem, dstElem := src.Elem(), dst.Elem()
	if !srcElem.IsValid() || srcElem.IsZero() {
		return
	}
	if tim, ok := srcElem.Interface().(time.Time); ok {
		dstElem.Set(reflect.ValueOf(tim))
		return
	}
	t, length := srcElem.Type(), srcElem.NumField()
	for i := 0; i < length; i++ {
		srcField, dstField, field := srcElem.Field(i), dstElem.Field(i), t.Field(i)
		if field.IsExported() {
			deepcopy(srcField, dstField)
		} else {
			copyNonExportedFields(srcField, dstField)
		}
	}
}

// copyNonExportedFields copy the value of non-exported fields in the struct.
func copyNonExportedFields(src, dst reflect.Value) {
	if src.CanAddr() {
		// associate the addresses of the source and target fields.
		srcUnPtr, dstUnPtr := src.Addr().UnsafePointer(), dst.Addr().UnsafePointer()
		// returns the number of bytes consumed by the source field data.
		size := int(src.Type().Size())
		copy(toBytes(dstUnPtr, size), toBytes(srcUnPtr, size))
	}
}

// copyInterface copy the value of a interface type.
func copyInterface(src, dst reflect.Value) {
	if src.IsNil() {
		return
	}
	// When calling the Elem method to get the underlying value of the interface, it should be noted that the underlying
	// value of the interface is not addressable. That is, executing the srcElem.CanAddr() method returns false.
	srcElem := src.Elem()
	dstElem := rnew(srcElem).Elem()
	deepcopy(srcElem, dstElem)
	dst.Set(dstElem)
}

// copyPointer copy the value of a pointer type.
func copyPointer(src, dst reflect.Value) {
	srcElem := src.Elem()
	if !srcElem.IsValid() { // intercept nil pointers
		return
	}
	dst.Set(rnew(srcElem))
	dstElem := dst.Elem()
	deepcopy(srcElem, dstElem)
}

func isBasicType(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Bool,
		reflect.String:
		return true
	default:
		return false
	}
}

// rnew returns a reflect.Value instance.
func rnew(value reflect.Value) reflect.Value {
	return reflect.New(value.Type())
}

func toBytes(ptr unsafe.Pointer, size int) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(ptr), Len: size, Cap: size}))
}

func toPointer(value reflect.Value) reflect.Value {
	valuePtr := reflect.New(value.Type())
	srcPointerValue := valuePtr.Elem()
	srcPointerValue.Set(value)
	return valuePtr
}

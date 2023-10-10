package clone

import (
	"fmt"
	"reflect"
	"unsafe"
)

/**
 *
 * @Author AiTao
 * @Date 2023/10/8 5:47
 **/

// For detailed code ideas, please refer to https://github.com/huandu/go-clone/blob/master/clone.go#L372
func shallowCopy(src reflect.Value, ptr unsafe.Pointer) {
	switch src.Kind() {
	case reflect.Bool:
		*(*bool)(ptr) = src.Bool()
	case reflect.Int:
		*(*int)(ptr) = int(src.Int())
	case reflect.Int8:
		*(*int8)(ptr) = int8(src.Int())
	case reflect.Int16:
		*(*int16)(ptr) = int16(src.Int())
	case reflect.Int32:
		*(*int32)(ptr) = int32(src.Int())
	case reflect.Int64:
		*(*int64)(ptr) = src.Int()
	case reflect.Uint:
		*(*uint)(ptr) = uint(src.Uint())
	case reflect.Uint8:
		*(*uint8)(ptr) = uint8(src.Uint())
	case reflect.Uint16:
		*(*uint16)(ptr) = uint16(src.Uint())
	case reflect.Uint32:
		*(*uint32)(ptr) = uint32(src.Uint())
	case reflect.Uint64:
		*(*uint64)(ptr) = src.Uint()
	case reflect.Uintptr:
		*(*uintptr)(ptr) = uintptr(src.Uint())
	case reflect.Float32:
		*(*float32)(ptr) = float32(src.Float())
	case reflect.Float64:
		*(*float64)(ptr) = src.Float()
	case reflect.Complex64:
		*(*complex64)(ptr) = complex64(src.Complex())
	case reflect.Complex128:
		*(*complex128)(ptr) = src.Complex()
	case reflect.Array:
		t := src.Type()
		if src.CanAddr() {
			srcPtr := src.Addr().UnsafePointer()
			size := t.Size()
			copy((*[1 << 30]byte)(ptr)[:size:size], (*[1 << 30]byte)(srcPtr)[:size:size])
		} else {
			// source data is not addressable.
			dst := reflect.NewAt(t, ptr).Elem()
			if src.CanInterface() {
				dst.Set(src)
			} else {
				size, length := t.Elem().Size(), src.Len()
				for i := 0; i < length; i++ {
					elemPtr := unsafe.Pointer(uintptr(ptr) + uintptr(i)*size)
					shallowCopy(src.Index(i), elemPtr)
				}
			}
		}
	case reflect.Chan, reflect.Map, reflect.Pointer, reflect.UnsafePointer:
		*((*uintptr)(ptr)) = src.Pointer()
	case reflect.Func:
		reflect.NewAt(src.Type(), ptr).Elem().Set(copyUnexportedValue(src))
	case reflect.Interface:
		*((*ifaceSource)(ptr)) = parseReflectValue(src)
	case reflect.Slice:
		*(*reflect.SliceHeader)(ptr) = reflect.SliceHeader{Data: src.Pointer(), Len: src.Len(), Cap: src.Cap()}
	case reflect.String:
		// String content is read-only in memory.
		// When copying a string, the string content is copied to a new object with a new address,
		// but the content is the same as the original string.
		reflect.NewAt(src.Type(), ptr).Elem().SetString(src.String())
	case reflect.Struct:
		t := src.Type()
		dst := reflect.NewAt(t, ptr).Elem()
		if src.CanInterface() {
			dst.Set(src)
		} else {
			length := t.NumField()
			for i := 0; i < length; i++ {
				shallowCopy(src.Field(i), unsafe.Pointer(uintptr(ptr)+t.Field(i).Offset))
			}
		}
	default:
		panic(fmt.Errorf("failed to cloning non-exported `%v` type field", src.Type()))
	}
}

// For detailed code ideas, please refer to https://github.com/huandu/go-clone/blob/master/structtype.go#L209
func copyUnexportedValue(src reflect.Value) (dst reflect.Value) {
	if src.CanInterface() {
		return src
	}
	// If the source data is an unexported field value, make a copy of its value.
	switch src.Kind() {
	case reflect.Bool:
		return reflect.ValueOf(src.Bool())
	case reflect.Int:
		return reflect.ValueOf(int(src.Int()))
	case reflect.Int8:
		return reflect.ValueOf(int8(src.Int()))
	case reflect.Int16:
		return reflect.ValueOf(int16(src.Int()))
	case reflect.Int32:
		return reflect.ValueOf(int32(src.Int()))
	case reflect.Int64:
		return reflect.ValueOf(src.Int())
	case reflect.Uint:
		return reflect.ValueOf(uint(src.Uint()))
	case reflect.Uint8:
		return reflect.ValueOf(uint8(src.Uint()))
	case reflect.Uint16:
		return reflect.ValueOf(uint16(src.Uint()))
	case reflect.Uint32:
		return reflect.ValueOf(uint32(src.Uint()))
	case reflect.Uint64:
		return reflect.ValueOf(src.Uint())
	case reflect.Uintptr:
		return reflect.ValueOf(uintptr(src.Uint()))
	case reflect.Float32:
		return reflect.ValueOf(float32(src.Float()))
	case reflect.Float64:
		return reflect.ValueOf(src.Float())
	case reflect.Complex64:
		return reflect.ValueOf(complex64(src.Complex()))
	case reflect.Complex128:
		return reflect.ValueOf(src.Complex())
	case reflect.String:
		return reflect.ValueOf(src.String())
	case reflect.Func:
		return getFuncReflectValue(src)
	case reflect.UnsafePointer:
		return reflect.ValueOf(unsafe.Pointer(src.Pointer()))
	}
	panic(fmt.Errorf("failed to cloning non-exported `%v` type field", src.Type()))
}

func getFuncReflectValue(src reflect.Value) reflect.Value {
	if src.IsNil() {
		return reflect.Zero(src.Type())
	}
	// It should be noted that the function type is non-addressable.
	// that is, calling the CanAddr() will return false
	src = src.Convert(typeOfAny) // convert to `any` type.
	var a any
	av := reflect.ValueOf(&a) // *any
	*(*ifaceSource)(av.UnsafePointer()) = parseReflectValue(src)
	// return the actual value pointed to by the pointer.
	avElem := av.Elem() // any
	// returns the underlying value of the interface.
	return avElem.Elem()
}

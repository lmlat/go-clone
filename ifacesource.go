package cloner

import (
	"reflect"
	"unsafe"
)

/**
 * @Author AiTao
 * @Date 2023/10/8 8:01
 *
 **/

// For detailed code ideas, please refer to https://github.com/huandu/go-clone/blob/master/interfacedata.go#L18
type ifaceSource struct {
	// the definition of the `_` field does not allocate any meaningful data, it is only used to take up memory space.
	_ [sizeOfPointers]unsafe.Pointer
}

// Returns the number of bytes occupied by the pointer.
//
// On 64-bit systems, the size of `uintptr` is 8 bytes, while on 32-bit systems, the size of `uintptr` is usually 4 bytes.
//
// On 64-bit systems, the size of `any` is 16 bytes, it is worth noting that the size of the `any` interface may vary
// under different systems and compilers.
const sizeOfPointers = unsafe.Sizeof((any)(0)) / unsafe.Sizeof(uintptr(0))

var reflectValuePtrOffset uintptr

func init() {
	reflectValuePtrOffset = getPtrFieldOffsetInReflectValue(reflect.Value{})
}

func getPtrFieldOffsetInReflectValue(src reflect.Value) uintptr {
	t := reflect.TypeOf(src)
	length := t.NumField()
	// traverse all fields in `reflect.Value`.
	for i := 0; i < length; i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.UnsafePointer {
			return field.Offset
		}
	}
	panic("failed to find internal ptr field in `reflect.Value`")
}

// For detailed code ideas, please refer to https://github.com/huandu/go-clone/blob/master/interfacedata.go#L46
func parseReflectValue(v reflect.Value) ifaceSource {
	pv := (unsafe.Pointer)(uintptr(unsafe.Pointer(&v)) + reflectValuePtrOffset)
	ptr := *(*unsafe.Pointer)(pv)
	return *(*ifaceSource)(ptr)
}

package cloner

import (
	"reflect"
	"sync"
	"time"
	"unsafe"
)

/**
 * @Author AiTao
 * @Date 2023/6/15 11:48
 **/

const Ignore = "ignore"

type cacheKey struct {
	ptr  uintptr
	typ  reflect.Type
	size int
}

type cache map[cacheKey]reflect.Value

type Empty struct{}

var (
	ctx             *Context
	typeOfString    = reflect.TypeOf("")
	typeOfAny       = reflect.TypeOf((*any)(nil)).Elem()
	typeOfByteSlice = reflect.TypeOf([]byte(nil))
)

type Context struct {
	enableCache      bool
	cache            cache
	shallowCopyTypes sync.Map
	flags            OpFlags
}

type Options func(*Context)

func WithEnableCache(enableCache bool) Options {
	return func(ctx *Context) {
		ctx.enableCache = enableCache
	}
}

func WithTypes(types ...reflect.Type) Options {
	return func(ctx *Context) {
		for _, t := range types {
			ctx.shallowCopyTypes.Store(t, Empty{})
		}
	}
}
func WithOpFlags(flags OpFlags) Options {
	return func(ctx *Context) {
		ctx.flags = flags
	}
}

func init() {
	ctx = &Context{
		enableCache:      true,
		cache:            cache{},
		shallowCopyTypes: sync.Map{},
		flags:            AllFields | DeepString | DeepArray,
	}

	// reflect.rtype is immutable because the Go language's type system is static at compile time,
	// and the type information is determined at compile time.
	ctx.shallowCopyTypes.Store(reflect.TypeOf(reflect.TypeOf(Empty{})), Empty{})
	ctx.shallowCopyTypes.Store(reflect.TypeOf(time.Time{}), Empty{})
}

func deepCopy(src, dst reflect.Value) {
	kind := src.Kind()
	if isShallowCopyType(kind) {
		dst.Set(copyUnexportedValue(src))
		return
	}
	switch kind {
	case reflect.Struct:
		if cloneable, ok := src.Interface().(Cloneable); ok {
			dst.Set(reflect.ValueOf(cloneable.DeepClone()))
			return
		}
		copyStruct(src, dst)
	case reflect.Array:
		copyArray(src, dst)
	case reflect.Slice:
		copySlice(src, dst)
	case reflect.String:
		copyString(src, dst)
	case reflect.Map:
		copyMap(src, dst)
	case reflect.Interface:
		copyInterface(src, dst)
	case reflect.Pointer:
		copyPointer(src, dst)
	case reflect.Chan:
		copyChan(src, dst)
	case reflect.Func:
		copyFunc(src, dst)
	default:
		dst.Set(copyUnexportedValue(src))
	}
}

// copyString copy the value of a string type.
func copyString(src, dst reflect.Value) {
	if src.IsZero() {
		return
	}
	t, size := src.Type(), src.Len()
	bytes := reflect.MakeSlice(typeOfByteSlice, size, size)
	if !src.CanInterface() {
		src = reflect.ValueOf(src.String())
	}
	reflect.Copy(bytes, src)
	dstPtr := reflect.New(t)
	slice := bytes.Interface().([]byte)
	*(*reflect.StringHeader)(dstPtr.UnsafePointer()) = *(*reflect.StringHeader)(unsafe.Pointer(&slice))
	dst.Set(dstPtr.Elem())
}

// copyChan copy the value of a chan type.
func copyChan(src reflect.Value, dst reflect.Value) {
	if src.IsNil() {
		return
	}
	dst.Set(reflect.MakeChan(src.Type(), src.Cap()))
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
	t := src.Type()
	if ctx.enableCache && ctx.cache != nil {
		k := cacheKey{ptr: src.Pointer(), typ: t}
		if val, ok := ctx.cache[k]; ok {
			dst.Set(val)
			return
		}
	}
	m := reflect.MakeMapWithSize(t, src.Len())
	if ctx.enableCache && ctx.cache != nil {
		k := cacheKey{ptr: src.Pointer(), typ: t}
		ctx.cache[k] = m // mark map accessed.
	}
	dst.Set(m)
	for itr := src.MapRange(); itr.Next(); {
		srcValue := itr.Value()
		dstValue := rnew(srcValue).Elem()
		deepCopy(srcValue, dstValue)

		srcKey := itr.Key()
		dstKey := rnew(srcKey).Elem()
		deepCopy(srcKey, dstKey)

		dst.SetMapIndex(srcKey, dstValue)
	}
}

// copySlice copy the value of a slice type.
func copySlice(src, dst reflect.Value) {
	if src.IsNil() {
		return
	}
	t := src.Type()
	srcLen := src.Len()
	if ctx.enableCache && ctx.cache != nil {
		k := cacheKey{src.Pointer(), t, srcLen}
		if val, ok := ctx.cache[k]; ok {
			dst.Set(val)
			return
		}
	}
	s := reflect.MakeSlice(src.Type(), srcLen, src.Cap())
	if ctx.enableCache && ctx.cache != nil {
		k := cacheKey{src.Pointer(), t, srcLen}
		ctx.cache[k] = s // mark slice accessed.
	}
	dst.Set(s)
	for i := 0; i < srcLen; i++ {
		deepCopy(src.Index(i), dst.Index(i))
	}
}

func copyArray(src, dst reflect.Value) {
	if src.IsZero() {
		return
	}
	size := src.Len()
	elementKind := src.Type().Elem().Kind()
	if isShallowCopyType(elementKind) {
		dst.Set(copyUnexportedValue(src))
		return
	}
	for i := 0; i < size; i++ {
		deepCopy(src.Index(i), dst.Index(i))
	}
}

// copyStruct copy the value of a struct type, containing the fields in the struct that are not exported.
func copyStruct(src, dst reflect.Value) {
	srcValuePtr, dstValuePtr := asPointer(src), asPointer(dst)
	copyStructPointer(srcValuePtr, dstValuePtr)
	dst.Set(dstValuePtr.Elem())
}

// copyStructPointer copy the value of a struct pointer type, containing the fields in the struct that are not exported.
func copyStructPointer(src, dst reflect.Value) {
	srcElem, dstElem := src.Elem(), dst.Elem()
	if !srcElem.IsValid() || srcElem.IsZero() {
		return
	}
	t := srcElem.Type()
	if _, ok := ctx.shallowCopyTypes.Load(t); ok {
		if srcElem.CanInterface() {
			dstElem.Set(srcElem)
		} else {
			srcPtr := reflect.New(t)
			shallowCopy(srcElem, srcPtr.UnsafePointer())
			dstElem.Set(srcPtr.Elem())
		}
		return
	}
	length := srcElem.NumField()
	for i := 0; i < length; i++ {
		srcField, dstField, field := srcElem.Field(i), dstElem.Field(i), t.Field(i)
		// skip all fields with ignore tag in the struct.
		if _, ok := field.Tag.Lookup(Ignore); ok {
			continue
		}
		if field.IsExported() { // copy the exported fields.
			if ctx.flags == OnlyPrivateField {
				continue
			}
			deepCopy(srcField, dstField)
		} else { // copy the non-exported fields.
			if ctx.flags == OnlyPublicField {
				continue
			}
			copyUnexportedFields(srcField, dstField)
		}
	}
}

// copyUnexportedFields copy the value of non-exported fields in the struct.
func copyUnexportedFields(src, dst reflect.Value) {
	// field type must be addressable.
	if src.CanAddr() {
		t := src.Type()
		// fixed memory sharing issue where copying did not export reference type fields.
		srcValue := reflect.NewAt(t, src.Addr().UnsafePointer())
		srcElem := srcValue.Elem()
		dstElem := rnew(srcElem).Elem()
		// deep copy reference type fields that are not exported.
		if isDeepCopyType(srcElem.Kind()) {
			deepCopy(srcElem, dstElem)
		} else {
			dstElem.Set(copyUnexportedValue(srcElem))
		}
		// memory copy.
		memcopy(dstElem.Addr().UnsafePointer(), dst.Addr().UnsafePointer(), int(t.Size()))
	}
}

// copyInterface copy the value of a interface type.
func copyInterface(src, dst reflect.Value) {
	// nil interface are not handled.
	if src.IsNil() {
		return
	}
	// when calling the Elem method to get the underlying value of the interface, it should be noted that the underlying
	// value of the interface is not addressable.
	// that is, executing the srcElem.CanAddr() method returns false.
	srcElem := src.Elem()
	if !srcElem.IsValid() {
		return
	}
	dstElem := rnew(srcElem).Elem()
	deepCopy(srcElem, dstElem)
	dst.Set(dstElem)
}

// copyPointer copy the value of a pointer type.
func copyPointer(src, dst reflect.Value) {
	if src.IsNil() {
		return
	}
	t := src.Type()
	// returns the actual value pointed to by the pointer.
	srcElem := src.Elem()
	// intercept nil pointers.
	if !srcElem.IsValid() {
		return
	}
	if _, ok := ctx.shallowCopyTypes.Load(t); ok {
		if src.CanInterface() {
			dst.Set(src)
		} else {
			srcPtr := reflect.New(t)
			shallowCopy(src, srcPtr.UnsafePointer())
			dst.Set(srcPtr.Elem())
		}
		return
	}
	// If the current pointer has been deeply cloned, get the reflect.Value instance directly from the cache.
	if ctx.enableCache && ctx.cache != nil {
		k := cacheKey{ptr: src.Pointer(), typ: t}
		if val, ok := ctx.cache[k]; ok {
			dst.Set(val)
			return
		}
	}
	p := rnew(srcElem)
	if ctx.enableCache && ctx.cache != nil {
		k := cacheKey{ptr: src.Pointer(), typ: t}
		ctx.cache[k] = p // make pointer accessed.
	}
	dst.Set(p)
	dstElem := dst.Elem()
	deepCopy(srcElem, dstElem)
}

func isDeepCopyType(kind reflect.Kind) bool {
	switch kind {
	case reflect.Slice, reflect.Map, reflect.Interface, reflect.Pointer, reflect.Chan:
		return true
	case reflect.Func:
		return ctx.flags.Has(DeepFunc)
	case reflect.String:
		return ctx.flags.Has(DeepString)
	case reflect.Array:
		return ctx.flags.Has(DeepArray)
	default:
		return false
	}
}

func isShallowCopyType(kind reflect.Kind) bool {
	switch kind {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.UnsafePointer:
		return true
	case reflect.Func:
		return !ctx.flags.Has(DeepFunc)
	case reflect.String:
		return !ctx.flags.Has(DeepString)
	case reflect.Array:
		return !ctx.flags.Has(DeepArray)
	}
	return false
}

// rnew returns a reflect.Value instance.
func rnew(value reflect.Value) reflect.Value {
	return reflect.New(value.Type())
}

// memcopy memory copy.
// src  - source pointer
// dst  - target pointer
// size - the number of bytes consumed by the source field data.
func memcopy(src, dst unsafe.Pointer, size int) {
	copy(asBytes(dst, size), asBytes(src, size))
}

func asBytes(ptr unsafe.Pointer, size int) (bytes []byte) {
	slice := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	slice.Data = uintptr(ptr)
	slice.Len = size
	slice.Cap = size
	return bytes
}

func asPointer(value reflect.Value) reflect.Value {
	valuePtr := rnew(value)
	srcPtrValue := valuePtr.Elem()
	srcPtrValue.Set(value)
	return valuePtr
}

package cloner

import (
	"bytes"
	"fmt"
	"github.com/huandu/go-assert"
	"io"
	"reflect"
	"testing"
	"time"
	"unsafe"
)

/**
 *
 * @Author AiTao
 * @Date 2023/6/17 11:13
 * @Url
 **/

type PrivateStruct struct {
	string      string
	strings     []string
	stringArr   [4]string
	bool        bool
	bools       []bool
	byte        byte
	bytes       []byte
	int         int
	ints        []int
	int8        int8
	int8s       []int8
	int16       int16
	int16s      []int16
	int32       int32
	int32s      []int32
	int64       int64
	int64s      []int64
	uint        uint
	uints       []uint
	uint8       uint8
	uint8s      []uint8
	uint16      uint16
	uint16s     []uint16
	uint32      uint32
	uint32s     []uint32
	uint64      uint64
	uint64s     []uint64
	float32     float32
	float32s    []float32
	float64     float64
	float64s    []float64
	complex64   complex64
	complex64s  []complex64
	complex128  complex128
	complex128s []complex128
	iface       any
	ifaces      []any
}

type PublicStruct struct {
	String      string
	Strings     []string
	StringArr   [4]string
	Bool        bool
	Bools       []bool
	Byte        byte
	Bytes       []byte
	Int         int
	Ints        []int
	Int8        int8
	Int8s       []int8
	Int16       int16
	Int16s      []int16
	Int32       int32
	Int32s      []int32
	Int64       int64
	Int64s      []int64
	Uint        uint
	Uints       []uint
	Uint8       uint8
	Uint8s      []uint8
	Uint16      uint16
	Uint16s     []uint16
	Uint32      uint32
	Uint32s     []uint32
	Uint64      uint64
	Uint64s     []uint64
	Float32     float32
	Float32s    []float32
	Float64     float64
	Float64s    []float64
	Complex64   complex64
	Complex64s  []complex64
	Complex128  complex128
	Complex128s []complex128
	Interface   any
	Interfaces  []any
}

func TestClone_PrivateStruct_ZeroValue(t *testing.T) {
	src := PrivateStruct{}
	dst := Deep(src).(PrivateStruct)
	fmt.Printf("%+v\n", src)
	fmt.Printf("\n%+v\n", dst)
}

func TestClone_PublicStruct_ZeroValue(t *testing.T) {
	src := PublicStruct{}
	dst := Deep(src).(PublicStruct)
	fmt.Printf("%+v\n", src)
	fmt.Printf("\n%+v\n", dst)
}

func TestClone_PrivateStruct(t *testing.T) {
	src := PrivateStruct{
		"kimchi",
		[]string{"uni", "ika"},
		[4]string{"malort", "barenjager", "fernet", "salmiakki"},
		true,
		[]bool{true, false, true},
		'z',
		[]byte("abc"),
		42,
		[]int{0, 1, 3, 4},
		8,
		[]int8{8, 9, 10},
		16,
		[]int16{16, 17, 18, 19},
		32,
		[]int32{32, 33},
		64,
		[]int64{64},
		420,
		[]uint{11, 12, 13},
		81,
		[]uint8{81, 82},
		160,
		[]uint16{160, 161, 162, 163, 164},
		320,
		[]uint32{320, 321},
		640,
		[]uint64{6400, 6401, 6402, 6403},
		32.32,
		[]float32{32.32, 33},
		64.1,
		[]float64{64, 65, 66},
		complex64(-64 + 12i),
		[]complex64{complex64(-65 + 11i), complex64(66 + 10i)},
		complex128(-128 + 12i),
		[]complex128{complex128(-128 + 11i), complex128(129 + 10i)},
		nil,
		[]any{42, true, "pan-galactic"},
	}

	dst := Deep(src).(PrivateStruct)
	fmt.Printf("%+v\n", src)
	fmt.Printf("\n%+v\n", dst)
}

func TestClone_PublicStruct(t *testing.T) {
	src := PublicStruct{
		String:      "kimchi",
		Strings:     []string{"uni", "ika"},
		StringArr:   [4]string{"malort", "barenjager", "fernet", "salmiakki"},
		Bool:        true,
		Bools:       []bool{true, false, true},
		Byte:        'z',
		Bytes:       []byte("abc"),
		Int:         42,
		Ints:        []int{0, 1, 3, 4},
		Int8:        8,
		Int8s:       []int8{8, 9, 10},
		Int16:       16,
		Int16s:      []int16{16, 17, 18, 19},
		Int32:       32,
		Int32s:      []int32{32, 33},
		Int64:       64,
		Int64s:      []int64{64},
		Uint:        420,
		Uints:       []uint{11, 12, 13},
		Uint8:       81,
		Uint8s:      []uint8{81, 82},
		Uint16:      160,
		Uint16s:     []uint16{160, 161, 162, 163, 164},
		Uint32:      320,
		Uint32s:     []uint32{320, 321},
		Uint64:      640,
		Uint64s:     []uint64{6400, 6401, 6402, 6403},
		Float32:     32.32,
		Float32s:    []float32{32.32, 33},
		Float64:     64.1,
		Float64s:    []float64{64, 65, 66},
		Complex64:   complex64(-64 + 12i),
		Complex64s:  []complex64{complex64(-65 + 11i), complex64(66 + 10i)},
		Complex128:  complex128(-128 + 12i),
		Complex128s: []complex128{complex128(-128 + 11i), complex128(129 + 10i)},
		Interfaces:  []any{42, true, "pan-galactic"},
	}
	dst := Deep(src).(PublicStruct)
	fmt.Printf("%+v\n", src)
	fmt.Printf("\n%+v\n", dst)

	fmt.Printf("%+v\n", src)
	fmt.Printf("\n%+v\n", dst)
}

func TestClone_Func(t *testing.T) {
	otherFunc := func(a, b int) int { return a - b }
	srcFunc := func(a, b int) int { return a + b }
	dstFunc := Deep(srcFunc).(func(a, b int) int)
	fmt.Printf("%+v, %p\n", srcFunc(1, 2), srcFunc)
	fmt.Printf("%+v, %p\n", dstFunc(1, 2), dstFunc)

	dstFunc = otherFunc
	fmt.Printf("%+v, %p\n", srcFunc(1, 2), srcFunc)
	fmt.Printf("%+v, %p\n", dstFunc(1, 2), dstFunc)
}

func TestClone_Slice(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}
	dst := Deep(src).([]int)
	fmt.Printf("%+v, %p\n", src, src)
	fmt.Printf("\n%+v, %p\n", dst, dst)

	dst[3] = 100
	fmt.Printf("%+v, %p\n", src, src)
	fmt.Printf("\n%+v, %p\n", dst, dst)
}

func TestClone_Map(t *testing.T) {
	m := map[string]*Test{
		"a": {
			F: 123,
			S: map[string]any{
				"cba": 321,
			},
		},
		"b": {
			F: 456,
			S: map[string]any{
				"ghi": 789,
			},
		},
	}
	dst := Deep(m).(map[string]*Test)

	fmt.Printf("%+v, %p, %v, %v\n", m, m, m["a"].S, m["b"].F)
	fmt.Printf("%+v, %p, %v, %v\n", dst, dst, dst["a"].S, dst["b"].F)

	dst["a"].S["cba"] = 999
	dst["b"].F = 888
	fmt.Printf("%+v, %p, %v, %v\n", m, m, m["a"].S, m["b"].F)
	fmt.Printf("%+v, %p, %v, %v\n", dst, dst, dst["a"].S, dst["b"].F)
}

func TestClone_Nil_Pointer(t *testing.T) {
	var uintPtr *uint
	dstUintPtr := Deep(uintPtr).(*uint)

	fmt.Printf("%+v, %p\n", uintPtr, uintPtr)
	fmt.Printf("%+v, %p\n", dstUintPtr, dstUintPtr)
}

func TestClone_Pointer(t *testing.T) {
	fmt.Println("====================克隆string指针====================")
	var str string
	strPtr := &str
	dstStrPtr := Deep(strPtr).(*string)
	fmt.Printf("%+v, %v\n", str, *strPtr)
	fmt.Printf("%+v, %v\n", str, *dstStrPtr)

	*dstStrPtr = "aitao"
	fmt.Printf("%+v, %v\n", str, *strPtr)
	fmt.Printf("%+v, %v\n", str, *dstStrPtr)

	fmt.Println("====================克隆int指针====================")
	ints := 100
	intPtr := &ints
	dstIntPtr := Deep(intPtr).(*int)
	fmt.Printf("%+v, %v\n", ints, *intPtr)
	fmt.Printf("%+v, %v\n", ints, *dstIntPtr)

	*dstIntPtr = 999
	fmt.Printf("%+v, %v\n", ints, *intPtr)
	fmt.Printf("%+v, %v\n", ints, *dstIntPtr)
}

type C struct {
	name     string
	Age      int
	birthday time.Time
	hobby    []string
}

func TestClone_StructC(t *testing.T) {
	src := C{"aitao", 100, time.Now(), []string{"ping pong", "badminton", "football"}}
	dst := Deep(src).(C)
	fmt.Printf("%+v\n", src)
	fmt.Printf("\n%+v\n", dst)

	dst.name = "小阿梦" // modify an unreferenced type
	fmt.Printf("%+v\n", src)
	fmt.Printf("\n%+v\n", dst)

	dst.hobby[0] = "乒乓球" // modify a referenced type
	fmt.Printf("%+v\n", src)
	fmt.Printf("\n%+v\n", dst)
}

func TestClone_Struct_Pointer(t *testing.T) {
	src := &C{"aitao", 100, time.Now(), []string{"ping pong", "badminton", "football"}}
	dst := Deep(src).(*C)
	fmt.Printf("%+v\n", src)
	fmt.Printf("%+v\n\n", dst)

	dst.name = "kqai"
	fmt.Printf("%+v\n", src)
	fmt.Printf("%+v\n", dst)
}

func TestClone_PrivateStruct_Pointer_ZeroValue(t *testing.T) {
	src := &PrivateStruct{}
	dst := Deep(src).(*PrivateStruct)
	fmt.Printf("%+v\n", src)
	fmt.Printf("%+v\n\n", dst)

	dst.string = "aitao"
	dst.strings = append(dst.strings, "aitao")
	fmt.Printf("%+v\n", src)
	fmt.Printf("%+v\n", dst)
}

func TestClone_PublicStruct_Pointer_ZeroValue(t *testing.T) {
	src := &PublicStruct{}
	dst := Deep(src).(*PublicStruct)
	fmt.Printf("%+v\n", src)
	fmt.Printf("%+v\n\n", dst)

	dst.String = "aitao"
	dst.Strings = append(dst.Strings, "aitao")
	fmt.Printf("%+v\n", src)
	fmt.Printf("%+v\n", dst)
}

func TestClone_Struct_Time(t *testing.T) {
	src, _ := time.ParseInLocation(time.DateTime, "2023-04-22 10:30:00", time.Local)
	dst := Deep(src).(time.Time)
	fmt.Printf("%v,%p\n", src, &src)
	fmt.Printf("%v,%p\n\n", dst, &dst)

	dst = time.Now()
	fmt.Printf("%+v,%p\n", src, &src)
	fmt.Printf("%+v,%p\n", dst, &dst)
}

type A struct {
	Int     int
	String  string
	uints   []uint
	strings []string
	Map     map[string]int
	MapB    map[string]*B
	bs      []B
	B
	time time.Time
}

type B struct {
	Vals []string
}

func TestClone_StructA(t *testing.T) {
	src := A{
		Int:    42,
		String: "Konichiwa",
		uints:  []uint{0, 1, 2, 3},
		Map:    map[string]int{"a": 1, "b": 2},
		MapB: map[string]*B{
			"hi":  &B{Vals: []string{"hello", "bonjour"}},
			"bye": &B{Vals: []string{"good-bye", "au revoir"}},
		},
		bs: []B{
			{Vals: []string{"Ciao", "Aloha"}},
		},
		B:    B{Vals: []string{"42"}},
		time: time.Now(),
	}
	dst := Deep(src).(A)
	fmt.Printf("%+v %p\n", src, &src)
	fmt.Printf("%+v %p\n\n", dst, &dst)
	dst.uints[0] = 1000
	fmt.Printf("%+v %p\n", src, &src)
	fmt.Printf("%+v %p\n", dst, &dst)
}

func TestDeep_BasicType(t *testing.T) {
	srcInt := 42
	dstInt := Deep(srcInt).(int)
	fmt.Printf("%+v %p\n", srcInt, &srcInt)
	fmt.Printf("\n%+v %p\n", dstInt, &dstInt)

	srcString := "Hello, World!"
	dstString := Deep(srcString).(string)
	fmt.Printf("%+v %p\n", srcString, &srcString)
	fmt.Printf("\n%+v %p\n", dstString, &dstString)

	srcSlice := []int{1, 2, 3, 4, 5}
	dstSlice := Deep(srcSlice).([]int)
	fmt.Printf("%+v %p\n", srcSlice, &srcSlice)
	fmt.Printf("\n%+v %p\n", dstSlice, &dstSlice)

	type Person struct {
		name string
		age  int
	}
	srcPerson := Person{"Alice", 30}
	dstPerson := Deep(srcPerson).(Person)
	fmt.Printf("%+v %p\n", srcPerson, &srcPerson)
	fmt.Printf("\n%+v %p\n", dstPerson, &dstPerson)
}

type D struct {
	name     string `ignore:"name"` // 忽略name字段
	Age      int
	hobby    []string
	birthday time.Time
}

func TestIgnoreTag(t *testing.T) {
	// name与birthday字段将不会被拷贝
	src := &D{"aitao", 100, []string{"pingpong", "badminton"}, time.Now()}
	dst := Deep(src).(*D)
	fmt.Printf("%+v %p\n", src, &src)
	fmt.Printf("%+v %p\n\n", dst, &dst)

	dst.hobby[0] = "乒乓球"
	fmt.Printf("%+v %p\n", src, &src)
	fmt.Printf("%+v %p\n", dst, &dst)
}

type F struct {
	strings []string
	m       map[string]int
	fn      func(int, int) int
	f       *F
	val     any
}

func TestMemorySharing(t *testing.T) {
	f := &F{[]string{"football"}, map[string]int{"cccc": 1, "dddd": 2}, func(a int, b int) int {
		return a - b
	}, &F{val: 888}, 999}

	src := &F{[]string{"pingpong", "badminton"}, map[string]int{"a": 1, "b": 2}, func(a int, b int) int {
		return a + b
	}, f, f}

	mul := func(a int, b int) int {
		return a * b
	}

	dst := Deep(src).(*F)
	fmt.Printf("%+v %p\n", src, &src)
	fmt.Printf("%+v %p\n\n", dst, &dst)

	fmt.Printf("%+v %p\n", dst.f, &dst.f)
	fmt.Printf("%+v %p\n\n", dst.f.f, &dst.f.f)
	dst.strings[0] = "乒乓球"
	dst.m["b"] = 99999

	dst.f.fn = mul
	fmt.Println(src.f.fn(20, 50), dst.f.fn(20, 50)) // -30 1000

	dst.f.f.strings = append(dst.f.f.strings, "aitao")
	dst.f.f.val = 999999999999999

	fmt.Printf("%+v %p\n", dst.f, &dst.f)
	fmt.Printf("%+v %p\n\n", dst.f.f, &dst.f.f)

	fmt.Printf("%+v %p\n", src, &src)
	fmt.Printf("%+v %p\n", dst, &dst)
}

type G struct {
	name  string
	Age   int
	hobby []string
	Sex   bool
}

func TestCopyProperties(t *testing.T) {
	// name与birthday字段将不会被拷贝
	var src any = &G{"aitao", 100, []string{"pingpong", "badminton"}, true}
	dst := CopyProperties(src).(*G)
	fmt.Printf("%+v %p\n", src, &src)
	fmt.Printf("%+v %p\n\n", dst, &dst)

	dst.name = "kqai"
	fmt.Printf("%+v %p\n", src, &src)
	fmt.Printf("%+v %p\n", dst, &dst)
}

func TestDeepCopy_Array(t *testing.T) {
	src := [2]*Test{
		{
			F: 123,
			S: map[string]any{
				"abc": 123,
			},
		},
		{
			F: 456,
			S: map[string]any{
				"def": 456,
				"ghi": 789,
			},
		},
	}

	dst := Shallow(src).([2]*Test)

	fmt.Printf("%+v,%p,%+v,%p,%+v,%p\n", src, &src, src[0], &src[0], src[1], &src[1])
	fmt.Printf("%+v,%p,%+v,%p,%+v,%p\n", dst, &dst, dst[0], &dst[0], dst[1], &dst[1])

	dst[0].F = 987
	dst[1].S["ghi"] = 321
	fmt.Printf("%+v,%p,%+v,%p,%+v,%p\n", src, &src, src[0], &src[0], src[1], &src[1])
	fmt.Printf("%+v,%p,%+v,%p,%+v,%p\n", dst, &dst, dst[0], &dst[0], dst[1], &dst[1])
}

func TestCloneReflectType(t *testing.T) {
	foo := reflect.TypeOf("aitao")
	dst := Deep(foo).(reflect.Type)

	from := reflect.ValueOf(foo)
	to := reflect.ValueOf(dst)

	fmt.Println(from.Pointer(), to.Pointer())
}

type scalarWriter int8

func (scalarWriter) Write(p []byte) (n int, err error) { return }

type Unexported struct {
	insider
}
type Simple struct {
	Foo int
	Bar string
}
type insider struct {
	i             int
	i8            int8
	i16           int16
	i32           int32
	i64           int64
	u             uint
	u8            uint8
	u16           uint16
	u32           uint32
	u64           uint64
	uptr          uintptr
	b             bool
	s             string
	f32           float32
	f64           float64
	c64           complex64
	c128          complex128
	arr           [4]string
	m             map[string]any
	ptr           *Unexported
	_             *Unexported
	slice         []*Unexported
	st            Simple
	unsafePointer unsafe.Pointer
	t             reflect.Type
	ch            chan bool
	fn            func(s string) string
	method        func([]byte) (int, error)
	iface         io.Writer
	ifaceScalar   io.Writer
	_             any
	arrPtr        *[10]byte

	Simple
}

func TestClone_UnexportedFields(t *testing.T) {
	unexported := &Unexported{
		insider: insider{
			i:    -1,
			i8:   -8,
			i16:  -16,
			i32:  -32,
			i64:  -64,
			u:    1,
			u8:   8,
			u16:  16,
			u32:  32,
			u64:  64,
			uptr: uintptr(0xDEADC0DE),
			b:    true,
			s:    "hello",
			f32:  3.2,
			f64:  6.4,
			c64:  complex(6, 4),
			c128: complex(12, 8),
			arr: [4]string{
				"a", "b", "c", "d",
			},
			arrPtr: &[10]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			ch:     make(chan bool, 5),
			fn: func(s string) string {
				return s + ", world!"
			},
			method:      bytes.NewBufferString("method").Write,
			iface:       bytes.NewBufferString("interface"),
			ifaceScalar: scalarWriter(123),
			m: map[string]any{
				"key": "value",
			},
			unsafePointer: unsafe.Pointer(&Unexported{}),
			st: Simple{
				Foo: 123,
				Bar: "bar1",
			},
			Simple: Simple{
				Foo: 456,
				Bar: "bar2",
			},
			t: reflect.TypeOf(&Simple{}),
		},
	}
	a := assert.New(t)
	unexported.m["loop"] = &unexported.m

	// Make pointer cycles.
	unexported.ptr = unexported
	unexported.slice = []*Unexported{unexported}
	dst := Deep(unexported).(*Unexported)
	a.Use(&unexported, &dst)

	// unsafe.Pointer is shadow copied.
	a.Assert(dst.unsafePointer == unexported.unsafePointer)
	unexported.unsafePointer = nil
	dst.unsafePointer = nil

	// chan cannot be compared, but its buffer can be verified.
	a.Equal(cap(dst.ch), cap(unexported.ch))
	unexported.ch = nil
	dst.ch = nil

	// fn cannot be compared, but it can be called.
	a.Equal(dst.fn("Hello"), unexported.fn("Hello"))
	unexported.fn = nil
	dst.fn = nil

	// method cannot be compared, but it can be called.
	a.Assert(dst.method != nil)
	a.NilError(dst.method([]byte("1234")))
	unexported.method = nil
	dst.method = nil

	// cloned.m["loop"] must be exactly the same map of cloned.m.
	a.Assert(reflect.ValueOf(dst.m["loop"]).Elem().Pointer() == reflect.ValueOf(dst.m).Pointer())

	// Don't test this map in reflect.DeepEqual due to bug in Go.
	// https://github.com/golang/go/issues/33907
	unexported.m["loop"] = nil
	dst.m["loop"] = nil

	// reflect.Type should be copied by value.
	a.Equal(reflect.ValueOf(dst.t).Pointer(), reflect.ValueOf(unexported.t).Pointer())

	// Finally, everything else should equal.
	a.Equal(unexported, dst)

}
func TestClone_UnexportedStructMethod(t *testing.T) {
	a := assert.New(t)

	// Another complex case: clone a struct and a map of struct instead of ptr to a struct.
	src := insider{
		m: map[string]any{
			"insider": insider{
				method: bytes.NewBufferString("method").Write,
			},
		},
	}
	dst := Deep(src).(insider)
	a.Use(&src, &dst)

	// For a struct copy, there is a tricky way to copy method. Test it.
	a.Assert(dst.m["insider"].(insider).method != nil)
	n, err := dst.m["insider"].(insider).method([]byte("1234"))
	a.NilError(err)
	a.Equal(n, 4)
}

func TestShallow(t *testing.T) {
	a := 100
	var b *int = &a
	var c **int = &b
	z := Shallow(c).(**int)
	fmt.Printf("%+v %p\n", a, &a)
	fmt.Printf("%+v %p\n", **z, *z)

	**z = 300
	fmt.Printf("%+v %p\n", a, &a)
	fmt.Printf("%+v %p\n", **z, *z)

	fmt.Println("=====================================================")
	ti, _ := time.ParseInLocation(time.DateTime, "2021-04-22 10:30:00", time.Local)
	src := &D{"aitao", 100, []string{"pingpong", "badminton"}, ti}
	dst := Shallow(src).(*D)
	fmt.Printf("%+v %p %p\n", src, &src, &src.birthday)
	fmt.Printf("%+v %p %p\n\n", dst, &dst, &dst.birthday)

	dst.hobby[0] = "乒乓球"
	dst.birthday = time.Now()
	fmt.Printf("%+v %p %p\n", src, &src, &src.birthday)
	fmt.Printf("%+v %p %p\n", dst, &dst, &dst.birthday)

	fmt.Println("=====================================================")
}

type H struct {
	fn func(a, b int) int
	s  string
}

func TestCopyUnexportedFields(t *testing.T) {
	src := H{func(a, b int) int {
		return a + b
	}, "aitao"}

	dst := Deep(src).(H)
	fmt.Printf("%p %p\n", src.fn, &src.s)
	fmt.Printf("%p %p\n\n", dst.fn, &dst.s)

	dst.s = "lml" // 因为结构体是值类型, 修改dst中的字符串并不会影响src中的字符串
	dst.fn = func(a, b int) int {
		return a * b
	}
	fmt.Println(src.fn(10, 20))
	fmt.Printf("%p %p\n", src.fn, &src.s)
	fmt.Printf("%p %p\n\n", dst.fn, &dst.s)

	fmt.Printf("%+v %p\n", src, &src)
	fmt.Printf("%+v %p\n", dst, &dst)
}

func TestCopy_String(t *testing.T) {
	src := ""
	dst := Deep(src).(string)

	fmt.Printf("%+v %p\n", src, &src)
	fmt.Printf("%+v %p\n", dst, &dst)

	dst = "kqai"
	fmt.Printf("%+v %p\n", src, &src)
	fmt.Printf("%+v %p\n", dst, &dst)
}

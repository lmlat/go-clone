package clone

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

func TestDeepClone_Struct_ZeroValue(t *testing.T) {
	src := A{}
	dst := Deep(src).(A)
	fmt.Printf("%+v\n", src)
	fmt.Printf("%+v\n\n", dst)

	dst.Int = 100
	dst.uints = append(src.uints, 999)
	fmt.Printf("%+v\n", src)
	fmt.Printf("%+v\n", dst)
}

func TestDeepClone_Func(t *testing.T) {
	otherFunc := func(a, b int) int { return a - b }
	srcFunc := func(a, b int) int { return a + b }
	dstFunc := Deep(srcFunc).(func(a, b int) int)
	fmt.Printf("%+v, %p\n", srcFunc(1, 2), srcFunc)
	fmt.Printf("%+v, %p\n", dstFunc(1, 2), dstFunc)

	dstFunc = otherFunc
	fmt.Printf("%+v, %p\n", srcFunc(1, 2), srcFunc)
	fmt.Printf("%+v, %p\n", dstFunc(1, 2), dstFunc)
}

func TestDeepClone_Simple_Slice(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}
	dst := Deep(src).([]int)
	fmt.Printf("%+v, %p\n", src, src)
	fmt.Printf("\n%+v, %p\n", dst, dst)

	dst[3] = 100
	fmt.Printf("%+v, %p\n", src, src)
	fmt.Printf("\n%+v, %p\n", dst, dst)
}

func TestDeepClone_Complex_Slice(t *testing.T) {
	src := []*Test{{100,
		map[string]any{
			"a": 10,
			"b": 9,
			"c": 8,
		}}, {200,
		map[string]any{
			"d": 7,
			"e": 6,
			"f": 5,
		}}, {300,
		map[string]any{
			"g": 4,
			"h": 3,
			"i": 2,
		}}}
	dst := Deep(src).([]*Test)
	dst[0].I = 9999
	fmt.Printf("%+v, %p, %v, %p, %v, %p\n", src, src, src[0].I, &src[0].I, src[0].M, &src[0].M)
	fmt.Printf("%+v, %p, %v, %p, %v, %p\n", dst, dst, dst[0].I, &dst[0].I, dst[0].M, &dst[0].M)
}

func TestDeepClone_Map(t *testing.T) {
	src := map[string]*Test{
		"a": {
			I: 123,
			M: map[string]any{
				"cba": 321,
			},
		},
		"b": {
			I: 456,
			M: map[string]any{
				"ghi": 789,
			},
		},
	}
	dst := Deep(src).(map[string]*Test)

	fmt.Printf("%+v, %p, %v, %v\n", src, src, src["a"].M, src["b"].I)
	fmt.Printf("%+v, %p, %v, %v\n\n", dst, dst, dst["a"].M, dst["b"].I)

	dst["a"].M["cba"] = 999
	dst["b"].I = 888
	fmt.Printf("%+v, %p, %v, %v\n", src, src, src["a"].M, src["b"].I)
	fmt.Printf("%+v, %p, %v, %v\n", dst, dst, dst["a"].M, dst["b"].I)
}

type M struct {
	M map[string]any
	m map[string]any
}

func TestDeepClone_InStruct_Map(t *testing.T) {
	a := assert.New(t)
	m := map[string]any{"a": 1000, "b": 2000, "c": "lml", "d": true}
	src := M{m, m}
	dst := Deep(src).(M)
	fmt.Printf("%+v %p %v %p\n", src.M, &src.M, src.m, &src.m)
	fmt.Printf("%+v %p %v %p\n\n", dst.M, &dst.M, dst.m, &dst.m)

	a.Use(&src, &dst)
	// 字段M和m共享同一个map, 深拷贝后也只会存在一份
	a.Assert(equalsAddr(src.m, src.M))

	// 深拷贝结构体后它们的地址不相等.
	a.Assert(nonEqualsAddr(&src, &dst))

	// 对结构体中的m和M字段都会进行深拷贝操作, 但它们引用的数据是相同的.
	a.Assert(nonEqualsAddr(&src.M, &dst.M))
	a.Assert(nonEqualsAddr(&src.m, &dst.m))

	// 深拷贝后也同样共享同一个map
	a.Assert(equalsAddr(dst.m, dst.M))

	// 修改dst.m的数据会直接影响dst.M中的数据, 但不会影响源数据.
	a.Equal(&dst.m, &dst.M)
	dst.m["b"] = 2222
	a.Equal(&dst.m, &dst.M)

	// 修改dst.m的数据会直接影响dst.M中的数据, 但不会影响源数据.
	a.Equal(&dst.m, &dst.M)
	dst.M["a"] = 1111
	a.Equal(&dst.m, &dst.M)

	fmt.Printf("%+v %p %v %p\n", src.M, &src.M, src.m, &src.m)
	fmt.Printf("%+v %p %v %p\n", dst.M, &dst.M, dst.m, &dst.m)
}
func TestDeepClone_Nil_Pointer(t *testing.T) {
	var uintPtr *uint
	dstUintPtr := Deep(uintPtr).(*uint)
	fmt.Printf("%+v, %p\n", uintPtr, uintPtr)
	fmt.Printf("%+v, %p\n\n", dstUintPtr, dstUintPtr)

	var uints uint = 100
	dstUintPtr = &uints
	fmt.Printf("%+v, %p\n", uintPtr, uintPtr)
	fmt.Printf("%+v, %p\n", *dstUintPtr, dstUintPtr)
}

func TestDeepClone_Pointer(t *testing.T) {
	fmt.Println("====================copy string pointer====================")
	var str string
	strPtr := &str
	dstStrPtr := Deep(strPtr).(*string)
	fmt.Printf("%+v, %v\n", str, *strPtr)
	fmt.Printf("%+v, %v\n\n", str, *dstStrPtr)

	*dstStrPtr = "aitao"
	fmt.Printf("%+v, %v\n", str, *strPtr)
	fmt.Printf("%+v, %v\n", str, *dstStrPtr)

	fmt.Println("====================copy int pointer====================")
	ints := 100
	intPtr := &ints
	dstIntPtr := Deep(intPtr).(*int)
	fmt.Printf("%+v, %v\n", ints, *intPtr)
	fmt.Printf("%+v, %v\n\n", ints, *dstIntPtr)

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

func TestDeepClone_Struct_Pointer(t *testing.T) {
	src := &C{"aitao", 100, time.Now(), []string{"ping pong", "badminton", "football"}}
	dst := Deep(src).(*C)
	fmt.Printf("%+v\n", src)
	fmt.Printf("%+v\n\n", dst)

	dst.name = "kqai"
	fmt.Printf("%+v\n", src)
	fmt.Printf("%+v\n", dst)
}

func TestDeepClone_PrivateStruct_Pointer_ZeroValue(t *testing.T) {
	src := &A{}
	dst := Deep(src).(*A)
	fmt.Printf("%+v\n", src)
	fmt.Printf("%+v\n\n", dst)

	dst.Int = 100
	dst.strings = append(dst.strings, "aitao")
	fmt.Printf("%+v\n", src)
	fmt.Printf("%+v\n", dst)
}

func TestDeepClone_Struct_Time(t *testing.T) {
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
	MapA    map[string]int
	MapB    map[string]*B
	bs      []B
	B
	time time.Time
}

type B struct {
	Vals []string
}

func TestDeepClone_StructA(t *testing.T) {
	src := A{
		Int:    42,
		String: "Aitao",
		uints:  []uint{0, 1, 2, 3},
		MapA:   map[string]int{"a": 1, "b": 2},
		MapB: map[string]*B{
			"hi":  {Vals: []string{"hello", "bonjour"}},
			"bye": {Vals: []string{"good-bye", "au revoir"}},
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

	dst.uints[2] = 1000
	fmt.Printf("%+v %p\n", src, &src)
	fmt.Printf("%+v %p\n", dst, &dst)
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

func TestDeepCopy_Array(t *testing.T) {
	src := [2]*Test{
		{
			I: 123,
			M: map[string]any{
				"abc": 123,
			},
		},
		{
			I: 456,
			M: map[string]any{
				"def": 456,
				"ghi": 789,
			},
		},
	}

	dst := Shallow(src).([2]*Test)

	fmt.Printf("%+v,%p,%+v,%p,%+v,%p\n", src, &src, src[0], &src[0], src[1], &src[1])
	fmt.Printf("%+v,%p,%+v,%p,%+v,%p\n", dst, &dst, dst[0], &dst[0], dst[1], &dst[1])

	dst[0].I = 987
	dst[1].M["ghi"] = 321
	fmt.Printf("%+v,%p,%+v,%p,%+v,%p\n", src, &src, src[0], &src[0], src[1], &src[1])
	fmt.Printf("%+v,%p,%+v,%p,%+v,%p\n", dst, &dst, dst[0], &dst[0], dst[1], &dst[1])
}

func TestDeepCloneReflectType(t *testing.T) {
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

func equalsAddr(a, b any) bool {
	return reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer()
}

func nonEqualsAddr(a, b any) bool {
	return reflect.ValueOf(a).Pointer() != reflect.ValueOf(b).Pointer()
}

func TestDeepClone_UnexportedFields(t *testing.T) {
	src := &Unexported{
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
	src.m["loop"] = &src.m
	src.ptr = src
	src.slice = []*Unexported{src}
	dst := Deep(src).(*Unexported)

	// array is deep copied.
	a.Assert(!equalsAddr(&src.arr, &dst.arr))

	fmt.Printf("address: %p %p\n", &src.m, &dst.m)
	fmt.Println("map:", dst.m, src.m)

	fmt.Printf("%+v %p\n\n", src, &src)
	fmt.Printf("%+v %p\n", dst, &dst)

	a.Use(&src, &dst)

	// unsafe.Pointer is shallow copied.
	a.Assert(dst.unsafePointer == src.unsafePointer)
	src.unsafePointer = nil
	dst.unsafePointer = nil

	// chan cannot be directly compared.
	fmt.Println("channel capacity:", cap(src.ch), cap(dst.ch))
	a.Equal(cap(dst.ch), cap(src.ch))
	src.ch = nil
	dst.ch = nil

	// func cannot be directly compared.
	a.Equal(dst.fn("Hello"), src.fn("Hello"))
	src.fn = nil
	dst.fn = nil

	// method cannot be compared.
	a.Assert(dst.method != nil)
	a.NilError(dst.method([]byte("1234")))
	src.method = nil
	dst.method = nil

	// dst.m["loop"] must be exactly the same map of dst.m.
	a.Assert(reflect.ValueOf(dst.m["loop"]).Elem().Pointer() == reflect.ValueOf(dst.m).Pointer())
	// Don't test this map in reflect.DeepEqual due to bug in Go.
	// https://github.com/golang/go/issues/33907
	src.m["loop"] = nil
	dst.m["loop"] = nil

	// reflect.Type is shallow copied.
	a.Assert(equalsAddr(src.t, dst.t))

	a.Equal(src, dst)
}
func TestDeepClone_UnexportedStructMethod(t *testing.T) {
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
}

type H struct {
	fn func(a, b int) int
	s  string
}

func TestCopyUnexportedFields(t *testing.T) {
	a := assert.New(t)

	src := H{func(a, b int) int {
		return a + b
	}, "aitao"}
	dst := Shallow(src).(H)
	a.Use(&src, &dst)

	// struct is deep copied.
	a.Assert(nonEqualsAddr(&src, &dst))
	// string is shallow copied.
	a.Assert(nonEqualsAddr(&src.s, &dst.s))
	// func is shallow copied.
	a.Assert(equalsAddr(src.fn, dst.fn))

	// 结构体是值类型, 修改 dst.s 或 dst.fn 都不会影响 src.
	a.Equal(src.s, dst.s)
	dst.s = "lml"
	a.NotEqual(src.s, dst.s)

	a.Equal(src.fn(10, 20), dst.fn(10, 20))
	dst.fn = func(a, b int) int {
		return a * b
	}
	a.NotEqual(src.fn(10, 20), dst.fn(10, 20))
}

func TestDeepClone_String(t *testing.T) {
	src := "aitao"
	dst := Shallow(src).(string)

	fmt.Printf("%+v %p\n", src, &src)
	fmt.Printf("%+v %p\n", dst, &dst)

	dst = "kqai"
	fmt.Printf("%+v %p\n", src, &src)
	fmt.Printf("%+v %p\n", dst, &dst)
}

func TestDeepClone_StructC(t *testing.T) {
	src := C{"aitao", 100, time.Now(), []string{"ping pong", "badminton", "football"}}
	dst := Deep(src).(C)
	fmt.Printf("%+v\n", src)
	fmt.Printf("%+v\n\n", dst)

	dst.name = "哆啦A梦" // modify an unreferenced type value.
	fmt.Printf("%+v\n", src)
	fmt.Printf("%+v\n\n", dst)

	dst.hobby[0] = "乒乓球" // modify a referenced type value.
	fmt.Printf("%+v\n", src)
	fmt.Printf("%+v\n\n", dst)

	// 只拷贝结构体中的已导出字段
	dst = Deep(src, WithOpFlags(OnlyPublicField)).(C)
	fmt.Printf("%+v %p\n", src, &src)
	fmt.Printf("%+v %p\n\n", dst, &dst)

	dst.name = "kqai"
	fmt.Printf("%+v %p\n", src, &src)
	fmt.Printf("%+v %p\n\n", dst, &dst)

	dst = CopyProperties(src).(C)
	fmt.Printf("%+v %p\n", src, &src)
	fmt.Printf("%+v %p\n\n", dst, &dst)
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

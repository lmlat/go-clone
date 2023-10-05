package clone

import (
	"fmt"
	"testing"
	"time"
)

/**
 *
 * @Author AiTao
 * @Date 2023/8/17 11:13
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
	iface       interface{}
	ifaces      []interface{}
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
	Interface   interface{}
	Interfaces  []interface{}
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
		[]interface{}{42, true, "pan-galactic"},
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
		Interfaces:  []interface{}{42, true, "pan-galactic"},
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
	m := make(map[string]any)
	m = map[string]any{"name": "aitao", "age": 100, "sex": true, "birthday": time.Now().Format("2006-01-02")}
	dst := Deep(m).(map[string]any)
	fmt.Printf("%+v, %p\n", m, m)
	fmt.Printf("%+v, %p\n", dst, dst)

	dst["name"] = "kqai"
	fmt.Printf("%+v, %p\n", m, m)
	fmt.Printf("%+v, %p\n", dst, dst)
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
	fmt.Printf("\n%+v\n", dst)

	dst.name = "kqai"
	fmt.Printf("%+v\n", src)
	fmt.Printf("\n%+v\n", dst)
}

func TestClone_PrivateStruct_Pointer_ZeroValue(t *testing.T) {
	src := &PrivateStruct{}
	dst := Deep(src).(*PrivateStruct)
	fmt.Printf("%+v\n", src)
	fmt.Printf("\n%+v\n", dst)
	fmt.Println("=====================modified====================")
	dst.string = "aitao"
	fmt.Printf("%+v\n", src)
	fmt.Printf("\n%+v\n", dst)
}

func TestClone_PublicStruct_Pointer_ZeroValue(t *testing.T) {
	src := &PublicStruct{}
	dst := Deep(src).(*PublicStruct)
	fmt.Printf("%+v\n", src)
	fmt.Printf("\n%+v\n", dst)
	fmt.Println("=====================modified====================")
	dst.String = "aitao"
	fmt.Printf("%+v\n", src)
	fmt.Printf("\n%+v\n", dst)
}

func TestClone_Struct_Time(t *testing.T) {
	src, _ := time.ParseInLocation(time.DateTime, "2023-04-22 10:30:00", time.Local)
	dst := Deep(src).(time.Time)
	fmt.Printf("%+v %p\n", src, &src)
	fmt.Printf("\n%+v %p\n", dst, &dst)
	fmt.Println("=====================modified====================")
	dst = time.Now()
	fmt.Printf("%+v %p\n", src, &src)
	fmt.Printf("\n%+v %p\n", dst, &dst)
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
	dst := Shallow(src).(A)
	fmt.Printf("%+v %p\n", src, &src)
	fmt.Printf("\n%+v %p\n", dst, &dst)
	fmt.Println("================modified================")
	dst.Map["a"] = 1000
	fmt.Printf("%+v %p\n", src, &src)
	fmt.Printf("\n%+v %p\n", dst, &dst)
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

func TestDeep(t *testing.T) {
	originalInt := 42
	clonedInt := Shallow(originalInt).(int)
	fmt.Println("Original Int:", originalInt)
	fmt.Println("Cloned Int:", clonedInt)

	originalStr := "Hello, World!"
	clonedStr := Deep(originalStr).(string)
	fmt.Println("Original String:", originalStr)
	fmt.Println("Cloned String:", clonedStr)

	originalSlice := []int{1, 2, 3, 4, 5}
	clonedSlice := Deep(originalSlice).([]int)
	fmt.Println("Original Slice:", originalSlice)
	fmt.Println("Cloned Slice:", clonedSlice)

	type Person struct {
		name string
		age  int
	}
	originalPerson := Person{"Alice", 30}
	clonedPerson := Deep(originalPerson).(Person)
	fmt.Println("Original Person:", originalPerson)
	fmt.Println("Cloned Person:", clonedPerson)
}

func TestShallow(t *testing.T) {
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
	dst := Shallow(src).(A)
	fmt.Printf("%+v %p\n", src, &src)
	fmt.Printf("\n%+v %p\n", dst, &dst)
	fmt.Println("================modified================")
	dst.Map["a"] = 1000
	fmt.Printf("%+v %p\n", src, &src)
	fmt.Printf("\n%+v %p\n", dst, &dst)
}

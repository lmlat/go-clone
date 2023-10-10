package clone

import "testing"

/**
 *
 * @Author AiTao
 * @Date 2023/10/8 4:55
 * @Url
 **/

type Test struct {
	I int
	M map[string]any
}

type testSimple struct {
	Foo int
	Bar string
}

func BenchmarkSimpleDeepCopy(b *testing.B) {
	orig := &testSimple{
		Foo: 123,
		Bar: "abcd",
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Deep(orig)
	}
}

func BenchmarkComplexDeepCopy(b *testing.B) {
	m := map[string]*Test{
		"abc": {
			I: 123,
			M: map[string]any{
				"abc": 321,
			},
		},
		"def": {
			I: 456,
			M: map[string]any{
				"def": 789,
			},
		},
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Deep(m, AddOpFlags(DeepString))
	}
}

func BenchmarkSimpleShallowCopy(b *testing.B) {
	orig := &testSimple{
		Foo: 123,
		Bar: "abcd",
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Shallow(orig)
	}
}

func BenchmarkComplexShallowCopy(b *testing.B) {
	m := map[string]*Test{
		"abc": {
			I: 123,
			M: map[string]any{
				"abc": 321,
			},
		},
		"def": {
			I: 456,
			M: map[string]any{
				"def": 789,
			},
		},
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Shallow(m)
	}
}

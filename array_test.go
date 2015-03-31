// Copyright 2015 mparaiso<mparaiso@online.fr>. All rights reserve
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package datastruct

import (
	"errors"
	"testing"
)

func expect(t *testing.T, actual interface{}, expected interface{}) {

	if actual != expected {
		t.Error(actual, "should be", expected)
	}
}

func TestLength(t *testing.T) {
	a := NewArray(1, 2, 3)
	expected := 3
	actual := a.Length()
	expect(t, expected, actual)
}

// TestNewArray tests array constructor
func TestNewArray(t *testing.T) {
	a := NewArray(1, 2, 3)
	expect(t, a.At(0), 1)
	// test if pointers are properly supported
	b := NewArray(&struct{ x int }{x: 1})
	expect(t, b.At(0).(*struct{ x int }).x, 1)

}

func TestFilter(t *testing.T) {
	type fixture struct {
		array    ArrayInterface
		callback func(interface{}, int) bool
		expected ArrayInterface
	}
	fixtures := NewArray(
		&fixture{
			NewArray(1, 2, 3, 4),
			func(val interface{}, i int) bool {
				return val.(int)%2 == 0
			},
			NewArray(2, 4),
		},
	)
	fixtures.ForEach(func(el interface{}, i int) {
		fixture := el.(*fixture)
		filtered := fixture.array.Filter(fixture.callback)
		fixture.expected.ForEach(func(el interface{}, i int) {
			expect(t, filtered.At(i), el)
		})
	})
}

func TestNewArrayFrom(t *testing.T) {
	type fixture struct {
		args []interface{}
	}
	myStruct := &struct{ foo int }{foo: 1}
	myOtherStruct := &struct{ foo int }{foo: 2}
	type test struct{ letter string }
	fixtures := NewArray(
		&fixture{[]interface{}{NewArray(1, 2, 3), []int{1, 2, 3}}},
		&fixture{[]interface{}{NewArray("a", "b", "c"), []string{"a", "b", "c"}}},
		&fixture{[]interface{}{NewArray(test{"a"}), []interface{}{test{"a"}}}},
		&fixture{[]interface{}{NewArray(true, false), []bool{true, false}}},
		&fixture{[]interface{}{NewArray(myStruct, myOtherStruct), []*struct{ foo int }{myStruct, myOtherStruct}, func(v interface{}, a ArrayInterface) error {
			var err error
			switch v := v.(type) {
			case []*struct{ foo int }:
				for _, el := range v {
					a.Push(el)
				}
			default:
				err = errors.New("Cant tell type")
			}
			return err
		}}},
	)
	t.Log(fixtures.At(4).(*fixture).args[1])
	fixtures.ForEach(func(val interface{}, i int) {
		var a ArrayInterface
		fix := val.(*fixture)
		if len(fix.args) == 2 {
			a = NewArrayFrom(fix.args[1])
		} else {
			a = NewArrayFrom(fix.args[1], fix.args[2].(func(interface{}, ArrayInterface) error))
		}
		val.(*fixture).args[0].(*Array).ForEach(func(el interface{}, i int) {

			expect(t, a.At(i), el)
		})
	})
}

func TestPush(t *testing.T) {
	a := NewArray()
	a.Push("foo", "bar")
	NewArray("foo", "bar").ForEach(func(val interface{}, i int) {
		expect(t, a.At(i), val)
	})

}

func TestPop(t *testing.T) {
	popped := "bar"
	a := NewArray("foo", popped)
	b := a.Pop()
	if b != "bar" {
		t.Error(b, "should be", popped)
	}
}

func TestForEach(t *testing.T) {
	result := 0
	expected := 10
	a := NewArray(1, 2, 3, 4)
	a.ForEach(func(v interface{}, i int) {
		result += v.(int)
	})
	if result != 10 {
		t.Error(result, "should be", expected)
	}
}

func TestReduce(t *testing.T) {
	a := NewArray(1, 2, 3, 4)
	b := a.Reduce(func(r, v interface{}, i int) interface{} {
		return r.(int) + v.(int)
	}, 0)
	if b != 10 {
		t.Error(b, "should be 10")
	}
}

func TestMap(t *testing.T) {
	a := NewArray(1, 2, 3)
	expected := NewArray(2, 4, 6)
	b := a.Map(func(v interface{}, i int) interface{} {
		return v.(int) * 2
	})
	b.ForEach(func(v interface{}, i int) {
		if v != expected.At(i) {
			t.Errorf("%v should be %v", v, expected.At(i))
		}
	})

}

func TestAt(t *testing.T) {
	a := NewArray("foo", "bar", "baz")
	for i, v := range []string{"foo", "bar", "baz"} {
		if a.At(i) != v {
			t.Errorf("%v should be %v", a.At(i), v)
		}
	}
}

func TestUnshift(t *testing.T) {
	a := NewArray("baz")
	a.Unshift("foo", "bar")
	expected := []string{"bar", "foo", "baz"}
	for i, v := range expected {
		if a.At(i) != v {
			t.Errorf("%v should be %v", a.At(i), v)
		}
	}
}

func TestShift(t *testing.T) {
	a := NewArray("foo", "bar")
	b := a.Shift()
	if a.At(0) != "bar" {
		t.Errorf("a[0] should be bar")
	}
	if b != "foo" {
		t.Error("b should be 'foo'")
	}
}

func TestSplice(t *testing.T) {
	type fixture struct {
		array     ArrayInterface
		arguments []interface{}
		expected  ArrayInterface
	}
	fixtures := NewArray(
		&fixture{
			NewArray(1, 2, 3),
			[]interface{}{0, 1, 4, 5},
			NewArray(4, 5, 2, 3),
		},
	)
	fixtures.ForEach(func(fix interface{}, index int) {
		args := fix.(*fixture).arguments
		fix.(*fixture).array.Splice(args[0].(int), args[1].(int), args[2:]...)
		fix.(*fixture).expected.ForEach(func(val interface{}, i int) {
			expect(t, fix.(*fixture).array.At(i), val)
		})
	})
}

func TestSlice(t *testing.T) {
	type fixture struct {
		array ArrayInterface
		Slice struct {
			Begin int
			End   int
		}
		Expected []int
	}
	array := NewArray(1, 2, 3, 4, 5)
	fixtures := []*fixture{
		&fixture{
			array: array,
			Slice: struct {
				Begin int
				End   int
			}{Begin: 1, End: 4},
			Expected: []int{2, 3, 4},
		},
		&fixture{
			array:    array,
			Slice:    struct{ Begin, End int }{Begin: -3, End: array.Length()},
			Expected: []int{3, 4, 5},
		},
		&fixture{
			array:    array,
			Slice:    struct{ Begin, End int }{Begin: 0, End: -2},
			Expected: []int{1, 2, 3},
		},
	}
	for _, fixture := range fixtures {
		//t.Logf("%+v", fixture)
		actual := fixture.array.Slice(fixture.Slice.Begin, fixture.Slice.End)
		for i, v := range fixture.Expected {
			expect(t, actual.At(i), v)
		}
	}

}

func TestSome(t *testing.T) {
	type fixture struct {
		array    ArrayInterface
		cb       func(interface{}, int) bool
		expected bool
	}
	fixtures := []*fixture{
		&fixture{
			array: NewArray(1, 2, 3, 4, 5, -6),
			cb: func(value interface{}, index int) bool {
				return value.(int) >= 0
			},
			expected: true,
		},
		&fixture{
			array: NewArray(1, 2, 3),
			cb: func(v interface{}, i int) bool {
				return v.(int) == 0
			},
			expected: false,
		},
	}
	for _, fixture := range fixtures {
		expect(t, fixture.array.Some(fixture.cb), fixture.expected)
	}

}

func TestEvery(t *testing.T) {
	type fixture struct {
		array    ArrayInterface
		cb       func(value interface{}, index int) bool
		expected bool
	}
	isOdd := func(v interface{}, i int) bool {
		return v.(int)%2 == 0
	}
	fixtures := NewArray(
		fixture{
			NewArray(2, 4, 6),
			isOdd,
			true,
		},
		fixture{
			NewArray(0, 2, 4, 5),
			isOdd,
			false,
		},
	)
	fixtures.ForEach(func(v interface{}, i int) {
		fixture := v.(fixture)
		expect(t, fixture.array.Every(fixture.cb), fixture.expected)
	})
}

func TestReverse(t *testing.T) {
	type fixture struct {
		array    ArrayInterface
		expected ArrayInterface
	}
	fixtures := NewArray(
		fixture{
			NewArray(1, 2, 3),
			NewArray(3, 2, 1),
		},
	)
	fixtures.ForEach(func(v interface{}, i int) {
		reversed := v.(fixture).array.Reverse()
		reversed.ForEach(func(val interface{}, i int) {
			expect(t, val, v.(fixture).expected.At(i))
		})
	})
}

func TestConcat(t *testing.T) {
	a := NewArray(1, 2)
	b := NewArray(3, 4)
	c := NewArray(5, 6)
	expected := NewArray(1, 2, 3, 4, 5, 6)
	actual := a.Concat(b, c)
	expected.ForEach(func(val interface{}, i int) {
		expect(t, actual.At(i), val)
	})
}

func TestSort(t *testing.T) {
	type fixture struct {
		array    ArrayInterface
		expected ArrayInterface
		callback func(a, b interface{}) bool
	}
	fixtures := NewArray(&fixture{
		NewArray(1, 2, 3),
		NewArray(1, 2, 3),
		func(a, b interface{}) bool {
			if a.(int) <= b.(int) {
				return true
			}
			return false
		},
	}, &fixture{
		NewArray(1, 2, 3),
		NewArray(3, 2, 1),
		func(a, b interface{}) bool {
			if a.(int) <= b.(int) {
				return false
			}
			return true
		},
	})
	fixtures.ForEach(func(v interface{}, i int) {
		fix := v.(*fixture)
		sorted := fix.array.Sort(fix.callback)
		fix.expected.ForEach(func(v interface{}, i int) {
			expect(t, sorted.At(i), v)
		})
	})

}

func TestIndexOf(t *testing.T) {
	type fixture struct {
		array    ArrayInterface
		args     []interface{}
		expected int
	}
	fixtures := NewArray(&fixture{
		NewArray(1, 2, 3),
		[]interface{}{1, 0},
		0,
	},
		&fixture{
			NewArray(4, 5, 6),
			[]interface{}{1, 0},
			-1,
		},
		&fixture{
			NewArray(1, 2, "a", "b"),
			[]interface{}{"a", 1},
			2,
		})
	fixtures.ForEach(func(el interface{}, index int) {
		fixture := el.(*fixture)
		i := fixture.array.IndexOf(fixture.args[0], fixture.args[1].(int))
		expect(t, i, fixture.expected)
	})
}

func TestReduceRight(t *testing.T) {
	type fixture struct {
		array    ArrayInterface
		reducer  func(interface{}, interface{}, int) interface{}
		initial  interface{}
		expected interface{}
	}
	fixtures := NewArray(&fixture{
		NewArray("m", "e", "s", "s", "a", "g", "e"),
		func(word interface{}, letter interface{}, index int) interface{} {
			return word.(string) + letter.(string)
		},
		"",
		"egassem",
	})
	fixtures.ForEach(func(el interface{}, index int) {
		fixture := el.(*fixture)
		expect(t, fixture.array.ReduceRight(fixture.reducer, fixture.initial), fixture.expected)
	})
}

func TestLastIndexOf(t *testing.T) {
	type fixture struct {
		array    ArrayInterface
		args     []interface{}
		expected int
	}
	fixtures := NewArray(&fixture{
		NewArray(1, 2, 3, 1),
		[]interface{}{1, 4},
		3,
	})
	fixtures.ForEach(func(el interface{}, index int) {
		fixture := el.(*fixture)
		i := fixture.array.LastIndexOf(fixture.args[0], fixture.args[1].(int))
		expect(t, i, fixture.expected)
	})
}

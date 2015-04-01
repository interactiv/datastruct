// Copyright 2015 mparaiso<mparaiso@online.fr>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package array

import (
	"fmt"
	"sort"
)

// Array is a alternative structure for the default array implementation
type Array struct {
	array []interface{}
}

// ArrayInterface represents all methods of an array
type ArrayInterface interface {
	Push(values ...interface{}) int
	Pop() interface{}
	At(int) interface{}
	Length() int
	Shift() interface{}
	Unshift(values ...interface{}) int
	Filter(func(interface{}, int) bool) ArrayInterface
	ForEach(func(interface{}, int))
	Reduce(func(interface{}, interface{}, int) interface{}, interface{}) interface{}
	ReduceRight(func(interface{}, interface{}, int) interface{}, interface{}) interface{}
	Map(func(interface{}, int) interface{}) ArrayInterface
	Slice(v ...int) ArrayInterface
	Splice(start int, deleteCount int, items ...interface{}) ArrayInterface
	Some(func(interface{}, int) bool) bool
	Every(func(interface{}, int) bool) bool
	Reverse() ArrayInterface
	Concat(arrays ...ArrayInterface) ArrayInterface
	Sort(func(a, b interface{}) bool) ArrayInterface
	IndexOf(interface{}, int) int
	LastIndexOf(interface{}, int) int
	String() string
	ArrayInterface() []interface{}
}

// New returns a new array
func New(values ...interface{}) ArrayInterface {
	array := &Array{}
	for i := 0; i < len(values); i++ {
		array.Push(values[i])
	}
	return ArrayInterface(array)

}

// At get a value at index
func (a *Array) At(index int) interface{} {
	var result interface{}
	if index < len(a.array) && index >= 0 {
		result = a.array[index]
	}
	return result
}

// Push put values at the end of array
// returns the number of values added
func (a *Array) Push(values ...interface{}) int {
	a.array = append(a.array, values...)
	return len(values)
}

// Pop remove the last value of the array
func (a *Array) Pop() interface{} {
	var result interface{}
	if len(a.array) > 0 {
		result = a.array[len(a.array)-1]
		a.array = a.array[:len(a.array)-1]
	}
	return result
}

// Length returns the number of elements of the array
func (a *Array) Length() int {
	return len(a.array)
}

// Shift removes the first element of the array and returns it
func (a *Array) Shift() interface{} {
	var result interface{}
	if len(a.array) > 0 {
		result = a.array[0]
	}
	a.array = a.array[1:]
	return result
}

// Unshift add elements at index 0 and returns the number of added elements
func (a *Array) Unshift(values ...interface{}) int {

	for _, value := range values {
		a.array = append(append([]interface{}{}, value), a.array...)
	}
	return len(values)
}

// ForEach execute callback on each element of the array
func (a *Array) ForEach(callback func(value interface{}, i int)) {
	for i := 0; i < len(a.array); i++ {
		callback(a.array[i], i)
	}
}

// Reduce folds the array into a single value
func (a *Array) Reduce(callback func(result interface{}, value interface{}, index int) interface{}, initial interface{}) interface{} {
	a.ForEach(func(value interface{}, i int) {
		initial = callback(initial, value, i)
	})
	return initial

}

func (a *Array) ReduceRight(callback func(result interface{}, value interface{}, index int) interface{}, initial interface{}) interface{} {
	result := initial
	for i := len(a.array) - 1; i >= 0; i-- {
		result = callback(result, a.At(i), i)
	}
	return result
}

// Map iterate over array and push the result of callback into a new Array
func (a *Array) Map(callback func(value interface{}, i int) interface{}) ArrayInterface {
	return a.Reduce(func(result interface{}, value interface{}, index int) interface{} {
		result.(ArrayInterface).Push(callback(value, index))
		return result
	}, New()).(ArrayInterface)
}

// Sort sorts an array given a compare function
func (a *Array) Sort(compareFunc func(a, b interface{}) bool) ArrayInterface {
	result := a.Slice()
	sort.Sort(&sorter{result, compareFunc})
	return result
}

// Splice remove elements from the array at a given index and optionally insert new elements
func (a *Array) Splice(start int, deleteCount int, items ...interface{}) ArrayInterface {
	var result ArrayInterface
	head := a.Slice(0, start)
	queue := a.Slice(start + deleteCount)
	result = a.Slice(start, start+deleteCount)
	a.array = []interface{}{}
	for i := 0; i < head.Length(); i++ {
		a.Push(head.At(i))
	}
	a.Push(items...)
	for i := 0; i < queue.Length(); i++ {
		a.Push(queue.At(i))
	}
	return result
}

// Slice returns a copy of a portion of the array
// It takes up to 2 arguments :
// 	- begin int
// 	- end int (excluded)
//
func (a *Array) Slice(beginAndEndValues ...int) ArrayInterface {
	var begin, end int
	values := beginAndEndValues
	if len(values) > 0 {
		begin = values[0]
		if len(values) > 1 {
			end = values[1]
		} else {
			end = a.Length()
		}
	} else {
		return New(a.array...)
	}
	if begin < 0 {
		if a.Length()+begin < 0 {
			begin = 0
		} else {
			begin = a.Length() + begin
		}
	}
	if end < 0 {
		if a.Length()-end < 0 {
			return New()
		}
		end = a.Length() + end

	}
	if end <= begin {
		return New()
	}

	return New(a.array[begin:end]...)
}

// Some returns true if the callback predicate is satisfied
func (a *Array) Some(callback func(v interface{}, index int) bool) bool {
	for i, v := range a.array {
		if callback(v, i) {
			return true
		}
	}
	return false
}

// Every returns true if the callback predicate is true for every element of the array
func (a *Array) Every(callback func(v interface{}, index int) bool) bool {
	for i, value := range a.array {
		if !callback(value, i) {
			return false
		}
	}
	return true
}

// Reverse reverse the order of the elements of the array and returns a new one
func (a *Array) Reverse() ArrayInterface {
	var result ArrayInterface = New()
	for i := a.Length() - 1; i >= 0; i-- {
		result.Push(a.array[i])
	}
	return result
}

// Concat adds arrays to the end of the array and returns an new array
func (a *Array) Concat(arrays ...ArrayInterface) ArrayInterface {
	result := New(a.array...)
	for _, array := range arrays {
		array.ForEach(func(val interface{}, i int) {
			result.Push(val)
		})
	}
	return result
}

//Filter filters elements given a predicate
func (a *Array) Filter(predicate func(interface{}, int) bool) ArrayInterface {
	return a.Reduce(func(result interface{}, el interface{}, index int) interface{} {
		if predicate(el, index) {
			result.(*Array).Push(el)
		}
		return result
	}, New()).(ArrayInterface)
}

func (a *Array) IndexOf(searchElement interface{}, fromIndex int) int {

	for i := fromIndex; i < a.Length(); i++ {
		if a.At(i) == searchElement {
			return i
		}
	}
	return -1
}

// The lastIndexOf() method returns the last index at which a given element
// can be found in the array, or -1 if it is not present. The array is searched backwards, starting at fromIndex.
func (a *Array) LastIndexOf(searchElement interface{}, fromIndex int) int {
	for i := fromIndex; i >= 0; i-- {
		if a.At(i) == searchElement {
			return i
		}
	}

	return -1
}

func (a *Array) ArrayInterface() []interface{} {
	return append([]interface{}{}, a.array...)
}

func (a *Array) String() string {
	return "ArrayInterface[" + a.Reduce(func(result interface{}, el interface{}, index int) interface{} {
		if index == a.Length()-1 {
			return result.(string) + fmt.Sprintf("%+v", el)
		}
		return result.(string) + fmt.Sprintf("%+v", el) + ", "
	}, "").(string) + "]"

}

// NewFrom creates an Array from builtin Go arrays
// support the following types :
// []Bool, []Int, []int8, []int16, []int32, []int64, []uint,
// []uint8, []uint16, []uint32, []uint64, []float32, []float64, []complex64, []complex128 , []string
//
// CAN PANIC
func NewFrom(collection interface{}, delegate ...func(interface{}, ArrayInterface) error) ArrayInterface {
	a := New()

	switch collection := collection.(type) {
	case Array:
		a = collection.Slice().(*Array)
	case []bool:
		for _, el := range collection {
			a.Push(el)
		}
	case []int:
		for _, el := range collection {
			a.Push(el)
		}
	case []int8:
		for _, el := range collection {
			a.Push(el)
		}
	case []int16:
		for _, el := range collection {
			a.Push(el)
		}
	case []int32:
		for _, el := range collection {
			a.Push(el)
		}
	case []int64:
		for _, el := range collection {
			a.Push(el)
		}
	case []uint:
		for _, el := range collection {
			a.Push(el)
		}
	case []uint8:
		for _, el := range collection {
			a.Push(el)
		}
	case []uint16:
		for _, el := range collection {
			a.Push(el)
		}
	case []uint32:
		for _, el := range collection {
			a.Push(el)
		}
	case []uint64:
		for _, el := range collection {
			a.Push(el)
		}
	case []float32:
		for _, el := range collection {
			a.Push(el)
		}
	case []float64:
		for _, el := range collection {
			a.Push(el)
		}
	case []complex64:
		for _, el := range collection {
			a.Push(el)
		}
	case []complex128:
		for _, el := range collection {
			a.Push(el)
		}
	case []struct{}:
		for _, el := range collection {
			a.Push(el)
		}
	case string:
		for _, el := range collection {
			a.Push(el)
		}
	case []string:
		for _, el := range collection {
			a.Push(el)
		}
	case []interface{}:
		for _, el := range collection {
			a.Push(el)
		}
	default:
		if len(delegate) > 0 {
			if err := delegate[0](collection, a); err != nil {
				panic(fmt.Sprintf("can't turn value %+v into an ArrayInterface", collection))
			}
		} else {
			panic(fmt.Sprintf("can't turn value %+v into an ArrayInterface", collection))
		}
	}
	return a
}

// sorter is used for array.Sort
type sorter struct {
	array       ArrayInterface
	compareFunc func(a, b interface{}) bool
}

// Len returns the length of the array
func (s *sorter) Len() int {
	return s.array.Length()
}

// Less compare 2 elements. if i is less than j return true ,else return false
func (s *sorter) Less(i, j int) bool {
	if s.compareFunc(s.array.At(i), s.array.At(j)) == true {
		return true
	}
	return false
}

// Swap swaps 2 elements
func (s *sorter) Swap(i, j int) {
	a := s.array.Splice(i, 1, s.array.At(j))
	s.array.Splice(j, 1, a.At(0))
}

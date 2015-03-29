// Copyright 2015 mparaiso<mparaiso@online.fr>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.


package datastruct

import (
	"reflect"
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
	ForEach(func(interface{}, int))
	Reduce(func(interface{}, interface{}, int) interface{}, interface{}) interface{}
	Map(func(interface{}, int) interface{}) ArrayInterface
	Slice(v ...int) ArrayInterface
	Splice(start int, deleteCount int, items ...interface{}) ArrayInterface
	Some(func(interface{}, int) bool) bool
	Every(func(interface{}, int) bool) bool
	Reverse() ArrayInterface
	Concat(arrays ...ArrayInterface) ArrayInterface
	Sort(func(a, b interface{}) bool) ArrayInterface
}

// NewArray returns a new array
func NewArray(values ...interface{}) *Array {
	array := &Array{}
	for i := 0; i < len(values); i++ {
		array.Push(values[i])
	}
	return array
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

// Map iterate over array and push the result of callback into a new Array
func (a *Array) Map(callback func(value interface{}, i int) interface{}) ArrayInterface {
	return a.Reduce(func(result interface{}, value interface{}, index int) interface{} {
		result.(ArrayInterface).Push(callback(value, index))
		return result
	}, NewArray()).(ArrayInterface)
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
		return NewArray(a.array...)
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
			return NewArray()
		}
		end = a.Length() + end

	}
	if end <= begin {
		return NewArray()
	}

	return NewArray(a.array[begin:end]...)
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
	var result ArrayInterface = NewArray()
	for i := a.Length() - 1; i >= 0; i-- {
		result.Push(a.array[i])
	}
	return result
}

// Concat adds arrays to the end of the array and returns an new array
func (a *Array) Concat(arrays ...ArrayInterface) ArrayInterface {
	result := NewArray(a.array...)
	for _, array := range arrays {
		array.ForEach(func(val interface{}, i int) {
			result.Push(val)
		})
	}
	return result
}

//NewArrayFrom creates an Array from builtin Go arrays
func NewArrayFrom(array interface{}) *Array {
	a := NewArray()
	switch t := array.(type) {
	case Array, ArrayInterface:
		return array.Slice()
	default:
		for _, el := range array {
			a.Push(el)
		}
		return a
	}

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

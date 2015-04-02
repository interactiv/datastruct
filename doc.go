// Copyright 2015 mparaiso<mparaiso@online.fr>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
Package datastruct provides multiple data structures to complement Go default datas structures, such as
Array and many others

Array

Array is a alternative structure for the default array implementation
Array implements ArrayInterface


 	array:=NewArray(0,1,2)
creates an array of 3 elements , 0,1,2

    array.At(0)
get array value at index i ,returns 0

    array.Length()
get array length, returns 3

 	for _,v:=range []int{3,4}{
 		array.Push(v)
 	}
array is now 0,1,2,3,4

	array.Pop()
returns 4 , array is now 0,1,2,3

    array.Shift()
returns 0 , array is now 1,2,3

    array.Unshift(5,6,7)
returns 3(the number of arguments), array is now 7,6,5,1,2,3

	array.ForEach(func(element interface{},index int){
		fmt.Print(element,index)
	})
execute a function on each element of the array

	sum:=array.Reduce(func(value interface{}, element interface{},index int){
		return value.(int)+element.(int)
	},0)
folds the array into a single value,in that case returns the sum of all elements , 7+6+5+1+2+3 = 24

    doubleArray:=array.Map(func(element interface{},index int)interface{
		return element.(int)*2
	})
Map execute a function on each array element and return an new array containing all the results of each function
in the exemple: 14,12,10,2,4,6


*/
package datastruct

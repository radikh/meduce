// Pakcage meduce contains a primitive implementation of mapreduce functions
// similar to other programming languages.
// Simple as wooden stick.
package meduce

import (
	"runtime"
	"sync"
)

// Reduce is a type of function that takes two values
// and returns a result of their composition.
type Reducer[Value any] func(Value, Value) Value

// Reduce applies provided reducer to all values in the iterator.
// It returns first call result of iterator
// if there is only one or zero values in the interator buffer.
func Reduce[Value any](reducer Reducer[Value], input Iterator[Value]) Value {
	result, _ := input()
	for value, ok := input(); ok; value, ok = input() {
		result = reducer(result, value)
	}

	return result
}

// ParallelReduce applies reducer to all values in the iterator concurrently.
// The iterator should be concurrent safe and the reducer should be order-agnostic.
// Number of goroutines is equal to number of CPU cores * 8.
func ParallelReduce[Value any](reducer Reducer[Value], input Iterator[Value]) Value {
	concurrency := runtime.NumCPU() * 8

	results := make([]Value, concurrency)

	wg := sync.WaitGroup{}

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		i := i

		go func() {
			defer wg.Done()
			result, _ := input()
			for value, ok := input(); ok; value, ok = input() {
				result = reducer(result, value)
			}

			results[i] = result
		}()
	}

	wg.Wait()

	result := results[0]
	for i := 1; i < concurrency; i++ {
		result = reducer(result, results[i])
	}

	return result
}

// Mapper is a type of function that does something
// with a value and returns the result.
type Mapper[Value, Result any] func(Value) Result

// Map applies provided mapper to all values in the iterator
// and returns an iterator with resulting values.
func Map[Value, Result any](mapper Mapper[Value, Result], input Iterator[Value]) Iterator[Result] {
	result := make([]Result, 0)

	for value, ok := input(); ok; value, ok = input() {
		result = append(result, mapper(value))
	}

	return SliceIterator(result...)
}

// NewMapIterator returns an iterator that
// applies mapper on call to the iterator.
// It can be used to chain map operations.
// Usefull for lazy evaluation to save some memory.
func NewMapperIterator[Value, Result any](mapper Mapper[Value, Result], input Iterator[Value]) Iterator[Result] {
	return func() (Result, bool) {
		value, ok := input()
		if !ok {
			return *new(Result), false
		}

		return mapper(value), true
	}
}

// Filterer is a type of function that returns
// true if value matches some condition.
type Filterer[Value any] func(Value) bool

// Filter applies provided filterer to the iterator values
// and returns an iterator with values that match the condition.
func Filter[Value any](filterer Filterer[Value], input Iterator[Value]) Iterator[Value] {
	result := make([]Value, 0)

	for value, ok := input(); ok; value, ok = input() {
		if filterer(value) {
			result = append(result, value)
		}
	}

	return SliceIterator(result...)
}

// NewFilterIterator returns an iterator that
// applies filter on call to the iterator.
// It can be used to chain filter operations.
// Usefull for lazy evaluation to save some memory.
func NewFilterIterator[Value any](filterer Filterer[Value], input Iterator[Value]) Iterator[Value] {
	return func() (Value, bool) {
		for value, ok := input(); ok; value, ok = input() {
			if filterer(value) {
				return value, true
			}
		}
		return *new(Value), false
	}
}

// Iterator is a type of function that generates values.
// It returns value and true if there is a value to return
// and zero element of value type and false if there are no more values.
type Iterator[Value any] func() (Value, bool)

// Slice returns all values from the iterator as a slice.
func (iterator Iterator[Value]) Slice() []Value {
	result := make([]Value, 0)

	for value, ok := iterator(); ok; value, ok = iterator() {
		result = append(result, value)
	}

	return result
}

// SliceIterator returns an iterator that generates values from a slice.
// It returns zero value of Value type and false if there are no more values.
func SliceIterator[Value any](slice ...Value) Iterator[Value] {
	values := make(chan Value, len(slice))

	for _, value := range slice {
		values <- value
	}

	close(values)

	return func() (Value, bool) {
		value, ok := <-values

		return value, ok
	}
}

// JointIterator returns an iterator that generates values from other iterators.
// It returns zero value of Value type and false if there are no more values.
func JointIterator[Value any](iterators ...Iterator[Value]) Iterator[Value] {
	index := 0

	mu := sync.Mutex{}

	return func() (Value, bool) {
		mu.Lock()
		defer mu.Unlock()
		for index < len(iterators) {
			value, ok := iterators[index]()

			if ok {
				return value, true
			}

			index++
		}

		return *new(Value), false
	}
}

// ChannelIterator returns an iterator that generates values from a channel.
// It returns zero value of Value type and false if there are no more values.
func ChannelIterator[Value any](channel chan Value) Iterator[Value] {
	return func() (Value, bool) {
		value, ok := <-channel

		return value, ok
	}
}

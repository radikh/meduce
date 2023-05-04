// Pakcage meduce contains a primitive implementation of mapreduce functions
// similar to other programming languages.
// Simple as wooden stick.
package meduce

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
	index := 0

	return func() (Value, bool) {
		if index >= len(slice) {
			return *new(Value), false
		}

		value := slice[index]
		index++

		return value, true
	}
}

// JointIterator returns an iterator that generates values from other iterators.
// It returns zero value of Value type and false if there are no more values.
func JointIterator[Value any](iterators ...Iterator[Value]) Iterator[Value] {
	index := 0

	return func() (Value, bool) {
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

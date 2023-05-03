package meduce

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReduce(t *testing.T) {
	type testcase struct {
		input    []int
		expected int
	}

	testcases := []testcase{
		{[]int{1, 2, 3, 4, 5}, 15},
		{[]int{1, 2, 3, 4, 5, 6}, 21},
		{[]int{1, 2, 3, 4, 5, 6, 7}, 28},
		{[]int{1, 2, 3, 4, 5, 6, 7, 8}, 36},
		{[]int{}, 0},
	}

	reducer := sumReducer

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("Reduce(%v)", tc.input), func(t *testing.T) {
			result := Reduce(reducer, SliceIterator(tc.input...))
			assert.Equal(t, tc.expected, result)
		})
	}
}

func sumReducer(a, b int) int {
	return a + b
}

func TestMap(t *testing.T) {
	input := SliceIterator([]int{1, 2, 3, 4, 5}...)
	expected := []int{2, 4, 6, 8, 10}

	mapper := MultiplyByTwo

	result := Map(mapper, input)

	assert.Equal(t, expected, result.Slice())
}

func MultiplyByTwo(value int) int {
	return value * 2
}

func TestFilter(t *testing.T) {
	input := SliceIterator([]int{1, 2, 3, 4, 5}...)
	expected := []int{2, 4}

	filterer := Even

	result := Filter(filterer, input)

	assert.Equal(t, expected, result.Slice())
}

func TestReduceMapFilter(t *testing.T) {
	input := SliceIterator([]int{1, 2, 3, 4, 5}...)
	expected := 12

	result := Reduce(sumReducer, Map(MultiplyByTwo, Filter(Even, input)))

	assert.Equal(t, expected, result)
}

func Even(value int) bool {
	return value%2 == 0
}

func ExampleMap_string() {
	input := SliceIterator([]string{"a", "B", "c", "D", "e"}...)

	mapper := func(value string) string {
		return strings.ToUpper(value)
	}

	result := Map(mapper, input)

	fmt.Println(result.Slice())

	//Output:
	//[A B C D E]
}

func ExampleMap_itoa() {
	// You can transform values to different type as well.
	input := SliceIterator([]int{1, 2, 3, 4, 5}...)

	mapper := func(value int) string {
		return fmt.Sprintf("%d", value)
	}

	result := Map(mapper, input)

	fmt.Println(result.Slice())

	//Output:
	//[1 2 3 4 5]
}

func ExampleFilter() {
	input := SliceIterator([]int{1, 2, 3, 4, 5}...)

	filterer := func(value int) bool {
		return value%2 == 0
	}

	result := Filter(filterer, input)

	fmt.Println(result.Slice())

	//Output:
	//[2 4]

}

func ExampleReduce() {
	input := SliceIterator([]int{1, 2, 3, 4, 5}...)

	reducer := func(a, b int) int {
		return a + b
	}

	result := Reduce(reducer, input)

	fmt.Println(result)

	//Output:
	//15
}

func ExampleReduceMapFilter() {
	input := SliceIterator([]int{1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144}...)

	evenFilter := func(value int) bool {
		return value%2 == 0
	}

	itoaMap := func(value int) string {
		return fmt.Sprintf("%d", value)
	}

	longFilter := func(value string) bool {
		return len(value) > 1
	}

	concatenationReducer := func(a, b string) string {
		return fmt.Sprintf("%s%s", a, b)
	}

	result := Reduce(
		concatenationReducer, // will produce "34144"
		Filter(
			longFilter, // will produce ["34" "144"]
			Map(
				itoaMap, // will produce ["2" "8" "34" "144"]
				Filter(
					evenFilter, // will produce [2 8 34 144]
					input,
				),
			),
		),
	)

	fmt.Println(result)

	//Output:
	//34144
}
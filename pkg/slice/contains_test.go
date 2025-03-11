package slice

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

func TestContains(t *testing.T) {
	testCases := []struct {
		target   interface{}
		name     string
		slice    []interface{}
		expected bool
	}{
		{
			name:     "IntSliceFound",
			slice:    []interface{}{1, 2, 3},
			target:   2,
			expected: true,
		},
		{
			name:     "IntSliceNotFound",
			slice:    []interface{}{4, 5, 6},
			target:   7,
			expected: false,
		},
		{
			name:     "StringSliceFound",
			slice:    []interface{}{"apple", "banana", "cherry"},
			target:   "banana",
			expected: true,
		},
		{
			name:     "StringSliceNotFound",
			slice:    []interface{}{"mango", "orange", "pear"},
			target:   "grapes",
			expected: false,
		},
		{
			name:     "Float64SliceFound",
			slice:    []interface{}{1.0, 2.0, 3.0},
			target:   2.0,
			expected: true,
		},
		{
			name:     "Float64SliceNotFound",
			slice:    []interface{}{4.0, 5.0, 6.0},
			target:   7.0,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			actual := Contains(
				tc.slice,
				tc.target,
			)

			if actual != tc.expected {
				t.Errorf("Expected %v, got %v for slice %v and target %v", tc.expected, actual, tc.slice, tc.target)
			}
		})
	}
}

func BenchmarkContains(b *testing.B) {
	targets := []int{0, 1, 100}

	var testSlices = [][]int{
		{1, 2, 3},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
	}

	for _, s := range testSlices {
		for _, t := range targets {
			b.Run(fmt.Sprintf("SliceSize:%d_Target:%d", len(s), t), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					Contains(s, t)
				}
			})
		}
	}
}

func TestContainsByFunc(t *testing.T) {
	tests := []struct {
		name  string
		slice []interface{}
		pred  func(n interface{}) bool
		want  bool
	}{
		{
			name:  "Find number greater than 10",
			slice: []interface{}{5, 8, 12, 7},
			pred:  func(n interface{}) bool { return n.(int) > 10 },
			want:  true,
		},
		{
			name:  "Find string 'hello'",
			slice: []interface{}{"world", "hello", "goodbye"},
			pred:  func(s interface{}) bool { return s.(string) == "hello" },
			want:  true,
		},
		{
			name:  "No negative numbers",
			slice: []interface{}{1, 2, 3, 4},
			pred:  func(n interface{}) bool { return n.(int) < 0 },
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ContainsByFunc(tt.slice, tt.pred)
			if got != tt.want {
				t.Errorf("ContainsByFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

// BenchmarkContainsByFuncInt проверяет производительность функции ContainsByFunc
// при работе с большим количеством целочисленных элементов.
func BenchmarkContainsByFuncInt(b *testing.B) {
	size := 1000000
	random := rand.New(rand.NewSource(42))

	intSlice := make([]int, size)
	for i := 0; i < size; i++ {
		intSlice[i] = random.Intn(size)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = ContainsByFunc(intSlice, func(n int) bool { return n == size/2 })
	}
}

// BenchmarkContainsByFuncString проверяет производительность функции ContainsByFunc
// при работе с большим количеством строковых элементов.
func BenchmarkContainsByFuncString(b *testing.B) {
	size := 1000000
	random := rand.New(rand.NewSource(42))

	stringSlice := make([]string, size)
	for i := 0; i < size; i++ {
		stringSlice[i] = strconv.Itoa(random.Intn(size))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = ContainsByFunc(stringSlice, func(s string) bool { return s == "500000" })
	}
}

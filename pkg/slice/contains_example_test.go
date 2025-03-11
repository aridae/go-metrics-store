package slice

import "fmt"

// ExampleContains демонстрирует использование функции Contains.
func ExampleContains() {
	slice := []int{1, 2, 3}
	target := 2

	found := Contains(slice, target)
	fmt.Println(found)

	target = 4
	found = Contains(slice, target)
	fmt.Println(found)

	// Output:
	// true
	// false
}

func ExampleContainsByFunc() {
	// Проверка наличия числа больше 10 в срезе целых чисел
	nums := []int{5, 8, 12, 7}
	result := ContainsByFunc(nums, func(n int) bool { return n > 10 })
	fmt.Println(result)

	// Проверка наличия строки "hello" в срезе строк
	strs := []string{"world", "hello", "goodbye"}
	result = ContainsByFunc(strs, func(s string) bool { return s == "hello" })
	fmt.Println(result)

	// Проверка отсутствия числа меньше 0 в срезе целых чисел
	nums = []int{1, 2, 3, 4}
	result = ContainsByFunc(nums, func(n int) bool { return n < 0 })
	fmt.Println(result)

	// Output:
	// true
	// true
	// false
}

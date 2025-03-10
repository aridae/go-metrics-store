package slice

// Contains проверяет наличие заданного элемента в срезе.
//
// Эта функция проходит по каждому элементу среза и возвращает true, если находит
// целевой элемент, иначе возвращает false.
//
// Параметры:
//
//	slice - Срез значений параметризованного типа T.
//	target - Элемент, наличие которого проверяется в срезе.
//
// Возвращаемое значение:
//
//	True, если элемент найден в срезе, иначе False.
//
// Примеры:
//
//	Contains([]int{1, 2, 3}, 2)      // Возвращает true
//	Contains([]string{"a", "b"}, "c") // Возвращает false
//
// Ошибки:
//
//	Функция не вызывает ошибок.
func Contains[T comparable](slice []T, target T) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}

	return false
}

func ContainsByFunc[T any](slice []T, predicate func(item T) bool) bool {
	for _, v := range slice {
		if isFound := predicate(v); isFound {
			return true
		}
	}

	return false
}

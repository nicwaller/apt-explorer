package iter

type Iterator[T any] struct {
	Next  func() bool
	Value func() T
}

func Empty[T any]() Iterator[T] {
	return Iterator[T]{
		Next: func() bool {
			return false
		},
		Value: func() T {
			var none T
			return none
		},
	}
}

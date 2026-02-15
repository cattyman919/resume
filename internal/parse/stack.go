package parse

type Stack[T any] struct {
	elements []T
}

func NewStack[T any]() *Stack[T] {
	stack := &Stack[T]{
		elements: make([]T, 0, 20),
	}
	return stack
}

func (s *Stack[T]) Push(element T) {
	s.elements = append(s.elements, element)
}

func (s *Stack[T]) Pop() T {
	lastIndex := len(s.elements) - 1
	popElem := s.elements[lastIndex]
	s.elements = s.elements[:lastIndex]

	return popElem
}

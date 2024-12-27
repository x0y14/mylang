package runtime

import "fmt"

type Stack struct {
	objects []*Object
}

func NewStack(size int) *Stack {
	s := Stack{objects: make([]*Object, size)}
	return &s
}

func (s *Stack) Push(obj *Object) error {
	if len(s.objects) == 0 {
		return fmt.Errorf("failed to push to stack: reason=no space in stack: size=%v", len(s.objects))
	}

	for i := 0; i < len(s.objects); i++ {
		if s.objects[i] == nil {
			s.objects[i] = obj.Clone()
			return nil
		}
	}
	return fmt.Errorf("failed to push to stack: reason=no space in stack: size=%v", len(s.objects))
}
func (s *Stack) Pop() (*Object, error) {
	if len(s.objects) == 0 {
		return nil, fmt.Errorf("failed to pop to stack: reason=stack has no item")
	}
	for i := len(s.objects) - 1; 0 <= i; i-- {
		if s.objects[i] != nil {
			obj := s.objects[i].Clone()
			s.objects[i] = nil
			return obj, nil
		}
	}
	return nil, fmt.Errorf("failed to pop to stack: reason=stack item not found")
}
func (s *Stack) GetSize() int {
	return len(s.objects)
}

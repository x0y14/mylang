package runtime

type Stack struct {
	objects []*Object
}

func NewStack(size int) *Stack {
	return &Stack{objects: make([]*Object, size)}
}

func (s *Stack) Push(obj *Object) {
	s.objects = append(s.objects, obj)
}
func (s *Stack) Pop() *Object {
	last := s.objects[len(s.objects)-1]
	s.objects = s.objects[:len(s.objects)-1]
	return last
}
func (s *Stack) GetSize() int {
	return len(s.objects)
}

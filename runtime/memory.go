package runtime

import "fmt"

func NewMemory(size int) *Memory {
	mem := make(Memory, size)
	return &mem
}

type Memory []*Object

func (m *Memory) SetAt(addr int, obj *Object) error {
	if 0 <= addr && addr < len(*m) {
		(*m)[addr] = obj
		return nil
	}
	return fmt.Errorf("addr must be 0 <= $addr < %d", len(*m))
}
func (m *Memory) GetAt(addr int) *Object {
	return (*m)[addr]
}
func (m *Memory) DeleteAt(addr int) {
	(*m)[addr] = nil
}
func (m *Memory) IsEmptyAt(addr int) bool {
	return (*m)[addr] == nil
}

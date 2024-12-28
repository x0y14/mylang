package runtime

import "fmt"

func NewMemoryV2(size int) *MemoryV2 {
	mem := make(MemoryV2, size)
	return &mem
}

type MemoryV2 []*Object

func (m2 *MemoryV2) SetAt(addr int, obj *Object) error {
	if 0 <= addr && addr < len(*m2) {
		(*m2)[addr] = obj
		return nil
	}
	return fmt.Errorf("addr must be 0 <= $addr < %d", len(*m2))
}
func (m2 *MemoryV2) GetAt(addr int) *Object {
	return (*m2)[addr]
}
func (m2 *MemoryV2) DeleteAt(addr int) {
	(*m2)[addr] = nil
}
func (m2 *MemoryV2) IsEmptyAt(addr int) bool {
	return (*m2)[addr] == nil
}

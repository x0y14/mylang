package runtime

import (
	"fmt"
)

type Memory struct {
	mapping map[string]int
	data    []*Object
}

func NewMemory(size int) *Memory {
	return &Memory{
		mapping: map[string]int{},
		data:    make([]*Object, size),
	}
}

func (m *Memory) Set(key string, value *Object) error {
	for address := 0; address < len(m.data); address++ {
		if m.data[address] == nil {
			m.data[address] = value
			m.mapping[key] = address
			return nil
		}
	}
	return fmt.Errorf("failed to set value: reason=no space in memory: size=%v", len(m.data))
}

func (m *Memory) SetAt(address int, value *Object) error {
	_, ok := m.GetAt(address)
	if ok {
		return fmt.Errorf("failed to set value with an address: reason=non null address: address=%v", address)
	}
	m.data[address] = value
	return nil
}

func (m *Memory) Get(key string) (*Object, error) {
	address, ok := m.mapping[key]
	if !ok {
		return nil, fmt.Errorf("failed to get value: reason=non registered key: key=%v", key)
	}
	return m.data[address], nil
}

func (m *Memory) GetAddress(key string) (int, bool) {
	address, ok := m.mapping[key]
	if ok { // mappingに登録されていた
		if v := m.data[address]; v != nil { // 初期値のnilでなくデータが入っていることを確認
			return address, true
		}
	}
	return address, false
}

func (m *Memory) GetAt(address int) (*Object, bool) {
	v := m.data[address]
	if v == nil {
		return v, false
	}
	return v, true
}

func (m *Memory) Delete(key string) error {
	address, ok := m.mapping[key]
	if !ok {
		return fmt.Errorf("failed to delete value: reason=non registered key: key=%v", key)
	}
	m.data[address] = nil
	delete(m.mapping, key)
	return nil
}

func (m *Memory) DeleteAt(address int) error {
	_, ok := m.GetAt(address)
	if !ok {
		return fmt.Errorf("failed to delete value with an address: reason=null address: address=%v", address)
	}
	m.data[address] = nil
	return nil
}

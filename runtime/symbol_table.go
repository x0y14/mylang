package runtime

import "fmt"

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		symbols: make(map[string]int),
	}
}

type SymbolTable struct {
	symbols map[string]int
}

func (s *SymbolTable) Get(name string) (int, error) {
	v, ok := s.symbols[name]
	if !ok {
		return v, fmt.Errorf("failed to get symbol: not registered: %s", name)
	}
	return v, nil
}

func (s *SymbolTable) Set(name string, addr int) error {
	_, err := s.Get(name)
	if err == nil {
		return fmt.Errorf("failed to set symbol: already registered: %s", name)
	}
	s.symbols[name] = addr
	return nil
}

func (s *SymbolTable) Delete(name string) {
	delete(s.symbols, name)
}

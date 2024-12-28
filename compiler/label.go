package compiler

import "fmt"

type LabelCollector struct {
	label   map[string]int
	counter int
}

func NewLabelCollector() *LabelCollector {
	return &LabelCollector{
		label:   make(map[string]int),
		counter: 0,
	}
}
func (lc *LabelCollector) Init() {
	_, _ = lc.Set("main")
}

func (lc *LabelCollector) Get(name string) (int, bool) {
	no, ok := lc.label[name]
	if !ok {
		return 0, ok
	}
	return no, ok
}

func (lc *LabelCollector) Set(name string) (int, error) {
	no, ok := lc.Get(name)
	if ok {
		return no, fmt.Errorf("alredy registered label: %s(id=%d)", name, no)
	}
	lc.label[name] = lc.counter
	lc.counter++
	return lc.label[name], nil
}

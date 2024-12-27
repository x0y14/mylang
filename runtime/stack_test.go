package runtime

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewStack(t *testing.T) {
	assert.Equal(t, NewStack(1), &Stack{make([]*Object, 1)})
	assert.Equal(t, NewStack(10), &Stack{make([]*Object, 10)})
}

func TestStack_GetSize(t *testing.T) {
	stack := Stack{make([]*Object, 1)}
	assert.Equal(t, stack.GetSize(), 1)
	_ = stack.Push(NewNullObject())
	assert.Equal(t, stack.GetSize(), 1)
}

func TestStack_Push(t *testing.T) {
	stack := Stack{make([]*Object, 2)}
	assert.Equal(t, stack.GetSize(), 2)
	_ = stack.Push(NewNullObject())
	assert.Equal(t, stack.GetSize(), 2)
	_ = stack.Push(NewNullObject())
	assert.Equal(t, stack.GetSize(), 2)
}

func TestStack_Pop(t *testing.T) {
	stack := Stack{make([]*Object, 3)}
	err := stack.Push(NewObject('a'))
	assert.Equal(t, nil, err)
	pop, err := stack.Pop()
	assert.Equal(t, nil, err)
	assert.Equal(t, pop.String(), "a")
	err = stack.Push(NewObject(true))
	assert.Equal(t, nil, err)
	err = stack.Push(NewListObject(10))
	assert.Equal(t, nil, err)
	pop, err = stack.Pop()
	assert.Equal(t, nil, err)
	assert.Equal(t, pop.String(), "list(10)")
	pop, err = stack.Pop()
	assert.Equal(t, nil, err)
	assert.Equal(t, pop.String(), "true")
}

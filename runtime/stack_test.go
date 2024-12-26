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
	stack := Stack{make([]*Object, 0)}
	assert.Equal(t, stack.GetSize(), 0)
	stack.Push(NewNullObject())
	assert.Equal(t, stack.GetSize(), 1)
}

func TestStack_Push(t *testing.T) {
	stack := Stack{make([]*Object, 0)}
	assert.Equal(t, stack.GetSize(), 0)
	stack.Push(NewNullObject())
	assert.Equal(t, stack.GetSize(), 1)
	stack.Push(NewNullObject())
	assert.Equal(t, stack.GetSize(), 2)
}

func TestStack_Pop(t *testing.T) {
	stack := Stack{make([]*Object, 0)}
	stack.Push(NewObject('a'))
	assert.Equal(t, stack.Pop().String(), "a")
	stack.Push(NewObject(true))
	stack.Push(NewListObject(10))
	assert.Equal(t, stack.Pop().String(), "list(10)")
	assert.Equal(t, stack.Pop().String(), "true")
}

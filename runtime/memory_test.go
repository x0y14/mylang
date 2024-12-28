package runtime

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMemory(t *testing.T) {
	assert.Equal(t, make(Memory, 1), *NewMemory(1))
	assert.Equal(t, make(Memory, 100), *NewMemory(100))
	assert.True(t, (*NewMemory(1))[0] == nil) // 初期値がnilと同等であることを確認
	assert.True(t, (*NewMemory(2))[1] == nil)
	assert.Equal(t, len(*NewMemory(1)), 1) // サイズに応じてdataも変更されていることを確認
	assert.Equal(t, len(*NewMemory(10)), 10)
	assert.Equal(t, len(*NewMemory(100)), 100)
}

func TestMemory_SetAt(t *testing.T) {
	memory := NewMemory(10)
	assert.Nil(t, memory.SetAt(0, NewObject('x')))
	assert.Nil(t, memory.GetAt(1))
	assert.Equal(t, NewObject('x'), memory.GetAt(0))
	assert.Nil(t, memory.SetAt(0, NewObject('y')))
	assert.Equal(t, NewObject('y'), memory.GetAt(0))
}

func TestMemory_GetAt(t *testing.T) {
	memory := NewMemory(10)
	assert.Nil(t, memory.SetAt(0, NewObject('x')))
	assert.Equal(t, NewObject('x'), memory.GetAt(0))
}

func TestMemory_DeleteAt(t *testing.T) {
	memory := NewMemory(10)
	assert.Nil(t, memory.SetAt(0, NewObject('x')))
	assert.Equal(t, NewObject('x'), memory.GetAt(0))
	memory.DeleteAt(0)
	assert.Nil(t, memory.GetAt(0))
}

func TestMemory_IsEmptyAt(t *testing.T) {
	memory := NewMemory(2)
	assert.Nil(t, memory.SetAt(0, NewObject('x')))
	assert.False(t, memory.IsEmptyAt(0))
	assert.True(t, memory.IsEmptyAt(1))
}

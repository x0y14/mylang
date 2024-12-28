package runtime

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMemoryV2(t *testing.T) {
	assert.Equal(t, make(MemoryV2, 1), *NewMemoryV2(1))
	assert.Equal(t, make(MemoryV2, 100), *NewMemoryV2(100))
	assert.True(t, (*NewMemoryV2(1))[0] == nil) // 初期値がnilと同等であることを確認
	assert.True(t, (*NewMemoryV2(2))[1] == nil)
	assert.Equal(t, len(*NewMemoryV2(1)), 1) // サイズに応じてdataも変更されていることを確認
	assert.Equal(t, len(*NewMemoryV2(10)), 10)
	assert.Equal(t, len(*NewMemoryV2(100)), 100)
}

func TestMemoryV2_SetAt(t *testing.T) {
	memory := NewMemoryV2(10)
	assert.Nil(t, memory.SetAt(0, NewObject('x')))
	assert.Nil(t, memory.GetAt(1))
	assert.Equal(t, NewObject('x'), memory.GetAt(0))
	assert.Nil(t, memory.SetAt(0, NewObject('y')))
	assert.Equal(t, NewObject('y'), memory.GetAt(0))
}

func TestMemoryV2_GetAt(t *testing.T) {
	memory := NewMemoryV2(10)
	assert.Nil(t, memory.SetAt(0, NewObject('x')))
	assert.Equal(t, NewObject('x'), memory.GetAt(0))
}

func TestMemoryV2_DeleteAt(t *testing.T) {
	memory := NewMemoryV2(10)
	assert.Nil(t, memory.SetAt(0, NewObject('x')))
	assert.Equal(t, NewObject('x'), memory.GetAt(0))
	memory.DeleteAt(0)
	assert.Nil(t, memory.GetAt(0))
}

func TestMemoryV2_IsEmptyAt(t *testing.T) {
	memory := NewMemoryV2(2)
	assert.Nil(t, memory.SetAt(0, NewObject('x')))
	assert.False(t, memory.IsEmptyAt(0))
	assert.True(t, memory.IsEmptyAt(1))
}

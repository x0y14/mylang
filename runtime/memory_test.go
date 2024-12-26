package runtime

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMemory(t *testing.T) {
	assert.Equal(t, NewMemory(1), &Memory{mapping: make(map[string]int), data: make([]*Object, 1)})
	assert.Equal(t, NewMemory(100), &Memory{mapping: make(map[string]int), data: make([]*Object, 100)})
	assert.True(t, NewMemory(1).data[0] == nil) // 初期値がnilと同等であることを確認
	assert.True(t, NewMemory(2).data[1] == nil)
	assert.Equal(t, len(NewMemory(1).data), 1) // サイズに応じてdataも変更されていることを確認
	assert.Equal(t, len(NewMemory(10).data), 10)
	assert.Equal(t, len(NewMemory(100).data), 100)
}

func TestMemory_Set(t *testing.T) {
	memory := NewMemory(10)
	assert.Equal(t, memory.Set("a", NewObject('x')), nil)
	assert.Equal(t, memory.mapping["a"], 0)
	assert.Equal(t, memory.data[0], NewObject('x'))
}

func TestMemory_SetAt(t *testing.T) {
	memory := NewMemory(3)
	assert.Equal(t, memory.SetAt(2, NewObject('x')), nil)
	assert.Equal(t, len(memory.mapping), 0)         // mapping側に変更が入っていないことを確認
	assert.True(t, memory.data[0] == nil)           // 0番目に変更がないことを確認
	assert.True(t, memory.data[1] == nil)           // 1番目に変更がないことを確認
	assert.Equal(t, memory.data[2], NewObject('x')) // 2番目だけ変更されていることを確認
}

func TestMemory_Get(t *testing.T) {
	memory := NewMemory(3)
	// 存在しない値の取得を試みる, 当然エラー
	v, err := memory.Get("a")
	assert.True(t, v == nil)
	assert.Equal(t, err, fmt.Errorf("failed to get value: reason=non registered key: key=%v", "a"))
	// セットしてみる
	err = memory.Set("b", NewObject('x'))
	assert.True(t, err == nil)
	// 取得してみる
	v, err = memory.Get("b")
	assert.Equal(t, v, NewObject('x'))
	assert.True(t, err == nil)
	// mappingの確認
	assert.Equal(t, len(memory.mapping), 1)
	assert.Equal(t, memory.mapping, map[string]int{"b": 0})
	// dataの確認
	assert.Equal(t, len(memory.data), 3)
	assert.Equal(t, memory.data, []*Object{NewObject('x'), nil, nil})
}

func TestMemory_GetAddress(t *testing.T) {
	memory := NewMemory(3)
	// 存在しない取り出し
	v, ok := memory.GetAddress("a")
	assert.Equal(t, v, 0)
	assert.Equal(t, ok, false)
	// セットしてみる
	err := memory.Set("a", NewObject('x'))
	assert.Equal(t, err, nil)
	// 取得してみる
	v, ok = memory.GetAddress("a")
	assert.Equal(t, v, 0)
	assert.Equal(t, ok, true)
	assert.Equal(t, memory.data[v], NewObject('x'))
}

func TestMemory_GetAt(t *testing.T) {
	memory := NewMemory(2)
	// 存在しないアクセスを試みる
	v, ok := memory.GetAt(0)
	assert.True(t, v == nil)
	assert.Equal(t, ok, false)
	// セットしてみる
	err := memory.Set("a", NewObject('x'))
	assert.Equal(t, err, nil)
	// 取得
	v, ok = memory.GetAt(0)
	assert.Equal(t, v, NewObject('x'))
	assert.Equal(t, ok, true)
	assert.Equal(t, v, memory.data[0])
	// Getで取ったものと同じになることを確認
	v2, err := memory.Get("a")
	assert.Equal(t, v2, v)
	assert.Equal(t, err, nil)
}

func TestMemory_Delete(t *testing.T) {
	memory := NewMemory(2)
	// 存在しないものを削除する
	err := memory.Delete("a")
	assert.Equal(t, err, fmt.Errorf("failed to delete value: reason=non registered key: key=%v", "a"))
	assert.Equal(t, len(memory.mapping), 0) // 存在しないので変わらず
	assert.Equal(t, len(memory.data), 2)    // サイズは変わらないことを確認
	// セットしてみる
	err = memory.Set("a", NewObject('x'))
	assert.Equal(t, err, nil)
	assert.Equal(t, len(memory.mapping), 1) // 追加されたので変わったことを確認
	assert.Equal(t, len(memory.data), 2)    // サイズは変わらないことを確認
	// 削除する
	err = memory.Delete("a")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(memory.mapping), 0) // 削除されたので初期と同じことを確認
	assert.Equal(t, len(memory.data), 2)    // サイズは変わらないことを確認
}

func TestMemory_DeleteAt(t *testing.T) {
	memory := NewMemory(2)
	// 存在しないものを削除してみる
	err := memory.DeleteAt(0)
	assert.Equal(t, err, fmt.Errorf("failed to delete value with an address: reason=null address: address=%v", 0))
	assert.Equal(t, len(memory.mapping), 0)
	assert.Equal(t, len(memory.data), 2)
	// セットしてみる
	err = memory.Set("a", NewObject('x'))
	assert.Equal(t, err, nil)
	assert.Equal(t, len(memory.mapping), 1)
	assert.Equal(t, len(memory.data), 2)
	assert.Equal(t, memory.data[0], NewObject('x')) // データがセットされていることを確認
	// 削除してみる
	err = memory.DeleteAt(0)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(memory.mapping), 1) // Atでデータだけ消したのでマッピングは消えない
	assert.Equal(t, len(memory.data), 2)
	assert.True(t, memory.data[0] == nil)   // 消えていることを確認
	assert.Equal(t, memory.mapping["a"], 0) // マッピングの生存を確認
}

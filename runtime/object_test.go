package runtime

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewObject(t *testing.T) {
	assert.Equal(t, NewObject(39), &Object{kind: OBJ_INT, data: 39})
	assert.Equal(t, NewObject('c'), &Object{kind: OBJ_CHAR, data: int('c')})
	assert.Equal(t, NewObject(false), &Object{kind: OBJ_BOOL, data: 0})
	assert.Equal(t, NewObject(true), &Object{kind: OBJ_BOOL, data: 1})
}

func TestNewNullObject(t *testing.T) {
	assert.Equal(t, NewNullObject(), &Object{kind: OBJ_NULL, data: 0})
}

func TestNewListObject(t *testing.T) {
	assert.Equal(t, NewListObject(3), &Object{kind: OBJ_LIST, data: 3})
}

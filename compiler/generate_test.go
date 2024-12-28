package compiler

import (
	"github.com/stretchr/testify/assert"
	"mylang/runtime"
	"testing"
)

func TestGenerate_Return(t *testing.T) {
	n := &Node{kind: ST_RETURN, lhs: &Node{kind: ST_PRIMITIVE, lhs: &Node{kind: ST_INTEGER, leaf: NewToken(TK_INT, "100")}}}
	prog, err := Generate(n)
	assert.Nil(t, err)
	assert.Equal(t, runtime.Program{
		runtime.NewMoveOp(runtime.REG_STATUS, runtime.NewObject(100)),
		runtime.NewReturnOp(),
	}, prog)
}

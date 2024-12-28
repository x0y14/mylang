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
		runtime.NewMoveOp(runtime.NewRegisterObject(runtime.REG_STATUS), runtime.NewObject(100)),
		runtime.NewReturnOp(),
	}, prog)
}

func TestGenerate_DefineFunction(t *testing.T) {
	n := &Node{
		kind: ST_DEFINE_FUNCTION,
		lhs: &Node{
			kind: ST_FUNCTION_DECLARATION,
			lhs: &Node{
				kind: ST_FUNCTION_HEADER,
				lhs:  &Node{kind: ST_IDENT, leaf: NewToken(TK_IDENT, "main")},
				rhs:  &Node{kind: ST_FUNCTION_ARGUMENTS},
			},
			rhs: &Node{kind: ST_FUNCTION_RETURNS},
		},
		rhs: &Node{
			kind: ST_BLOCK,
			lhs:  &Node{kind: ST_RETURN, lhs: &Node{kind: ST_PRIMITIVE, lhs: &Node{kind: ST_INTEGER, leaf: NewToken(TK_INT, "100")}}},
		},
	}
	prog, err := Generate(n)
	assert.Nil(t, err)
	assert.Equal(t, runtime.Program{
		runtime.NewDefLabelOp(runtime.NewLabelObject(0)),
		runtime.NewMoveOp(runtime.NewRegisterObject(runtime.REG_STATUS), runtime.NewObject(100)),
		runtime.NewReturnOp(),
	}, prog)

	n = &Node{
		kind: ST_DEFINE_FUNCTION,
		lhs: &Node{
			kind: ST_FUNCTION_DECLARATION,
			lhs: &Node{
				kind: ST_FUNCTION_HEADER,
				lhs:  &Node{kind: ST_IDENT, leaf: NewToken(TK_IDENT, "main")},
				rhs: &Node{
					kind: ST_FUNCTION_ARGUMENTS,
					lhs: &Node{
						kind: ST_IDENT,
						leaf: NewToken(TK_IDENT, "arg1"),
						next: &Node{
							kind: ST_IDENT,
							leaf: NewToken(TK_IDENT, "arg2"),
						},
					},
				},
			},
			rhs: &Node{kind: ST_FUNCTION_RETURNS},
		},
		rhs: &Node{
			kind: ST_BLOCK,
			lhs:  &Node{kind: ST_RETURN, lhs: &Node{kind: ST_PRIMITIVE, lhs: &Node{kind: ST_INTEGER, leaf: NewToken(TK_INT, "100")}}},
		},
	}
	prog, err = Generate(n)
	assert.Nil(t, err)
	assert.Equal(t, runtime.Program{
		runtime.NewDefLabelOp(runtime.NewLabelObject(0)),
		runtime.NewPopOp(runtime.NewRegisterObject(runtime.REG_TEMP_1)),
		runtime.NewMoveOp(runtime.NewReferenceObject(2), runtime.NewRegisterObject(runtime.REG_TEMP_1)),
		runtime.NewPopOp(runtime.NewRegisterObject(runtime.REG_TEMP_1)),
		runtime.NewMoveOp(runtime.NewReferenceObject(1), runtime.NewRegisterObject(runtime.REG_TEMP_1)),
		runtime.NewMoveOp(runtime.NewRegisterObject(runtime.REG_STATUS), runtime.NewObject(100)),
		runtime.NewReturnOp(),
	}, prog)
}

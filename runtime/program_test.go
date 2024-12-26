package runtime

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExport(t *testing.T) {
	program := Program{
		&Operation{kind: OP_EXIT},
	}
	assert.Equal(t, "EXIT", Export(program))
	program = Program{
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(30)},
		&Operation{kind: OP_EXIT},
	}
	assert.Equal(t, "MOVE register(GENERAL_1) 30\nEXIT", Export(program))
}

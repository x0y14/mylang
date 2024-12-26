package runtime

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRuntime_Run_Exit(t *testing.T) {
	runtime := NewRuntime(3, 3)
	assert.True(t, nil == runtime.register[REG_PROGRAM_COUNTER])
	err := runtime.Run(Program{
		&Operation{kind: OP_EXIT},
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, runtime.register[REG_PROGRAM_COUNTER].data)
	assert.Equal(t, STAT_SUCCESS, Status(runtime.register[REG_STATUS].data))
}

func TestRuntime_Run_Move(t *testing.T) {
	runtime := NewRuntime(3, 3)
	err := runtime.Run(Program{
		&Operation{kind: OP_MOVE, param1: NewReferenceObject(REG_STATUS), param2: NewObject(1)},
		&Operation{kind: OP_EXIT},
	})
	assert.Equal(t, STAT_ERR, Status(runtime.register[REG_STATUS].data))
	assert.Equal(t, nil, err)

	err = runtime.Run(Program{
		&Operation{kind: OP_MOVE, param1: NewReferenceObject(REG_STATUS), param2: NewObject(999)},
		&Operation{kind: OP_EXIT},
	})
	assert.Equal(t, 999, runtime.register[REG_STATUS].data)
	assert.Equal(t, nil, err)

	err = runtime.Run(Program{
		&Operation{kind: OP_MOVE, param1: NewReferenceObject(REG_GENERAL_1), param2: NewObject(888)},
		&Operation{kind: OP_MOVE, param1: NewReferenceObject(REG_STATUS), param2: NewReferenceObject(REG_GENERAL_1)},
		&Operation{kind: OP_EXIT},
	})
	assert.Equal(t, 888, runtime.register[REG_STATUS].data)
	assert.Equal(t, nil, err)

	err = runtime.Run(Program{
		&Operation{kind: OP_MOVE, param1: NewObject(1), param2: NewObject(1)},
		&Operation{kind: OP_EXIT},
	})
	assert.Equal(t, STAT_ERR, Status(runtime.register[REG_STATUS].data))
	assert.Equal(t, fmt.Errorf("unsupported move value: reason=dest is not REFERENCE: dest=%d", 1), err)
}

func TestRuntime_Run_Add(t *testing.T) {
	runtime := NewRuntime(3, 3)
	err := runtime.Run(Program{
		&Operation{kind: OP_MOVE, param1: NewReferenceObject(REG_GENERAL_1), param2: NewObject(30)},
		&Operation{kind: OP_ADD, param1: NewReferenceObject(REG_GENERAL_1), param2: NewObject(5)},
		&Operation{kind: OP_MOVE, param1: NewReferenceObject(REG_STATUS), param2: NewReferenceObject(REG_GENERAL_1)},
		&Operation{kind: OP_ADD, param1: NewReferenceObject(REG_GENERAL_1), param2: NewObject(5)},
		&Operation{kind: OP_EXIT},
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(35), runtime.register[REG_STATUS])
	assert.Equal(t, NewObject(40), runtime.register[REG_GENERAL_1])
}

func TestRuntime_Run_Sub(t *testing.T) {
	runtime := NewRuntime(3, 3)
	err := runtime.Run(Program{
		&Operation{kind: OP_MOVE, param1: NewReferenceObject(REG_GENERAL_1), param2: NewObject(30)},                  // g1 = 30
		&Operation{kind: OP_SUB, param1: NewReferenceObject(REG_GENERAL_1), param2: NewObject(5)},                    // g1 = 25
		&Operation{kind: OP_MOVE, param1: NewReferenceObject(REG_STATUS), param2: NewReferenceObject(REG_GENERAL_1)}, // status = 25
		&Operation{kind: OP_SUB, param1: NewReferenceObject(REG_GENERAL_1), param2: NewObject(5)},                    // g1 = 20
		&Operation{kind: OP_EXIT},
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(25), runtime.register[REG_STATUS])
	assert.Equal(t, NewObject(20), runtime.register[REG_GENERAL_1])
}
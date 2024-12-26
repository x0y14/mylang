package runtime

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRuntime_Run_Exit(t *testing.T) {
	runtime := NewRuntime(3, 3)
	assert.True(t, nil == runtime.register[REG_PROGRAM_COUNTER])
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewObject(0)},
		&Operation{kind: OP_EXIT},
	})
	err := runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, runtime.register[REG_PROGRAM_COUNTER].data) // [1]exitが読み込まれた後に+1されるから2なんだと思う多分．
	assert.Equal(t, STAT_SUCCESS, Status(runtime.register[REG_STATUS].data))
}

func TestRuntime_Run_Move(t *testing.T) {
	runtime := NewRuntime(3, 3)
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewObject(0)},
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_STATUS), param2: NewObject(1)},
		&Operation{kind: OP_EXIT},
	})
	err := runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, STAT_ERR, Status(runtime.register[REG_STATUS].data))
	assert.Equal(t, nil, err)

	_ = runtime.memory.Delete("l_0")
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewObject(0)},
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_STATUS), param2: NewObject(999)},
		&Operation{kind: OP_EXIT},
	})
	err = runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, 999, runtime.register[REG_STATUS].data)
	assert.Equal(t, nil, err)

	_ = runtime.memory.Delete("l_0")
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewObject(0)},
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(888)},
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_STATUS), param2: NewRegisterObject(REG_GENERAL_1)},
		&Operation{kind: OP_EXIT},
	})
	err = runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, 888, runtime.register[REG_STATUS].data)
	assert.Equal(t, nil, err)

	_ = runtime.memory.Delete("l_0")
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewObject(0)},
		&Operation{kind: OP_MOVE, param1: NewObject(1), param2: NewObject(1)},
		&Operation{kind: OP_EXIT},
	})
	err = runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, STAT_ERR, Status(runtime.register[REG_STATUS].data))
	assert.Equal(t, fmt.Errorf("unsupported move value: reason=dest is not REGISTER: dest=%d", 1), err)
}

func TestRuntime_Run_Add(t *testing.T) {
	runtime := NewRuntime(3, 3)
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewObject(0)},
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(30)},
		&Operation{kind: OP_ADD, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(5)},
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_STATUS), param2: NewRegisterObject(REG_GENERAL_1)},
		&Operation{kind: OP_ADD, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(5)},
		&Operation{kind: OP_EXIT},
	})
	err := runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(35), runtime.register[REG_STATUS])
	assert.Equal(t, NewObject(40), runtime.register[REG_GENERAL_1])
}

func TestRuntime_Run_Sub(t *testing.T) {
	runtime := NewRuntime(3, 3)
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewObject(0)},
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(30)},                 // g1 = 30
		&Operation{kind: OP_SUB, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(5)},                   // g1 = 25
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_STATUS), param2: NewRegisterObject(REG_GENERAL_1)}, // status = 25
		&Operation{kind: OP_SUB, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(5)},                   // g1 = 20
		&Operation{kind: OP_EXIT},
	})
	err := runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(25), runtime.register[REG_STATUS])
	assert.Equal(t, NewObject(20), runtime.register[REG_GENERAL_1])
}

func TestRuntime_Run_Jump(t *testing.T) {
	runtime := NewRuntime(3, 3)
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewObject(0)},
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(30)}, // [0] g1 = 30
		&Operation{kind: OP_JUMP, param1: NewLabelObject(1)},                                       // [1]
		&Operation{kind: OP_ADD, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(30)},  // [2] g1 += 30
		&Operation{kind: OP_DEF_LABEL, param1: NewObject(1)},                                       // 1
		&Operation{kind: OP_SUB, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(5)},   // [3] g1 -= 5
		&Operation{kind: OP_EXIT}, // [5] g1 = 25
	})
	err := runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(25), runtime.register[REG_GENERAL_1])
}

func TestRuntime_Run_Eq(t *testing.T) {
	runtime := NewRuntime(1, 1)
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewObject(0)}, // main:
		&Operation{kind: OP_EQ, param1: NewObject(100), param2: NewObject(100)},
		&Operation{kind: OP_EXIT},
	})
	err := runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(true), runtime.register[REG_BOOL_FLAG])

	_ = runtime.memory.Delete("l_0")
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewObject(0)}, // main:
		&Operation{kind: OP_EQ, param1: NewObject(99), param2: NewObject(100)},
		&Operation{kind: OP_EXIT},
	})
	err = runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(false), runtime.register[REG_BOOL_FLAG])
}

func TestRuntime_Run_Ne(t *testing.T) {
	runtime := NewRuntime(1, 1)
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewObject(0)}, // main:
		&Operation{kind: OP_NE, param1: NewObject(100), param2: NewObject(100)},
		&Operation{kind: OP_EXIT},
	})
	err := runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(false), runtime.register[REG_BOOL_FLAG])

	_ = runtime.memory.Delete("l_0")
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewObject(0)}, // main:
		&Operation{kind: OP_NE, param1: NewObject(99), param2: NewObject(100)},
		&Operation{kind: OP_EXIT},
	})
	err = runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(true), runtime.register[REG_BOOL_FLAG])
}

func TestRuntime_Run_Lt(t *testing.T) {
	runtime := NewRuntime(1, 1)
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewObject(0)}, // main:
		&Operation{kind: OP_LT, param1: NewObject(100), param2: NewObject(100)},
		&Operation{kind: OP_EXIT},
	})
	err := runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(false), runtime.register[REG_BOOL_FLAG])

	_ = runtime.memory.Delete("l_0")
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewObject(0)}, // main:
		&Operation{kind: OP_LT, param1: NewObject(100), param2: NewObject(99)},
		&Operation{kind: OP_EXIT},
	})
	err = runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(false), runtime.register[REG_BOOL_FLAG])

	_ = runtime.memory.Delete("l_0")
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewObject(0)}, // main:
		&Operation{kind: OP_LT, param1: NewObject(99), param2: NewObject(100)},
		&Operation{kind: OP_EXIT},
	})
	err = runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(true), runtime.register[REG_BOOL_FLAG])
}

func TestRuntime_Run_Le(t *testing.T) {
	runtime := NewRuntime(1, 1)
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewObject(0)}, // main:
		&Operation{kind: OP_LE, param1: NewObject(100), param2: NewObject(100)},
		&Operation{kind: OP_EXIT},
	})
	err := runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(true), runtime.register[REG_BOOL_FLAG])

	_ = runtime.memory.Delete("l_0")
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewObject(0)}, // main:
		&Operation{kind: OP_LE, param1: NewObject(100), param2: NewObject(99)},
		&Operation{kind: OP_EXIT},
	})
	err = runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(false), runtime.register[REG_BOOL_FLAG])

	_ = runtime.memory.Delete("l_0")
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewObject(0)}, // main:
		&Operation{kind: OP_LE, param1: NewObject(99), param2: NewObject(100)},
		&Operation{kind: OP_EXIT},
	})
	err = runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(true), runtime.register[REG_BOOL_FLAG])
}

func TestRuntime_Run_JumpTrue(t *testing.T) {
	runtime := NewRuntime(1, 3)
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewObject(0)},                                      // main:
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(0)}, // 	g1 = 0
		&Operation{kind: OP_DEF_LABEL, param1: NewObject(1)},                                      // loop:
		&Operation{kind: OP_ADD, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(1)},  //		g1 += 1
		&Operation{kind: OP_LT, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(5)},   // 	if (g1 < 5) bool_flag = true
		&Operation{kind: OP_JUMP_TRUE, param1: NewLabelObject(1)},                                 //	jt loop
		&Operation{kind: OP_EXIT},
	})
	err := runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(5), runtime.register[REG_GENERAL_1])
}

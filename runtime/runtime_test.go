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
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)},
		&Operation{kind: OP_RETURN},
	})
	err := runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, runtime.register[REG_PROGRAM_COUNTER].data) // [1]exitが読み込まれた後に+1されるから2なんだと思う多分．
	assert.Equal(t, STAT_SUCCESS, Status(runtime.register[REG_STATUS].data))
}

func TestRuntime_Run_Move(t *testing.T) {
	runtime := NewRuntime(3, 3)
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)},
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_STATUS), param2: NewObject(1)},
		&Operation{kind: OP_RETURN},
	})
	err := runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, STAT_ERR, Status(runtime.register[REG_STATUS].data))
	assert.Equal(t, nil, err)

	_ = runtime.memory.Delete("l_0")
	_ = runtime.memory.Delete("l_-1")
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)},
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_STATUS), param2: NewObject(999)},
		&Operation{kind: OP_RETURN},
	})
	err = runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, 999, runtime.register[REG_STATUS].data)
	assert.Equal(t, nil, err)

	_ = runtime.memory.Delete("l_0")
	_ = runtime.memory.Delete("l_-1")
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)},
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(888)},
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_STATUS), param2: NewRegisterObject(REG_GENERAL_1)},
		&Operation{kind: OP_RETURN},
	})
	err = runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, 888, runtime.register[REG_STATUS].data)
	assert.Equal(t, nil, err)

	_ = runtime.memory.Delete("l_0")
	_ = runtime.memory.Delete("l_-1")
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)},
		&Operation{kind: OP_MOVE, param1: NewObject(1), param2: NewObject(1)},
		&Operation{kind: OP_RETURN},
	})
	err = runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, STAT_ERR, Status(runtime.register[REG_STATUS].data))
	assert.Equal(t, fmt.Errorf("unsupported move value: reason=dest is not REGISTER: dest=%d", 1), err)
}

func TestRuntime_Run_Push(t *testing.T) {
	runtime := NewRuntime(4, 3)
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)},
		&Operation{kind: OP_PUSH, param1: NewObject(1)},
		&Operation{kind: OP_PUSH, param1: NewObject(2)},
		&Operation{kind: OP_PUSH, param1: NewObject(3)},
		&Operation{kind: OP_EXIT}, // stackの実験なので
	})
	err := runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(3), runtime.stack.objects[3])
	assert.Equal(t, NewObject(2), runtime.stack.objects[2])
	assert.Equal(t, NewObject(1), runtime.stack.objects[1])
	assert.Equal(t, NewReferenceObject(2), runtime.stack.objects[0])
}
func TestRuntime_Run_Pop(t *testing.T) {
	runtime := NewRuntime(4, 2)
	// 単純なpush, pop
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)},
		&Operation{kind: OP_PUSH, param1: NewObject(1)},
		&Operation{kind: OP_POP, param1: NewRegisterObject(REG_GENERAL_1)},
		&Operation{kind: OP_RETURN},
	})
	err := runtime.CollectLabel()
	assert.Nil(t, err)
	err = runtime.Run()
	assert.Nil(t, err)
	assert.Equal(t, NewObject(1), runtime.register[REG_GENERAL_1])
	// remove main
	_ = runtime.memory.Delete("l_0")
	_ = runtime.memory.Delete("l_-1")
	// popで上書き
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)},
		&Operation{kind: OP_PUSH, param1: NewObject(1)},
		&Operation{kind: OP_POP, param1: NewRegisterObject(REG_GENERAL_1)},
		&Operation{kind: OP_PUSH, param1: NewObject(2)},
		&Operation{kind: OP_POP, param1: NewRegisterObject(REG_GENERAL_1)},
		&Operation{kind: OP_PUSH, param1: NewObject(3)},
		&Operation{kind: OP_POP, param1: NewRegisterObject(REG_GENERAL_1)},
		&Operation{kind: OP_RETURN},
	})
	err = runtime.CollectLabel()
	assert.Nil(t, err)
	err = runtime.Run()
	assert.Nil(t, err)
	assert.Equal(t, NewObject(3), runtime.register[REG_GENERAL_1])
	// remove main
	_ = runtime.memory.Delete("l_0")
	_ = runtime.memory.Delete("l_-1")
	// G2に足してく
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)},
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_GENERAL_2), param2: NewObject(0)},                    // g2 = 0
		&Operation{kind: OP_PUSH, param1: NewObject(1)},                                                              // [1, 0, 0]
		&Operation{kind: OP_PUSH, param1: NewObject(2)},                                                              // [1, 2, 0]
		&Operation{kind: OP_PUSH, param1: NewObject(3)},                                                              // [1, 2, 3]
		&Operation{kind: OP_POP, param1: NewRegisterObject(REG_GENERAL_1)},                                           // g1 = 3
		&Operation{kind: OP_ADD, param1: NewRegisterObject(REG_GENERAL_2), param2: NewRegisterObject(REG_GENERAL_1)}, // g2 += g1 (0 += 3)
		&Operation{kind: OP_POP, param1: NewRegisterObject(REG_GENERAL_1)},                                           // g1 = 2
		&Operation{kind: OP_ADD, param1: NewRegisterObject(REG_GENERAL_2), param2: NewRegisterObject(REG_GENERAL_1)}, // g2 += g1 (3 += 2)
		&Operation{kind: OP_EXIT}, // stackの実験なので
	})
	err = runtime.CollectLabel()
	assert.Nil(t, err)
	err = runtime.Run()
	assert.Nil(t, err)
	assert.Equal(t, NewObject(5), runtime.register[REG_GENERAL_2])
}

func TestRuntime_Run_Call(t *testing.T) {
	runtime := NewRuntime(3, 3)
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)},                                 // [0] l_0
		&Operation{kind: OP_CALL, param1: NewLabelObject(9)},                                      // [1] call l_9
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_GENERAL_2), param2: NewObject(5)}, // [2] move g1 5 // skipされるはず
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(9)},                                 // [3] l_9
		&Operation{kind: OP_POP, param1: NewRegisterObject(REG_GENERAL_1)},                        // [4] g1 = *l_9
		&Operation{kind: OP_RETURN},
	})
	err := runtime.CollectLabel()
	assert.Nil(t, err)
	err = runtime.Run()
	assert.Nil(t, err)
	assert.Nil(t, runtime.register[REG_GENERAL_2])                            // skipされているからMOVEでの移動はないことを確認
	assert.Equal(t, NewReferenceObject(3+2), runtime.register[REG_GENERAL_1]) // 3はプロセスの方で勝手に挿入される
}

func TestRuntime_Run_Return(t *testing.T) {
	runtime := NewRuntime(3, 3)
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(1)},
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(5)},
		&Operation{kind: OP_RETURN},
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)}, // main(l_0):
		&Operation{kind: OP_CALL, param1: NewLabelObject(1)},
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_GENERAL_2), param2: NewRegisterObject(REG_GENERAL_1)},
		&Operation{kind: OP_RETURN},
	})
	err := runtime.CollectLabel()
	assert.Nil(t, err)
	err = runtime.Run()
	assert.Nil(t, err)
	assert.Equal(t, NewObject(5), runtime.register[REG_GENERAL_1])
	assert.Equal(t, NewObject(5), runtime.register[REG_GENERAL_2])
}

func TestRuntime_Run_Add(t *testing.T) {
	runtime := NewRuntime(3, 3)
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)},
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(30)},
		&Operation{kind: OP_ADD, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(5)},
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_STATUS), param2: NewRegisterObject(REG_GENERAL_1)},
		&Operation{kind: OP_ADD, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(5)},
		&Operation{kind: OP_RETURN},
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
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)},
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(30)},                 // g1 = 30
		&Operation{kind: OP_SUB, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(5)},                   // g1 = 25
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_STATUS), param2: NewRegisterObject(REG_GENERAL_1)}, // status = 25
		&Operation{kind: OP_SUB, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(5)},                   // g1 = 20
		&Operation{kind: OP_RETURN},
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
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)},
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(30)}, // [0] g1 = 30
		&Operation{kind: OP_JUMP, param1: NewLabelObject(1)},                                       // [1]
		&Operation{kind: OP_ADD, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(30)},  // [2] g1 += 30
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(1)},                                  // 1
		&Operation{kind: OP_SUB, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(5)},   // [3] g1 -= 5
		&Operation{kind: OP_RETURN}, // [5] g1 = 25
	})
	err := runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(25), runtime.register[REG_GENERAL_1])
}

func TestRuntime_Run_Eq(t *testing.T) {
	runtime := NewRuntime(1, 2)
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)}, // main:
		&Operation{kind: OP_EQ, param1: NewObject(100), param2: NewObject(100)},
		&Operation{kind: OP_RETURN},
	})
	err := runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(true), runtime.register[REG_BOOL_FLAG])

	_ = runtime.memory.Delete("l_0")
	_ = runtime.memory.Delete("l_-1")
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)}, // main:
		&Operation{kind: OP_EQ, param1: NewObject(99), param2: NewObject(100)},
		&Operation{kind: OP_RETURN},
	})
	err = runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(false), runtime.register[REG_BOOL_FLAG])
}

func TestRuntime_Run_Ne(t *testing.T) {
	runtime := NewRuntime(1, 2)
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)}, // main:
		&Operation{kind: OP_NE, param1: NewObject(100), param2: NewObject(100)},
		&Operation{kind: OP_RETURN},
	})
	err := runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(false), runtime.register[REG_BOOL_FLAG])

	_ = runtime.memory.Delete("l_0")
	_ = runtime.memory.Delete("l_-1")
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)}, // main:
		&Operation{kind: OP_NE, param1: NewObject(99), param2: NewObject(100)},
		&Operation{kind: OP_RETURN},
	})
	err = runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(true), runtime.register[REG_BOOL_FLAG])
}

func TestRuntime_Run_Lt(t *testing.T) {
	runtime := NewRuntime(1, 2)
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)}, // main:
		&Operation{kind: OP_LT, param1: NewObject(100), param2: NewObject(100)},
		&Operation{kind: OP_RETURN},
	})
	err := runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(false), runtime.register[REG_BOOL_FLAG])

	_ = runtime.memory.Delete("l_0")
	_ = runtime.memory.Delete("l_-1")
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)}, // main:
		&Operation{kind: OP_LT, param1: NewObject(100), param2: NewObject(99)},
		&Operation{kind: OP_RETURN},
	})
	err = runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(false), runtime.register[REG_BOOL_FLAG])

	_ = runtime.memory.Delete("l_0")
	_ = runtime.memory.Delete("l_-1")
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)}, // main:
		&Operation{kind: OP_LT, param1: NewObject(99), param2: NewObject(100)},
		&Operation{kind: OP_RETURN},
	})
	err = runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(true), runtime.register[REG_BOOL_FLAG])
}

func TestRuntime_Run_Le(t *testing.T) {
	runtime := NewRuntime(1, 2)
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)}, // main:
		&Operation{kind: OP_LE, param1: NewObject(100), param2: NewObject(100)},
		&Operation{kind: OP_RETURN},
	})
	err := runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(true), runtime.register[REG_BOOL_FLAG])

	_ = runtime.memory.Delete("l_0")
	_ = runtime.memory.Delete("l_-1")
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)}, // main:
		&Operation{kind: OP_LE, param1: NewObject(100), param2: NewObject(99)},
		&Operation{kind: OP_RETURN},
	})
	err = runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(false), runtime.register[REG_BOOL_FLAG])

	_ = runtime.memory.Delete("l_0")
	_ = runtime.memory.Delete("l_-1")
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)}, // main:
		&Operation{kind: OP_LE, param1: NewObject(99), param2: NewObject(100)},
		&Operation{kind: OP_RETURN},
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
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)},                                 // main:
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(0)}, // 	g1 = 0
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(1)},                                 // loop:
		&Operation{kind: OP_ADD, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(1)},  //		g1 += 1
		&Operation{kind: OP_LT, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(5)},   // 	if (g1 < 5) bool_flag = true
		&Operation{kind: OP_JUMP_TRUE, param1: NewLabelObject(1)},                                 //	jt loop
		&Operation{kind: OP_RETURN},
	})
	err := runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(5), runtime.register[REG_GENERAL_1])
}

func TestRuntime_Run_SyscallWrite(t *testing.T) {
	runtime := NewRuntime(1, 2)
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)},
		&Operation{kind: OP_SYSCALL_WRITE, param1: NewObject(STD_OUT), param2: NewObject('h')},
		&Operation{kind: OP_SYSCALL_WRITE, param1: NewObject(STD_OUT), param2: NewObject('e')},
		&Operation{kind: OP_SYSCALL_WRITE, param1: NewObject(STD_OUT), param2: NewObject('l')},
		&Operation{kind: OP_SYSCALL_WRITE, param1: NewObject(STD_OUT), param2: NewObject('l')},
		&Operation{kind: OP_SYSCALL_WRITE, param1: NewObject(STD_OUT), param2: NewObject('o')},
		&Operation{kind: OP_SYSCALL_WRITE, param1: NewObject(STD_OUT), param2: NewObject(',')},
		&Operation{kind: OP_SYSCALL_WRITE, param1: NewObject(STD_OUT), param2: NewObject('w')},
		&Operation{kind: OP_SYSCALL_WRITE, param1: NewObject(STD_OUT), param2: NewObject('o')},
		&Operation{kind: OP_SYSCALL_WRITE, param1: NewObject(STD_OUT), param2: NewObject('r')},
		&Operation{kind: OP_SYSCALL_WRITE, param1: NewObject(STD_OUT), param2: NewObject('l')},
		&Operation{kind: OP_SYSCALL_WRITE, param1: NewObject(STD_OUT), param2: NewObject('d')},
		&Operation{kind: OP_SYSCALL_WRITE, param1: NewObject(STD_OUT), param2: NewObject('!')},
		&Operation{kind: OP_SYSCALL_WRITE, param1: NewObject(STD_OUT), param2: NewObject(true)},
		&Operation{kind: OP_SYSCALL_WRITE, param1: NewObject(STD_OUT), param2: NewObject(30)},
		&Operation{kind: OP_SYSCALL_WRITE, param1: NewObject(STD_OUT), param2: NewNullObject()},
		&Operation{kind: OP_SYSCALL_WRITE, param1: NewObject(STD_OUT), param2: NewObject('\n')},
		&Operation{kind: OP_RETURN},
	})
	err := runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
}

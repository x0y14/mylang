package runtime

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
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

	runtime.symbolTable.Delete("l_0")
	runtime.symbolTable.Delete("l_-1")
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

	runtime.symbolTable.Delete("l_0")
	runtime.symbolTable.Delete("l_-1")
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

	runtime.symbolTable.Delete("l_0")
	runtime.symbolTable.Delete("l_-1")
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)},
		&Operation{kind: OP_MOVE, param1: NewObject(1), param2: NewObject(1)},
		&Operation{kind: OP_RETURN},
	})
	err = runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, STAT_ERR, Status(runtime.register[REG_STATUS].data))
	assert.Equal(t, fmt.Errorf("unsupported move value: reason=dest is nor REGISTER, REFERENCE: dest=%d", 1), err)

	runtime.symbolTable.Delete("l_0")
	runtime.symbolTable.Delete("l_-1")
	_ = runtime.Load(Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)},
		&Operation{kind: OP_MOVE, param1: NewReferenceObject(2), param2: NewObject(999)},
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_STATUS), param2: NewReferenceObject(2)},
		&Operation{kind: OP_RETURN},
	})
	err = runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, 999, runtime.register[REG_STATUS].data)
	assert.Equal(t, nil, err)
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
	runtime := NewRuntime(4, 3)
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
	runtime.symbolTable.Delete("l_0")
	runtime.symbolTable.Delete("l_-1")
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
	runtime.symbolTable.Delete("l_0")
	runtime.symbolTable.Delete("l_-1")
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

	runtime.symbolTable.Delete("l_0")
	runtime.symbolTable.Delete("l_-1")
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

	runtime.symbolTable.Delete("l_0")
	runtime.symbolTable.Delete("l_-1")
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

	runtime.symbolTable.Delete("l_0")
	runtime.symbolTable.Delete("l_-1")
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

	runtime.symbolTable.Delete("l_0")
	runtime.symbolTable.Delete("l_-1")
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

	runtime.symbolTable.Delete("l_0")
	runtime.symbolTable.Delete("l_-1")
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

	runtime.symbolTable.Delete("l_0")
	runtime.symbolTable.Delete("l_-1")
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

func TestRuntime_Run_JumpFalse(t *testing.T) {
	runtime := NewRuntime(1, 3)
	_ = runtime.Load(Program{
		// main:
		//   move g1 1
		//   eq g1 0
		//   jf l_1
		//   move g2 100
		//   jump l_2
		// l_1:
		//   move g2 1
		// l_2:
		//   ret
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)},
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(1)},
		&Operation{kind: OP_EQ, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(0)},
		&Operation{kind: OP_JUMP_FALSE, param1: NewLabelObject(1)},
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_GENERAL_2), param2: NewObject(100)},
		&Operation{kind: OP_JUMP, param1: NewLabelObject(2)},
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(1)},
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_GENERAL_2), param2: NewObject(1)},
		&Operation{kind: OP_RETURN},
	})
	err := runtime.CollectLabel()
	assert.Equal(t, nil, err)
	err = runtime.Run()
	assert.Equal(t, nil, err)
	assert.Equal(t, NewObject(1), runtime.register[REG_GENERAL_2])
}

func TestRuntime_Run_SyscallWrite(t *testing.T) {
	// 参考: https://chantsune.github.io/articles/320/
	tmpStdout := os.Stdout // 標準出力を元に戻せるように保存
	r, w, _ := os.Pipe()
	os.Stdout = w // 標準出力の書き込み先を変更

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

	_ = w.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	s := strings.TrimRight(buf.String(), "") // バッファーから文字列へ変換
	os.Stdout = tmpStdout
	assert.Equal(t, "hello,world!true30null\n", s)
}

func TestRuntime_Run_FizzBuzz(t *testing.T) {
	// 参考: https://chantsune.github.io/articles/320/
	tmpStdout := os.Stdout // 標準出力を元に戻せるように保存
	r, w, _ := os.Pipe()
	os.Stdout = w // 標準出力の書き込み先を変更

	runtime := NewRuntime(100, 100)
	_ = runtime.Load(Program{
		// check_x15(l_1):
		//   push g1 // fizzbuzzのメインの数字であるg1の保存
		// loop_c15(l_2):
		//   sub g1 15
		// 	 eq g1 0
		//   jt return_from_x15(l_3)
		//   lt g1 0 // g1 < 0の場合，x15でないことが確定なので, 終わる
		//   jt return_from_x15(l_3)
		//   jump loop_c15(l_2) // もう一回引き算して確認する
		// return_from_x15(l_3):
		//   eq g1 0
		//   pop g1 // g1の復元
		//   ret
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(1)},
		&Operation{kind: OP_PUSH, param1: NewRegisterObject(REG_GENERAL_1)},
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(2)},
		&Operation{kind: OP_SUB, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(15)},
		&Operation{kind: OP_EQ, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(0)},
		&Operation{kind: OP_JUMP_TRUE, param1: NewLabelObject(3)},
		&Operation{kind: OP_LT, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(0)},
		&Operation{kind: OP_JUMP_TRUE, param1: NewLabelObject(3)},
		&Operation{kind: OP_JUMP, param1: NewLabelObject(2)},
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(3)},
		&Operation{kind: OP_EQ, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(0)},
		&Operation{kind: OP_POP, param1: NewRegisterObject(REG_GENERAL_1)},
		&Operation{kind: OP_RETURN},

		// check_x5(l_4):
		//   push g1 // fizzbuzzのメインの数字であるg1の保存
		// loop_c5(l_5):
		//   sub g1 5
		// 	 eq g1 0
		//   jt return_from_x5(l_6)
		//   lt g1 0 // g1 < 0の場合，x5でないことが確定なので, 終わる
		//   jt return_from_x5(l_6)
		//   jump loop_c5(l_5) // もう一回引き算して確認する
		// return_from_x5(l_6):
		//   eq g1 0
		//   pop g1 // g1の復元
		//   ret
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(4)},
		&Operation{kind: OP_PUSH, param1: NewRegisterObject(REG_GENERAL_1)},
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(5)},
		&Operation{kind: OP_SUB, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(5)},
		&Operation{kind: OP_EQ, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(0)},
		&Operation{kind: OP_JUMP_TRUE, param1: NewLabelObject(6)},
		&Operation{kind: OP_LT, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(0)},
		&Operation{kind: OP_JUMP_TRUE, param1: NewLabelObject(6)},
		&Operation{kind: OP_JUMP, param1: NewLabelObject(5)},
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(6)},
		&Operation{kind: OP_EQ, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(0)},
		&Operation{kind: OP_POP, param1: NewRegisterObject(REG_GENERAL_1)},
		&Operation{kind: OP_RETURN},

		// check_x3(l_7):
		//   push g1 // fizzbuzzのメインの数字であるg1の保存
		// loop_c3(l_8):
		//   sub g1 3
		// 	 eq g1 0
		//   jt return_from_x3(l_9)
		//   lt g1 0 // g1 < 0の場合，x3でないことが確定なので, 終わる
		//   jt return_from_x3(l_9)
		//   jump loop_c3(l_8) // もう一回引き算して確認する
		// return_from_x3(l_9):
		//   eq g1 0
		//   pop g1 // g1の復元
		//   ret
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(7)},
		&Operation{kind: OP_PUSH, param1: NewRegisterObject(REG_GENERAL_1)},
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(8)},
		&Operation{kind: OP_SUB, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(3)},
		&Operation{kind: OP_EQ, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(0)},
		&Operation{kind: OP_JUMP_TRUE, param1: NewLabelObject(9)},
		&Operation{kind: OP_LT, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(0)},
		&Operation{kind: OP_JUMP_TRUE, param1: NewLabelObject(9)},
		&Operation{kind: OP_JUMP, param1: NewLabelObject(8)},
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(9)},
		&Operation{kind: OP_EQ, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(0)},
		&Operation{kind: OP_POP, param1: NewRegisterObject(REG_GENERAL_1)},
		&Operation{kind: OP_RETURN},

		// print_fizz(l_10):
		//   syscall_write stdout 'f'
		//   syscall_write stdout 'i'
		//   syscall_write stdout 'z'
		//   syscall_write stdout 'z'
		//   ret
		// print_buzz(l_11):
		//   syscall_write stdout 'b'
		//   syscall_write stdout 'u'
		//   syscall_write stdout 'z'
		//   syscall_write stdout 'z'
		//   ret
		// print_fizzbuzz(l_12):
		//   call print_fizz(l_10)
		//   call print_buzz(l_11)
		//   ret
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(10)},
		&Operation{kind: OP_SYSCALL_WRITE, param1: NewObject(STD_OUT), param2: NewObject('f')},
		&Operation{kind: OP_SYSCALL_WRITE, param1: NewObject(STD_OUT), param2: NewObject('i')},
		&Operation{kind: OP_SYSCALL_WRITE, param1: NewObject(STD_OUT), param2: NewObject('z')},
		&Operation{kind: OP_SYSCALL_WRITE, param1: NewObject(STD_OUT), param2: NewObject('z')},
		&Operation{kind: OP_RETURN},
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(11)},
		&Operation{kind: OP_SYSCALL_WRITE, param1: NewObject(STD_OUT), param2: NewObject('b')},
		&Operation{kind: OP_SYSCALL_WRITE, param1: NewObject(STD_OUT), param2: NewObject('u')},
		&Operation{kind: OP_SYSCALL_WRITE, param1: NewObject(STD_OUT), param2: NewObject('z')},
		&Operation{kind: OP_SYSCALL_WRITE, param1: NewObject(STD_OUT), param2: NewObject('z')},
		&Operation{kind: OP_RETURN},
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(12)},
		&Operation{kind: OP_CALL, param1: NewLabelObject(10)},
		&Operation{kind: OP_CALL, param1: NewLabelObject(11)},
		&Operation{kind: OP_RETURN},

		// fizzbuzz(l_13):
		//   syscall_write g1
		//   syscall_write stdout ' '
		//   call check_x15(l_1)
		//   jf its_is_not_x15(l_14)
		//   call print_fizzbuzz(l_12)
		//   jump its_it_not_x3(l_16) // 重複するので終わる
		// its_is_not_x15(l_14):
		//   call check_x5(l_4)
		// 	 jf its_is_not_x5(l_15)
		// 	 call print_buzz(l_11)
		//   jump its_it_not_x3(l_16) // 重複するので終わる
		// its_is_not_x5(l_15):
		//   call check_x3(l_7)
		//   jf its_it_not_x3(l_16)
		//   call print_fizz(l_10)
		// its_it_not_x3(l_16):
		//   syscall_write stdout '\n'
		//   ret
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(13)},
		&Operation{kind: OP_SYSCALL_WRITE, param1: NewObject(STD_OUT), param2: NewRegisterObject(REG_GENERAL_1)},
		&Operation{kind: OP_SYSCALL_WRITE, param1: NewObject(STD_OUT), param2: NewObject(' ')},
		&Operation{kind: OP_CALL, param1: NewLabelObject(1)},
		&Operation{kind: OP_JUMP_FALSE, param1: NewLabelObject(14)},
		&Operation{kind: OP_CALL, param1: NewLabelObject(12)},
		&Operation{kind: OP_JUMP, param1: NewLabelObject(16)},
		//
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(14)},
		&Operation{kind: OP_CALL, param1: NewLabelObject(4)},
		&Operation{kind: OP_JUMP_FALSE, param1: NewLabelObject(15)},
		&Operation{kind: OP_CALL, param1: NewLabelObject(11)},
		&Operation{kind: OP_JUMP, param1: NewLabelObject(16)},
		//
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(15)},
		&Operation{kind: OP_CALL, param1: NewLabelObject(7)},
		&Operation{kind: OP_JUMP_FALSE, param1: NewLabelObject(16)},
		&Operation{kind: OP_CALL, param1: NewLabelObject(10)},
		//
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(16)},
		&Operation{kind: OP_SYSCALL_WRITE, param1: NewObject(STD_OUT), param2: NewObject('\n')},
		&Operation{kind: OP_RETURN},

		// main(l_0):
		//   mov g1 1
		//   call loop(l_17)
		//   ret
		// loop(l_17):
		//   call fizzbuzz(l_13)
		//   add g1 1
		//   eq g1 101
		//   jf loop(l_17)
		//   ret
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(0)},
		&Operation{kind: OP_MOVE, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(1)},
		&Operation{kind: OP_CALL, param1: NewLabelObject(17)},
		&Operation{kind: OP_RETURN},
		//
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(17)},
		&Operation{kind: OP_CALL, param1: NewLabelObject(13)},
		&Operation{kind: OP_ADD, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(1)},
		&Operation{kind: OP_EQ, param1: NewRegisterObject(REG_GENERAL_1), param2: NewObject(101)},
		&Operation{kind: OP_JUMP_FALSE, param1: NewLabelObject(17)},
		&Operation{kind: OP_RETURN},
	})
	err := runtime.CollectLabel()
	assert.Nil(t, err)
	err = runtime.Run()
	assert.Nil(t, err)

	_ = w.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	s := strings.TrimRight(buf.String(), "") // バッファーから文字列へ変換
	os.Stdout = tmpStdout

	// PYTHON
	//for i in range(1, 101):
	//if i % 3 == 0 and i % 5 == 0:
	//print(f"{i} fizzbuzz")
	//elif i % 3 == 0:
	//print(f"{i} fizz")
	//elif i % 5 == 0:
	//print(f"{i} buzz")
	//else:
	//print(f"{i} ")

	fizzbuzz := `1 
2 
3 fizz
4 
5 buzz
6 fizz
7 
8 
9 fizz
10 buzz
11 
12 fizz
13 
14 
15 fizzbuzz
16 
17 
18 fizz
19 
20 buzz
21 fizz
22 
23 
24 fizz
25 buzz
26 
27 fizz
28 
29 
30 fizzbuzz
31 
32 
33 fizz
34 
35 buzz
36 fizz
37 
38 
39 fizz
40 buzz
41 
42 fizz
43 
44 
45 fizzbuzz
46 
47 
48 fizz
49 
50 buzz
51 fizz
52 
53 
54 fizz
55 buzz
56 
57 fizz
58 
59 
60 fizzbuzz
61 
62 
63 fizz
64 
65 buzz
66 fizz
67 
68 
69 fizz
70 buzz
71 
72 fizz
73 
74 
75 fizzbuzz
76 
77 
78 fizz
79 
80 buzz
81 fizz
82 
83 
84 fizz
85 buzz
86 
87 fizz
88 
89 
90 fizzbuzz
91 
92 
93 fizz
94 
95 buzz
96 fizz
97 
98 
99 fizz
100 buzz
`
	assert.Equal(t, fizzbuzz, s)
}

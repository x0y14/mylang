package runtime

import (
	"fmt"
)

type Runtime struct {
	stack    *Stack
	memory   *Memory
	program  Program
	register Register
}

func NewRuntime(stackSize int, memorySize int) *Runtime {
	return &Runtime{
		stack:    NewStack(stackSize),
		memory:   NewMemory(memorySize),
		program:  nil,
		register: NewRegister(),
	}
}

func (r *Runtime) setProgram(prog []*Operation) {
	r.program = prog
}
func (r *Runtime) setPC(newPC int) {
	r.register[REG_PROGRAM_COUNTER] = NewObject(newPC)
}
func (r *Runtime) setStatus(stat Status) {
	r.register[REG_STATUS] = NewObject(int(stat))
}

func (r *Runtime) consumeOp() *Operation {
	curt := r.program[r.register[REG_PROGRAM_COUNTER].data]
	r.advance()
	return curt
}
func (r *Runtime) advance() {
	r.register[REG_PROGRAM_COUNTER].data++
}

func (r *Runtime) doMove(dest, src *Object) error {
	if dest.kind != OBJ_REFERENCE {
		return fmt.Errorf("unsupported move value: reason=dest is not REFERENCE: dest=%v", dest)
	}
	switch src.kind {
	case OBJ_REFERENCE:
		r.register[RegisterKind(dest.data)] = r.register[RegisterKind(src.data)]
	default:
		r.register[RegisterKind(dest.data)] = src
	}
	return nil
}

func (r *Runtime) Run(program []*Operation) error {
	r.setProgram(program)
	r.setPC(0)
	r.setStatus(STAT_SUCCESS)
programLoop:
	for {
		switch curtOp := r.consumeOp(); {
		case curtOp.kind == OP_EXIT: // EXIT
			break programLoop
		case curtOp.kind == OP_MOVE: // MOVE DEST SRC
			if err := r.doMove(curtOp.param1, curtOp.param2); err != nil {
				r.setStatus(STAT_ERR)
				return err
			}
		default:
			r.setStatus(STAT_ERR)
			return fmt.Errorf("unsupported Op: %s", curtOp.kind.String())
		}
	}
	return nil
}

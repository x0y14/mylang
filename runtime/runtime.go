package runtime

import (
	"fmt"
	"strconv"
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

func (r *Runtime) setProgram(prog Program) {
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
	if dest.kind != OBJ_REGISTER {
		return fmt.Errorf("unsupported move value: reason=dest is not REGISTER: dest=%v", dest)
	}
	switch src.kind {
	case OBJ_REGISTER:
		r.register[RegisterKind(dest.data)] = r.register[RegisterKind(src.data)].Clone()
	default:
		r.register[RegisterKind(dest.data)] = src.Clone()
	}
	return nil
}

func (r *Runtime) doAdd(dest, src *Object) error {
	if dest.kind != OBJ_REGISTER {
		return fmt.Errorf("unsupported add value: reason=dest is not REGISTER: dest=%v", dest)
	}
	switch src.kind {
	case OBJ_REGISTER:
		r.register[RegisterKind(dest.data)].data += r.register[RegisterKind(src.data)].data
	default:
		r.register[RegisterKind(dest.data)].data += src.data
	}
	return nil
}

func (r *Runtime) doSub(dest, src *Object) error {
	if dest.kind != OBJ_REGISTER {
		return fmt.Errorf("unsupported sub value: reason=dest is not REGISTER: dest=%v", dest)
	}
	switch src.kind {
	case OBJ_REGISTER:
		r.register[RegisterKind(dest.data)].data -= r.register[RegisterKind(src.data)].data
	default:
		r.register[RegisterKind(dest.data)].data -= src.data
	}
	return nil
}

func (r *Runtime) doJump(dest *Object) error {
	if dest.kind != OBJ_LABEL {
		return fmt.Errorf("unsupported jump value: reason=dest is not label: dest=%v", dest)
	}
	destAddressObj, err := r.memory.Get("l_" + strconv.Itoa(dest.data))
	if err != nil {
		return err
	}
	r.setPC(destAddressObj.data)
	return nil
}

func (r *Runtime) Load(program Program) error {
	r.setProgram(program)
	return nil
}

func (r *Runtime) LabelCollect() error {
	for pc, op := range r.program {
		if op.kind == OP_DEF_LABEL {
			if err := r.memory.Set("l_"+strconv.Itoa(op.param1.data), NewObject(pc)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Runtime) Run() error {
	entryPointAddressObj, err := r.memory.Get("l_0")
	if err != nil {
		return err
	}
	r.setPC(entryPointAddressObj.data)
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
		case curtOp.kind == OP_ADD:
			if err := r.doAdd(curtOp.param1, curtOp.param2); err != nil {
				r.setStatus(STAT_ERR)
				return err
			}
		case curtOp.kind == OP_SUB:
			if err := r.doSub(curtOp.param1, curtOp.param2); err != nil {
				r.setStatus(STAT_ERR)
				return err
			}
		case curtOp.kind == OP_JUMP:
			if err := r.doJump(curtOp.param1); err != nil {
				r.setStatus(STAT_ERR)
				return err
			}
		case curtOp.kind == OP_DEF_LABEL:
			continue
		default:
			r.setStatus(STAT_ERR)
			return fmt.Errorf("unsupported Op: %s", curtOp.kind.String())
		}
	}
	return nil
}

package runtime

import (
	"fmt"
	"os"
	"strconv"
)

type Runtime struct {
	stack       *Stack
	memory      *Memory
	program     Program
	register    Register
	symbolTable *SymbolTable
}

func NewRuntime(stackSize int, memorySize int) *Runtime {
	return &Runtime{
		stack:       NewStack(stackSize),
		memory:      NewMemory(memorySize),
		program:     nil,
		register:    NewRegister(),
		symbolTable: NewSymbolTable(),
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
	switch dest.kind {
	case OBJ_REGISTER: // 代入先がレジスタ
		switch src.kind {
		case OBJ_REGISTER: // ソースがレジスタ
			r.register[RegisterKind(dest.data)] = r.register[RegisterKind(src.data)].Clone()
			return nil
		case OBJ_REFERENCE: // ソースがメモリ
			if yes := r.memory.IsEmptyAt(src.data); yes { // ソースメモリが空
				return fmt.Errorf("failed to move value: reason=src memory is empty: %v", src)
			}
			r.register[RegisterKind(dest.data)] = r.memory.GetAt(src.data).Clone()
			return nil
		default:
			r.register[RegisterKind(dest.data)] = src.Clone()
			return nil
		}
	case OBJ_REFERENCE: // 代入先がメモリ
		if yes := r.memory.IsEmptyAt(dest.data); !yes { // 宛先メモリにデータが入っている
			return fmt.Errorf("failed to move value: reason=dest memory is not empty: %v", dest)
		}
		switch src.kind {
		case OBJ_REGISTER: // ソースがレジスタ
			if err := r.memory.SetAt(dest.data, r.register[RegisterKind(src.data)].Clone()); err != nil {
				return err
			}
			return nil
		case OBJ_REFERENCE: // ソースがメモリ
			if yes := r.memory.IsEmptyAt(src.data); yes { // ソースメモリが空
				return fmt.Errorf("failed to move value: reason=src memory is empty: %v", src)
			}
			if err := r.memory.SetAt(dest.data, r.memory.GetAt(src.data).Clone()); err != nil {
				return err
			}
			return nil
		default:
			if err := r.memory.SetAt(dest.data, src.Clone()); err != nil {
				return err
			}
			return nil
		}
	default:
		return fmt.Errorf("unsupported move value: reason=dest is nor REGISTER, REFERENCE: dest=%v", dest)
	}
	switch src.kind {
	case OBJ_REGISTER:
		r.register[RegisterKind(dest.data)] = r.register[RegisterKind(src.data)].Clone()
		return nil
	case OBJ_REFERENCE:
		if !r.memory.IsEmptyAt(dest.data) {
			return fmt.Errorf("failed to move value: reason=dest memory is not empty: %v", dest)
		}
		err := r.memory.SetAt(dest.data, src.Clone())
		return err
	default:
		r.register[RegisterKind(dest.data)] = src.Clone()
	}
	return nil
}

func (r *Runtime) doPush(obj1 *Object) error {
	switch {
	case obj1.kind == OBJ_REGISTER:
		if err := r.stack.Push(r.register[RegisterKind(obj1.data)].Clone()); err != nil {
			return err
		}
	default:
		if err := r.stack.Push(obj1.Clone()); err != nil {
			return err
		}
	}
	return nil
}

func (r *Runtime) doPop(dest *Object) error {
	if dest.kind != OBJ_REGISTER {
		return fmt.Errorf("unsupported pop value: reason=dest is not REGISTER: dest=%v", dest)
	}
	pop, err := r.stack.Pop()
	if err != nil {
		return err
	}
	r.register[RegisterKind(dest.data)] = pop.Clone()
	return nil
}

func (r *Runtime) doCall(dest *Object) error {
	if dest.kind != OBJ_LABEL {
		return fmt.Errorf("unsupported call value: reason=dest is not LABEL: dest=%v", dest)
	}
	if err := r.stack.Push(NewReferenceObject(r.register[REG_PROGRAM_COUNTER].data)); err != nil {
		return err
	}
	// ラベル経由で宛先の取り出し
	destAddress, err := r.symbolTable.Get("l_" + strconv.Itoa(dest.data))
	if err != nil {
		return err
	}
	// PCの書き換え
	r.setPC(destAddress)
	return nil
}
func (r *Runtime) doReturn() error {
	dest, err := r.stack.Pop()
	if err != nil {
		return err
	}
	if dest.kind != OBJ_REFERENCE {
		return fmt.Errorf("unsupported retrun value: reason=dest is not REFERENCE: dest=%v", dest)
	}
	// PCの書き換え
	r.setPC(dest.data)
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
	destAddress, err := r.symbolTable.Get("l_" + strconv.Itoa(dest.data))
	if err != nil {
		return err
	}
	r.setPC(destAddress)
	return nil
}

func (r *Runtime) doJumpTrue(dest *Object) error {
	if dest.kind != OBJ_LABEL {
		return fmt.Errorf("unsupported jump_true value: reason=dest is not label: dest=%v", dest)
	}
	if r.register[REG_BOOL_FLAG].IsSame(NewObject(false)) {
		return nil
	} else if r.register[REG_BOOL_FLAG].IsSame(NewObject(true)) {
		destAddress, err := r.symbolTable.Get("l_" + strconv.Itoa(dest.data))
		if err != nil {
			return err
		}
		r.setPC(destAddress)
	} else {
		return fmt.Errorf("unsupported jump_true value: reason=bool_flag has not bool: %v", r.register[REG_BOOL_FLAG].String())
	}

	return nil
}

func (r *Runtime) doJumpFalse(dest *Object) error {
	if dest.kind != OBJ_LABEL {
		return fmt.Errorf("unsupported jump_false value: reason=dest is not label: dest=%v", dest)
	}
	if r.register[REG_BOOL_FLAG].IsSame(NewObject(true)) {
		return nil
	} else if r.register[REG_BOOL_FLAG].IsSame(NewObject(false)) {
		destAddress, err := r.symbolTable.Get("l_" + strconv.Itoa(dest.data))
		if err != nil {
			return err
		}
		r.setPC(destAddress)
	} else {
		return fmt.Errorf("unsupported jump_false value: reason=bool_flag has not bool: %v", r.register[REG_BOOL_FLAG].String())
	}
	return nil
}

func (r *Runtime) doEq(obj1, obj2 *Object) error {
	switch {
	case obj1.kind == OBJ_REGISTER && obj2.kind == OBJ_REGISTER:
		if r.register[RegisterKind(obj1.data)].data == r.register[RegisterKind(obj2.data)].data {
			r.register[REG_BOOL_FLAG] = NewObject(true)
			return nil
		}
	case obj1.kind == OBJ_REGISTER && obj2.kind != OBJ_REGISTER:
		if r.register[RegisterKind(obj1.data)].data == obj2.data {
			r.register[REG_BOOL_FLAG] = NewObject(true)
			return nil
		}
	case obj1.kind != OBJ_REGISTER && obj2.kind == OBJ_REGISTER:
		if obj1.data == r.register[RegisterKind(obj2.data)].data {
			r.register[REG_BOOL_FLAG] = NewObject(true)
			return nil
		}
	case obj1.kind != OBJ_REGISTER && obj2.kind != OBJ_REGISTER:
		if obj1.data == obj2.data {
			r.register[REG_BOOL_FLAG] = NewObject(true)
			return nil
		}
	}
	r.register[REG_BOOL_FLAG] = NewObject(false)
	return nil
}

func (r *Runtime) doNe(obj1, obj2 *Object) error {
	switch {
	case obj1.kind == OBJ_REGISTER && obj2.kind == OBJ_REGISTER:
		if r.register[RegisterKind(obj1.data)].data != r.register[RegisterKind(obj2.data)].data {
			r.register[REG_BOOL_FLAG] = NewObject(true)
			return nil
		}
	case obj1.kind == OBJ_REGISTER && obj2.kind != OBJ_REGISTER:
		if r.register[RegisterKind(obj1.data)].data != obj2.data {
			r.register[REG_BOOL_FLAG] = NewObject(true)
			return nil
		}
	case obj1.kind != OBJ_REGISTER && obj2.kind == OBJ_REGISTER:
		if obj1.data != r.register[RegisterKind(obj2.data)].data {
			r.register[REG_BOOL_FLAG] = NewObject(true)
			return nil
		}
	case obj1.kind != OBJ_REGISTER && obj2.kind != OBJ_REGISTER:
		if obj1.data != obj2.data {
			r.register[REG_BOOL_FLAG] = NewObject(true)
			return nil
		}
	}
	r.register[REG_BOOL_FLAG] = NewObject(false)
	return nil
}

func (r *Runtime) doLt(obj1, obj2 *Object) error {
	switch {
	case obj1.kind == OBJ_REGISTER && obj2.kind == OBJ_REGISTER:
		if r.register[RegisterKind(obj1.data)].data < r.register[RegisterKind(obj2.data)].data {
			r.register[REG_BOOL_FLAG] = NewObject(true)
			return nil
		}
	case obj1.kind == OBJ_REGISTER && obj2.kind != OBJ_REGISTER:
		if r.register[RegisterKind(obj1.data)].data < obj2.data {
			r.register[REG_BOOL_FLAG] = NewObject(true)
			return nil
		}
	case obj1.kind != OBJ_REGISTER && obj2.kind == OBJ_REGISTER:
		if obj1.data < r.register[RegisterKind(obj2.data)].data {
			r.register[REG_BOOL_FLAG] = NewObject(true)
			return nil
		}
	case obj1.kind != OBJ_REGISTER && obj2.kind != OBJ_REGISTER:
		if obj1.data < obj2.data {
			r.register[REG_BOOL_FLAG] = NewObject(true)
			return nil
		}
	}
	r.register[REG_BOOL_FLAG] = NewObject(false)
	return nil
}

func (r *Runtime) doLe(obj1, obj2 *Object) error {
	switch {
	case obj1.kind == OBJ_REGISTER && obj2.kind == OBJ_REGISTER:
		if r.register[RegisterKind(obj1.data)].data <= r.register[RegisterKind(obj2.data)].data {
			r.register[REG_BOOL_FLAG] = NewObject(true)
			return nil
		}
	case obj1.kind == OBJ_REGISTER && obj2.kind != OBJ_REGISTER:
		if r.register[RegisterKind(obj1.data)].data <= obj2.data {
			r.register[REG_BOOL_FLAG] = NewObject(true)
			return nil
		}
	case obj1.kind != OBJ_REGISTER && obj2.kind == OBJ_REGISTER:
		if obj1.data <= r.register[RegisterKind(obj2.data)].data {
			r.register[REG_BOOL_FLAG] = NewObject(true)
			return nil
		}
	case obj1.kind != OBJ_REGISTER && obj2.kind != OBJ_REGISTER:
		if obj1.data <= obj2.data {
			r.register[REG_BOOL_FLAG] = NewObject(true)
			return nil
		}
	}
	r.register[REG_BOOL_FLAG] = NewObject(false)
	return nil
}

func (r *Runtime) doSyscallWrite(dest, src *Object) error {
	var f *os.File
	switch {
	case dest.IsSame(NewObject(STD_OUT)):
		f = os.Stdout
	case dest.IsSame(NewObject(STD_ERR)):
		f = os.Stderr
	default:
		return fmt.Errorf("unsupported syscall_write value: reason=dest is nor 2 & 3: dest=%v", dest)
	}
	if src.kind == OBJ_REGISTER {
		_, err := fmt.Fprintf(f, r.register[RegisterKind(src.data)].StringData())
		return err
	}
	_, err := fmt.Fprintf(f, src.StringData())
	return err
}

func (r *Runtime) Load(program Program) error {
	// main(l_0)を叩くコード, exit
	// TODO: startupはコンパイラ側で挿入する
	startup := Program{
		&Operation{kind: OP_DEF_LABEL, param1: NewLabelObject(-1)}, // process root label
		&Operation{kind: OP_CALL, param1: NewLabelObject(0)},       // call main
		&Operation{kind: OP_EXIT},
	}
	program = append(startup, program...)
	r.setProgram(program)
	return nil
}

func (r *Runtime) CollectLabel() error {
	for pc, op := range r.program {
		if op.kind == OP_DEF_LABEL {
			if op.param1.kind != OBJ_LABEL {
				return fmt.Errorf("failed to collect label: failed to define label: reason=this is not label object: obj=%s", op.param1.String())
			}
			if err := r.symbolTable.Set("l_"+strconv.Itoa(op.param1.data), pc); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Runtime) Run() error {
	entryPointAddress, err := r.symbolTable.Get("l_-1")
	if err != nil {
		return err
	}
	r.setPC(entryPointAddress)
	r.setStatus(STAT_SUCCESS)
programLoop:
	for {
		switch curtOp := r.consumeOp(); {
		case curtOp.kind == OP_EXIT: // EXIT
			break programLoop
		case curtOp.kind == OP_MOVE: // MOVE $DEST $SRC
			if err := r.doMove(curtOp.param1, curtOp.param2); err != nil {
				r.setStatus(STAT_ERR)
				return err
			}
		case curtOp.kind == OP_PUSH:
			if err := r.doPush(curtOp.param1); err != nil {
				r.setStatus(STAT_ERR)
				return err
			}
		case curtOp.kind == OP_POP:
			if err := r.doPop(curtOp.param1); err != nil {
				r.setStatus(STAT_ERR)
				return err
			}
		case curtOp.kind == OP_CALL:
			if err := r.doCall(curtOp.param1); err != nil {
				r.setStatus(STAT_ERR)
				return err
			}
		case curtOp.kind == OP_RETURN:
			if err := r.doReturn(); err != nil {
				r.setStatus(STAT_ERR)
				return err
			}
		case curtOp.kind == OP_ADD: // ADD $DEST $SRC
			if err := r.doAdd(curtOp.param1, curtOp.param2); err != nil {
				r.setStatus(STAT_ERR)
				return err
			}
		case curtOp.kind == OP_SUB: // SUB $DEST $SRC
			if err := r.doSub(curtOp.param1, curtOp.param2); err != nil {
				r.setStatus(STAT_ERR)
				return err
			}
		case curtOp.kind == OP_JUMP: // JUMP $LABEL
			if err := r.doJump(curtOp.param1); err != nil {
				r.setStatus(STAT_ERR)
				return err
			}
		case curtOp.kind == OP_DEF_LABEL: // DEF_LABEL $LABEL_NO
			continue
		case curtOp.kind == OP_EQ: // EQ $OBJ1 $OBJ2
			if err := r.doEq(curtOp.param1, curtOp.param2); err != nil {
				r.setStatus(STAT_ERR)
				return err
			}
		case curtOp.kind == OP_NE: // NE $OBJ1 $OBJ2
			if err := r.doNe(curtOp.param1, curtOp.param2); err != nil {
				r.setStatus(STAT_ERR)
				return err
			}
		case curtOp.kind == OP_LT: // LT $OBJ1 $OBJ2
			if err := r.doLt(curtOp.param1, curtOp.param2); err != nil {
				r.setStatus(STAT_ERR)
				return err
			}
		case curtOp.kind == OP_LE: // LE $OBJ1 $OBJ2
			if err := r.doLe(curtOp.param1, curtOp.param2); err != nil {
				r.setStatus(STAT_ERR)
				return err
			}
		case curtOp.kind == OP_JUMP_TRUE: // JUMP_TRUE $LABEL_NO
			if err := r.doJumpTrue(curtOp.param1); err != nil {
				r.setStatus(STAT_ERR)
				return err
			}
		case curtOp.kind == OP_JUMP_FALSE: // JUMP_TRUE $LABEL_NO
			if err := r.doJumpFalse(curtOp.param1); err != nil {
				r.setStatus(STAT_ERR)
				return err
			}
		case curtOp.kind == OP_SYSCALL_WRITE:
			if err := r.doSyscallWrite(curtOp.param1, curtOp.param2); err != nil {
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

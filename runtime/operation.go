package runtime

type OperationKind int

const (
	OP_ILLEGAL OperationKind = iota
	OP_EXIT
	OP_MOVE
	OP_PUSH
	OP_POP
	OP_CALL
	OP_RETURN
	OP_ADD
	OP_SUB
	OP_JUMP
	OP_JUMP_TRUE
	OP_JUMP_FALSE
	OP_DEF_LABEL
	OP_EQ
	OP_NE
	OP_LT
	OP_LE
	OP_SYSCALL_WRITE
)

var opKinds = [...]string{
	OP_ILLEGAL:       "ILLEGAL",
	OP_EXIT:          "EXIT",
	OP_MOVE:          "MOVE",
	OP_PUSH:          "PUSH",
	OP_POP:           "POP",
	OP_CALL:          "CALL",
	OP_RETURN:        "RETURN",
	OP_ADD:           "ADD",
	OP_SUB:           "SUB",
	OP_JUMP:          "JUMP",
	OP_JUMP_TRUE:     "JUMP_TRUE",
	OP_JUMP_FALSE:    "JUMP_FALSE",
	OP_DEF_LABEL:     "DEF_LABEL",
	OP_EQ:            "EQ",
	OP_NE:            "NE",
	OP_LT:            "LT",
	OP_LE:            "LE",
	OP_SYSCALL_WRITE: "SYSCALL_WRITE",
}

func (opKind OperationKind) String() string {
	return opKinds[opKind]
}

type Operation struct {
	kind   OperationKind
	param1 *Object
	param2 *Object
	param3 *Object
	param4 *Object
}

func (op *Operation) String() string {
	str := op.kind.String()
	if op.param1 != nil {
		str += " " + op.param1.String()
	}
	if op.param2 != nil {
		str += " " + op.param2.String()
	}
	if op.param3 != nil {
		str += " " + op.param3.String()
	}
	if op.param4 != nil {
		str += " " + op.param4.String()
	}
	return str
}

func NewReturnOp() *Operation {
	return &Operation{kind: OP_RETURN}
}

func NewMoveOp(dest, src *Object) *Operation {
	return &Operation{kind: OP_MOVE, param1: dest, param2: src}
}

func NewDefLabelOp(label *Object) *Operation {
	return &Operation{kind: OP_DEF_LABEL, param1: label}
}

func NewPushOp(src *Object) *Operation {
	return &Operation{kind: OP_PUSH, param1: src}
}
func NewPopOp(dest *Object) *Operation {
	return &Operation{kind: OP_POP, param1: dest}
}

func NewAddOp(dest, src *Object) *Operation {
	return &Operation{kind: OP_ADD, param1: dest, param2: src}
}
func NewSubOp(dest, src *Object) *Operation {
	return &Operation{kind: OP_SUB, param1: dest, param2: src}
}

func NewCallOp(label *Object) *Operation {
	return &Operation{kind: OP_CALL, param1: label}
}

package runtime

type OperationKind int

const (
	OP_ILLEGAL OperationKind = iota
	OP_EXIT
	OP_MOVE
	OP_ADD
	OP_SUB
)

var opKinds = [...]string{
	OP_ILLEGAL: "ILLEGAL",
	OP_EXIT:    "EXIT",
	OP_MOVE:    "MOVE",
	OP_ADD:     "ADD",
	OP_SUB:     "SUB",
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

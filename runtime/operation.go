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

type Program []*Operation

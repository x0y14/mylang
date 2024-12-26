package runtime

type RegisterKind int

const (
	REG_RETURN_ADDRESS RegisterKind = iota
	REG_PROGRAM_COUNTER
	REG_STATUS
	REG_BOOL_FLAG
	REG_GENERAL_1
	REG_GENERAL_2
	REG_TEMP_1
)

var regKinds = [...]string{
	REG_RETURN_ADDRESS:  "RETURN_ADDRESS",
	REG_PROGRAM_COUNTER: "PROGRAM_COUNTER",
	REG_STATUS:          "STATUS",
	REG_BOOL_FLAG:       "BOOL_FLAG",
	REG_GENERAL_1:       "GENERAL_1",
	REG_GENERAL_2:       "GENERAL_2",
	REG_TEMP_1:          "TEMP_1",
}

func (regKind RegisterKind) String() string {
	return regKinds[regKind]
}

type Register []*Object

func NewRegister() Register {
	reg := make([]*Object, len(regKinds))
	return reg
}

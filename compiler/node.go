package compiler

type Syntax int

const (
	ST_ILLEGAL Syntax = iota

	ST_DEFINE_FUNCTION

	ST_PRIMITIVE
	ST_INTEGER

	ST_BLOCK
	ST_RETURN
)

var stKinds = [...]string{
	ST_ILLEGAL:         "ILLEGAL",
	ST_DEFINE_FUNCTION: "DEFINE_FUNCTION",
	ST_PRIMITIVE:       "PRIMITIVE",
	ST_INTEGER:         "INTEGER",
	ST_BLOCK:           "BLOCK",
	ST_RETURN:          "RETURN",
}

func (st Syntax) String() string {
	return stKinds[st]
}

type Node struct {
	kind Syntax
	leaf *Token
	lhs  *Node // 1個しか要素がないならLHSを使う
	rhs  *Node
	next *Node
}

func (n *Node) String() string {
	return ""
}
package compiler

import (
	"fmt"
	"mylang/runtime"
)

var curt *Node

func nextNode() error {
	if curt.next == nil {
		return fmt.Errorf("end of node")
	}
	curt = curt.next
	return nil
}
func genPrimitive(nd *Node) (*runtime.Object, error) {
	switch primValue := nd.lhs; primValue.kind {
	case ST_INTEGER:
		i, err := primValue.leaf.GetInt()
		if err != nil {
			return nil, err
		}
		return runtime.NewObject(i), nil
	default:
		return nil, fmt.Errorf("genPrimitive: unsupported value: %s", primValue.String())
	}
}

func genReturn(nd *Node) (runtime.Program, error) {
	prog := runtime.Program{}
	switch retValue := nd.lhs; {
	case retValue == nil:
		prog = append(prog, runtime.Program{
			runtime.NewReturnOp(),
		}...)
	case retValue.kind == ST_PRIMITIVE:
		retObj, err := genPrimitive(retValue)
		if err != nil {
			return nil, err
		}
		prog = append(prog, runtime.Program{
			runtime.NewMoveOp(runtime.REG_STATUS, retObj),
			runtime.NewReturnOp(),
		}...)
	default:
		return nil, fmt.Errorf("genReturn: unsupported value: %s", retValue.String())
	}
	return prog, nil
}

func Generate(node *Node) (runtime.Program, error) {
	curt = &Node{next: node} // dummy

	program := runtime.Program{}
	for {
		if err := nextNode(); err != nil { // end of node
			break
		}
		switch curt.kind {
		case ST_RETURN:
			prog, err := genReturn(curt)
			if err != nil {
				return nil, err
			}
			program = append(program, prog...)
		default:
			return nil, fmt.Errorf("unsupported syntax: %v", curt.kind.String())
		}
	}
	return program, nil
}

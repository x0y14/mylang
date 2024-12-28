package compiler

import (
	"fmt"
	"mylang/runtime"
	"slices"
)

var curt *Node
var lc *LabelCollector

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
			runtime.NewMoveOp(runtime.NewRegisterObject(runtime.REG_STATUS), retObj),
			runtime.NewReturnOp(),
		}...)
	default:
		return nil, fmt.Errorf("genReturn: unsupported value: %s", retValue.String())
	}
	return prog, nil
}

func genBlock(nd *Node) (runtime.Program, error) {
	backup := *curt
	prog, err := Generate(nd.lhs)
	if err != nil {
		return nil, err
	}
	curt = &backup
	return prog, nil
}

func genIdent(nd *Node) (int, error) {
	id, err := nd.leaf.GetIdent()
	if err != nil {
		return 0, err
	}
	no, ok := lc.Get(id)
	if ok {
		return no, nil
	}
	no, err = lc.Set(id)
	if err != nil {
		return 0, err
	}
	return no, nil
}

func genFunctionArguments(nd *Node, fnNameLabel int) runtime.Program {
	// argumentsとして幾つpushされるのか
	prog := runtime.Program{}
	count := 0
	argHead := nd.lhs
	args := []Node{}
	for {
		if argHead == nil {
			break
		}
		count++
		args = append(args, *argHead)
		argHead = argHead.next
	}
	// 引数なし
	if count == 0 {
		return prog
	}
	// 逆にする
	slices.Reverse(args)
	for _, _ = range args {
		prog = append(prog, runtime.Program{
			runtime.NewPopOp(runtime.NewRegisterObject(runtime.REG_TEMP_1)),
			runtime.NewMoveOp(runtime.NewReferenceObject(fnNameLabel+count), runtime.NewRegisterObject(runtime.REG_TEMP_1)),
		}...)
		count--
	}

	return prog
}

func analyzeFunctionHeader(nd *Node) (int, runtime.Program, error) {
	fnNameLabel, err := genIdent(nd.lhs)
	if err != nil {
		return 0, nil, err
	}
	fnArgsProg := genFunctionArguments(nd.rhs, fnNameLabel)
	return fnNameLabel, fnArgsProg, nil
}

func analyzeFunctionDeclaration(nd *Node) (int, runtime.Program, int, error) {
	// fnHeader
	fnNameLabel, fnArgsProg, err := analyzeFunctionHeader(nd.lhs)
	if err != nil {
		return 0, nil, 0, err
	}
	// fnReturns
	//fnReturnsCount := analyzeFunctionReturns(nd.rhs)
	return fnNameLabel, fnArgsProg, 0, nil
}

func genDefineFunction(nd *Node) (runtime.Program, error) {
	nameLabel, argsProg, _, err := analyzeFunctionDeclaration(nd.lhs)
	if err != nil {
		return nil, err
	}
	blockProg, err := genBlock(nd.rhs)
	if err != nil {
		return nil, err
	}

	prog := runtime.Program{
		runtime.NewDefLabelOp(runtime.NewLabelObject(nameLabel)),
	}
	prog = append(prog, argsProg...)
	prog = append(prog, blockProg...)
	//prog = append(prog, runtime.NewReturnOp())

	return prog, nil
}

func Generate(node *Node) (runtime.Program, error) {
	curt = &Node{next: node} // dummy
	lc = NewLabelCollector()
	lc.Init()

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
		case ST_DEFINE_FUNCTION:
			prog, err := genDefineFunction(curt)
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

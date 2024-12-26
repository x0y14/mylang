package runtime

type Program []*Operation

func Export(prog Program) string {
	str := ""

	for i, op := range prog {
		str += op.String()
		if i != len(prog)-1 { // 最後の行でなかったら
			str += "\n"
		}
	}

	return str
}

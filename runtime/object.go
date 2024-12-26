package runtime

import (
	"fmt"
	"log"
	"strconv"
)

type ObjectKind int

const (
	OBJ_INVALID ObjectKind = iota
	OBJ_NULL
	OBJ_INT
	OBJ_CHAR
	OBJ_BOOL
	OBJ_LIST
	OBJ_REFERENCE
	OBJ_LABEL
)

var objectKinds = [...]string{
	OBJ_INVALID:   "INVALID",
	OBJ_NULL:      "NULL",
	OBJ_INT:       "INT",
	OBJ_CHAR:      "CHAR",
	OBJ_BOOL:      "BOOL",
	OBJ_LIST:      "LIST",
	OBJ_REFERENCE: "REFERENCE",
	OBJ_LABEL:     "LABEL",
}

func (objKind ObjectKind) String() string {
	return objectKinds[objKind]
}

func NewObject[T int | rune | bool](data T) *Object {
	obj := Object{}

	switch any(data).(type) {
	case int:
		obj.kind = OBJ_INT
		obj.data = any(data).(int)
	case rune:
		obj.kind = OBJ_CHAR
		obj.data = int(any(data).(rune))
	case bool:
		obj.kind = OBJ_BOOL
		tORf := any(data).(bool)
		if tORf == true {
			obj.data = 1
		} else {
			obj.data = 0
		}
	default:
		log.Fatalf("unsupported object: data: %v", data)
	}
	return &obj
}

func NewNullObject() *Object {
	return &Object{kind: OBJ_NULL, data: 0}
}

func NewListObject(size int) *Object {
	return &Object{kind: OBJ_LIST, data: size}
}

func NewReferenceObject(reg RegisterKind) *Object {
	return &Object{kind: OBJ_REFERENCE, data: int(reg)}
}

func NewLabelObject(destPC int) *Object {
	return &Object{kind: OBJ_LABEL, data: destPC}
}

type Object struct {
	kind ObjectKind
	data int
}

func (o *Object) String() string {
	switch o.kind {
	case OBJ_INVALID:
		return "invalid"
	case OBJ_NULL:
		return "null"
	case OBJ_INT:
		return strconv.Itoa(o.data)
	case OBJ_CHAR:
		return string(rune(o.data))
	case OBJ_BOOL:
		if o.data == 1 {
			return "true"
		} else {
			return "false"
		}
	case OBJ_LIST:
		return fmt.Sprintf("list(%s)", strconv.Itoa(o.data))
	case OBJ_REFERENCE:
		return fmt.Sprintf("reference(%s)", RegisterKind(o.data).String())
	case OBJ_LABEL:
		return fmt.Sprintf("label(%d)", o.data)

	default:
		log.Fatalf("unsupported object kind: %s", o.kind)
	}
	return ""
}

func (o *Object) GetKind() ObjectKind {
	return o.kind
}

func (o *Object) Clone() *Object {
	newObj := Object{kind: o.kind, data: o.data}
	return &newObj
}

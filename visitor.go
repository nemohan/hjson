package hjson

import (
	"bytes"
	"fmt"
)

//visitor pattern
type Walker interface {
	walkObject(JObject)
	walkArray(JArray)
	walkString(JString)
	walkNumber(JNumber)
	walkBool(JBool)
	walkNull(JNull)
}

func (o JObject) accept(w Walker) {
	w.walkObject(o)
}
func (a JArray) accept(w Walker) {
	w.walkArray(a)
}
func (s JString) accept(w Walker) {
	w.walkString(s)
}
func (o JNumber) accept(w Walker) {
	w.walkNumber(o)
}
func (o JNull) accept(w Walker) {
	w.walkNull(o)
}
func (b JBool) accept(w Walker) {
	w.walkBool(b)
}

type nodeVisitor struct {
	buf    *bytes.Buffer
	indent int
}

func newNodeVisitor() *nodeVisitor {
	return &nodeVisitor{
		buf: bytes.NewBufferString(""),
	}
}
func (n *nodeVisitor) walkObject(obj JObject) {
	n.buf.WriteString("{")
	num := len(obj.values)
	i := 0
	for k, v := range obj.values {
		n.buf.WriteString(k)
		n.buf.WriteString(":")
		v.accept(n)
		if i == num-1 {
			break
		}
		n.buf.WriteString(",")
		i++
	}
	n.buf.WriteString("}")
}
func (n *nodeVisitor) walkArray(array JArray) {
	n.buf.WriteString("[")
	num := len(array.elements)
	for i, v := range array.elements {
		v.accept(n)
		if i == num-1 {
			break
		}
		n.buf.WriteString(",")
	}
	n.buf.WriteString("]")
}

func (n *nodeVisitor) walkString(str JString) {
	n.buf.WriteString(string(str))
}
func (n *nodeVisitor) walkNumber(number JNumber) {
	n.buf.WriteString(fmt.Sprint(number))
}
func (n *nodeVisitor) walkBool(v JBool) {
	n.buf.WriteString(fmt.Sprintf("%t", v))
}
func (n *nodeVisitor) walkNull(v JNull) {
	n.buf.WriteString(v.String())
}

package node

import "github.com/susruth/immutable/field"

type Node struct {
	Fields []field.Field
}

func New(fields ...field.Field) Node {
	return Node{
		fields,
	}
}

func (node *Node) AddField(f field.Field) {
	node.Fields = append(node.Fields, f)
}

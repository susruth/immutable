package list

import (
	"github.com/susruth/immutable/field"
	"github.com/susruth/immutable/node"
)

type List struct {
	CurrentVersion uint64
	Nodes          []node.Node
}

func New(vals ...interface{}) *List {
	nodes := make([]node.Node, len(vals))
	for i, val := range vals {
		nodes[i] = node.New(field.New(0, val))
	}
	return &List{
		Nodes: nodes,
	}
}

func (list *List) Append(value interface{}) {
	list.CurrentVersion++
	list.Nodes = append(list.Nodes, node.New(field.New(list.CurrentVersion, value)))
}

func (list *List) Update(index int, value interface{}) {
	list.CurrentVersion++
	list.Nodes[index].AddField(field.New(list.CurrentVersion, value))
}

func (list *List) Delete(index int) {
	list.Update(index, nil)
}

func (list *List) Values(version uint64) []interface{} {
	vals := []interface{}{}
	for _, node := range list.Nodes {
		field, stop := field.Search(node.Fields, version)
		if stop {
			return vals
		}
		if field.Value != nil {
			vals = append(vals, field.Value)
		}
	}
	return vals
}

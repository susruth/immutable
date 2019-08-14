package record

import (
	"github.com/susruth/immutable/field"
	"github.com/susruth/immutable/node"
)

type Record struct {
	CurrentVersion uint64
	Nodes          map[string]node.Node
}

func New(vals map[string]interface{}) *Record {
	nodes := make(map[string]node.Node, len(vals))
	for i, val := range vals {
		nodes[i] = node.New(field.New(0, val))
	}
	return &Record{
		Nodes: nodes,
	}
}

func (record *Record) Insert(key string, value interface{}) {
	record.CurrentVersion++
	record.Nodes[key] = node.New(field.New(record.CurrentVersion, value))
}

func (record *Record) Update(key string, value interface{}) {
	record.CurrentVersion++
	node := record.Nodes[key]
	node.AddField(field.New(record.CurrentVersion, value))
	record.Nodes[key] = node
}

func (record *Record) Delete(key string) {
	record.Update(key, nil)
}

func (record *Record) Values(version uint64) map[string]interface{} {
	vals := map[string]interface{}{}
	for key, node := range record.Nodes {
		field, stop := field.Search(node.Fields, version)
		if stop {
			return vals
		}
		if field.Value != nil {
			vals[key] = field.Value
		}
	}
	return vals
}

package node

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/susruth/immutable/field"
	"golang.org/x/crypto/sha3"
)

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

func (node *Node) LatestVersion() uint64 {
	return node.Fields[len(node.Fields)-1].Version
}

type Nodes map[string]Node

func (n Nodes) String() string {
	data, err := json.Marshal(n)
	if err != nil {
		panic(fmt.Sprintf("failed to marhsal nodes: %v", err))
	}
	hash := sha3.Sum256(data)
	return base64.StdEncoding.EncodeToString(hash[:])
}

func (n Nodes) Clone() Nodes {
	nodes := Nodes{}
	for key, node := range n {
		nodes[key] = node
	}
	return nodes
}

func (n Nodes) Insert(key string, f field.Field) error {
	node, ok := n[key]
	if ok {
		node.AddField(f)
	} else {
		node = New(f)
	}
	n[key] = node
	return nil
}

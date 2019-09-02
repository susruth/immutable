package record

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/renproject/kv"
	"github.com/susruth/immutable/field"
	"github.com/susruth/immutable/node"
)

// Record is a Persistent Map
type Record struct {
	CurrentVersion uint64
	Nodes          node.Nodes

	mu *sync.RWMutex
}

type Snapshot map[string]interface{}

func New() *Record {
	return &Record{
		Nodes: node.Nodes{},

		mu: new(sync.RWMutex),
	}
}

func (record *Record) Update(key string, value interface{}) (uint64, error) {
	record.mu.Lock()
	defer record.mu.Unlock()
	if err := record.Nodes.Insert(key, field.New(record.CurrentVersion+1, value)); err != nil {
		return record.CurrentVersion, err
	}
	record.CurrentVersion++
	return record.CurrentVersion, nil
}

func (record *Record) Get(ver uint64, key string) (interface{}, bool) {
	record.mu.RLock()
	defer record.mu.RUnlock()
	f, ok := field.Search(record.Nodes[key].Fields, ver)
	return f.Value, ok
}

func (record *Record) Delete(key string) (uint64, error) {
	record.mu.Lock()
	defer record.mu.Unlock()
	if _, ok := record.Nodes[key]; !ok {
		return record.CurrentVersion, fmt.Errorf("key not found")
	}
	if err := record.Nodes.Insert(key, field.New(record.CurrentVersion+1, nil)); err != nil {
		return record.CurrentVersion, err
	}
	record.CurrentVersion++
	return record.CurrentVersion, nil
}

func (record *Record) BatchUpdate(updates map[string]interface{}) (uint64, error) {
	record.mu.RLock()
	nodes := record.Nodes.Clone()
	record.mu.RUnlock()
	for key, val := range updates {
		if err := nodes.Insert(key, field.New(record.CurrentVersion+1, val)); err != nil {
			return record.CurrentVersion, err
		}
	}

	record.mu.Lock()
	defer record.mu.Unlock()
	record.Nodes = nodes

	record.CurrentVersion++
	return record.CurrentVersion, nil
}

func (record *Record) Snapshot(version uint64) Snapshot {
	snapshot := Snapshot{}
	for key, node := range record.Nodes {
		field, stop := field.Search(node.Fields, version)
		if stop {
			return snapshot
		}
		if field.Value != nil {
			snapshot[key] = field.Value
		}
	}
	return snapshot
}

func (record *Record) MarshalBinary() ([]byte, error) {
	return json.Marshal(record)
}

func (record *Record) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, record)
}

func (record *Record) StoreSingle(key string, table kv.Table) error {
	return table.Insert(key, record.Nodes[key])
}

func (record *Record) Store(table kv.Table) error {
	for key, node := range record.Nodes {
		if err := table.Insert(key, node); err != nil {
			return err
		}
	}
	return nil
}

func (record *Record) Restore(table kv.Table) error {
	iter, err := table.Iterator()
	if err != nil {
		return err
	}

	size, err := table.Size()
	if err != nil {
		return err
	}

	var currVersion uint64
	nodes := make(node.Nodes, size)
	for iter.Next() {
		key, err := iter.Key()
		if err != nil {
			return err
		}
		node := node.Node{}
		if err := iter.Value(&node); err != nil {
			return err
		}
		nodes[key] = node

		if currVersion < node.LatestVersion() {
			currVersion = node.LatestVersion()
		}
	}

	record.Nodes = nodes
	record.CurrentVersion = currVersion
	return nil
}

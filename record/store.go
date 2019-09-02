package record

import (
	"sync"

	"github.com/renproject/kv"
	"github.com/renproject/phi"
	"github.com/susruth/immutable/node"
)

type Store struct {
	storage kv.Table
	record  Record
}

func NewStore(table kv.Table) *Store {
	iter := table.Iterator()
	size, err := table.Size()
	if err != nil {
		panic(err)
	}

	var currVersion uint64
	nodes := make(node.Nodes, size)
	for iter.Next() {
		key, err := iter.Key()
		if err != nil {
			panic(err)
		}
		node := node.Node{}
		if err := iter.Value(&node); err != nil {
			panic(err)
		}
		nodes[key] = node

		if currVersion < node.LatestVersion() {
			currVersion = node.LatestVersion()
		}
	}

	return &Store{
		storage: table,
		record:  Record{currVersion, nodes, new(sync.RWMutex)},
	}
}

func (store *Store) Handle(_ phi.Task, msg phi.Message) {
	switch msg := msg.(type) {
	case GetRequest:
		store.handleGet(msg)
	case SnapshotRequest:
		store.handleSnapshot(msg)
	case BatchUpdateRequest:
		store.handleBatchUpdateRequest(msg)
	case UpdateRequest:
		store.handleUpdateRequest(msg)
	}
}

func (store *Store) handleGet(msg GetRequest) {
	val, ok := store.record.Get(msg.Version, msg.Key)
	msg.Responder <- GetResponse{Element: Element{Key: msg.Key, Value: val}, Ok: ok}
}

func (store *Store) handleSnapshot(msg SnapshotRequest) {
	msg.Responder <- store.record.Snapshot(msg.Version)
}

func (store *Store) handleBatchUpdateRequest(msg BatchUpdateRequest) {
	updates := map[string]interface{}{}
	for _, element := range msg.Elements {
		updates[element.Key] = element.Value
	}
	ver, err := store.record.BatchUpdate(updates)
	if err != nil {
		msg.Responder <- UpdateResponse{Version: ver, Err: err}
		return
	}
	err = store.record.Store(store.storage)
	msg.Responder <- UpdateResponse{Version: ver, Err: err}
}

func (store *Store) handleUpdateRequest(msg UpdateRequest) {
	ver, err := store.record.Update(msg.Element.Key, msg.Element.Value)
	if err != nil {
		msg.Responder <- UpdateResponse{Version: ver, Err: err}
		return
	}
	err = store.record.StoreSingle(msg.Element.Key, store.storage)
	msg.Responder <- UpdateResponse{Version: ver, Err: err}
}

func (store *Store) handleDeleteRequest(msg UpdateRequest) {
	ver, err := store.record.Delete(msg.Element.Key)
	if err != nil {
		msg.Responder <- UpdateResponse{Version: ver, Err: err}
		return
	}
	err = store.record.StoreSingle(msg.Element.Key, store.storage)
	msg.Responder <- UpdateResponse{Version: ver, Err: err}
}

type GetRequest struct {
	Version uint64
	Key     string

	Responder chan<- GetResponse
}

func (GetRequest) IsMessage() {
}

type SnapshotRequest struct {
	Version uint64

	Responder chan<- Snapshot
}

func (SnapshotRequest) IsMessage() {
}

type BatchUpdateRequest struct {
	Elements []Element

	Responder chan<- UpdateResponse
}

func (BatchUpdateRequest) IsMessage() {
}

type UpdateRequest struct {
	Element

	Responder chan<- UpdateResponse
}

func (UpdateRequest) IsMessage() {
}

type DeleteRequest struct {
	Key string

	Responder chan<- UpdateResponse
}

func (DeleteRequest) IsMessage() {
}

type Element struct {
	Key   string
	Value interface{}
}

type GetResponse struct {
	Element
	Ok bool
}

type UpdateResponse struct {
	Version uint64
	Err     error
}

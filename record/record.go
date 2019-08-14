package record

type Field struct {
	Version uint64
	Value   interface{}
}

type Node struct {
	Fields []Field
}

type Record struct {
	CurrentVersion uint64
	Nodes          map[string]Node
}

func New(vals map[string]interface{}) *Record {
	nodes := make(map[string]Node, len(vals))
	for i, val := range vals {
		nodes[i] = Node{[]Field{{Value: val}}}
	}
	return &Record{
		Nodes: nodes,
	}
}

func (record *Record) Insert(key string, value interface{}) {
	record.CurrentVersion++
	record.Nodes[key] = Node{Fields: []Field{{record.CurrentVersion, value}}}
}

func (record *Record) Update(key string, value interface{}) {
	record.CurrentVersion++
	node := record.Nodes[key]
	node.Fields = append(record.Nodes[key].Fields, Field{record.CurrentVersion, value})
	record.Nodes[key] = node
}

func (record *Record) Delete(key string) {
	record.Update(key, nil)
}

func (record *Record) Values(version uint64) map[string]interface{} {
	vals := map[string]interface{}{}
	for key, node := range record.Nodes {
		field, stop := getField(version, node.Fields)
		if stop {
			return vals
		}
		if field.Value != nil {
			vals[key] = field.Value
		}
	}
	return vals
}

func getField(version uint64, fields []Field) (Field, bool) {
	if fields[len(fields)-1].Version < version {
		return fields[len(fields)-1], false
	}
	if fields[0].Version > version {
		return Field{}, true
	}
	field := fields[len(fields)/2]
	if field.Version > version {
		return getField(version, fields[:len(fields)/2])
	}
	if field.Version < version {
		return getField(version, fields[len(fields)/2:])
	}
	return field, false
}

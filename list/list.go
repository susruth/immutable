package list

type Field struct {
	Version uint64
	Value   interface{}
}

type Node struct {
	Fields []Field
}

type List struct {
	CurrentVersion uint64
	Nodes          []Node
}

func New(vals ...interface{}) *List {
	nodes := make([]Node, len(vals))
	for i, val := range vals {
		nodes[i] = Node{[]Field{{Value: val}}}
	}
	return &List{
		Nodes: nodes,
	}
}

func (list *List) Append(value interface{}) {
	list.CurrentVersion++
	list.Nodes = append(list.Nodes, Node{Fields: []Field{{list.CurrentVersion, value}}})
}

func (list *List) Update(index int, value interface{}) {
	list.CurrentVersion++
	list.Nodes[index].Fields = append(list.Nodes[index].Fields, Field{list.CurrentVersion, value})
}

func (list *List) Remove(index int) {
	list.CurrentVersion++
	list.Nodes[index].Fields = append(list.Nodes[index].Fields, Field{list.CurrentVersion, nil})
}

func (list *List) Values(version uint64) []interface{} {
	vals := []interface{}{}
	for _, node := range list.Nodes {
		field, stop := getField(version, node.Fields)
		if stop {
			return vals
		}
		if field.Value != nil {
			vals = append(vals, field.Value)
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

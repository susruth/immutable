package field

type Field struct {
	Version uint64
	Value   interface{}
}

func New(ver uint64, val interface{}) Field {
	return Field{ver, val}
}

func Search(fields []Field, version uint64) (Field, bool) {
	if fields[len(fields)-1].Version < version {
		return fields[len(fields)-1], false
	}
	if fields[0].Version > version {
		return Field{}, true
	}
	field := fields[len(fields)/2]
	if field.Version > version {
		return Search(fields[:len(fields)/2], version)
	}
	if field.Version < version {
		return Search(fields[len(fields)/2:], version)
	}
	return field, false
}

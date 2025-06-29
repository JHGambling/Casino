package protocol

type SubChangedRecord struct {
	Operation  string
	TableID    string
	ResourceID interface{}
	Record     interface{}
}

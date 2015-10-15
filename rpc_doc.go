package dbdb

type RPCDoc struct {
	Store  string
	DocId  uint64
	DocVal map[string]interface{}
	DocIds []uint64
}

package dbdb

type QuerySet struct {
	Docs  []*Doc
	Count int
}

type RPCDoc struct {
	Store  string
	DocId  uint64
	DocVal map[string]interface{}
	DocIds []uint64
	Comps  []QueryComp
}

package dbdb

type RPCDoc struct {
	Store  string
	DocId  float64
	DocVal map[string]interface{}
	DocIds []float64
	Comps  []QueryComp
}

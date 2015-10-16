package dbdb

import (
	"encoding/gob"
	"log"
	"net/rpc"
)

func init() {
	gob.Register([]interface{}(nil))
}

type RPCClient struct {
	conn *rpc.Client
}

func NewRPCClient(dsn string) *RPCClient {
	conn, err := rpc.Dial("tcp", dsn)
	if err != nil {
		log.Fatalf("NewRPCClient() -> rpc.Dial() -> %v\n", err)
	}
	return &RPCClient{
		conn: conn,
	}
}

func (rpcc *RPCClient) GetAllStoreStats() ([]*StoreStat, error) {
	var stats []*StoreStat
	err := rpcc.conn.Call("RPCServer.GetAllStoreStats", struct{}{}, &stats)
	return stats, err
}

func (rpcc *RPCClient) GetStoreStat(store string) (*StoreStat, error) {
	var stat *StoreStat
	err := rpcc.conn.Call("RPCServer.GetStoreStat", store, &stat)
	return stat, err
}

func (rpcc *RPCClient) AddStore(store string) (bool, error) {
	var ok bool
	err := rpcc.conn.Call("RPCServer.AddStore", store, &ok)
	return ok, err
}

func (rpcc *RPCClient) HasStore(store string) (bool, error) {
	var ok bool
	err := rpcc.conn.Call("RPCServer.HasStore", store, &ok)
	return ok, err
}

func (rpcc *RPCClient) DelStore(store string) (bool, error) {
	var ok bool
	err := rpcc.conn.Call("RPCServer.DelStore", store, &ok)
	return ok, err
}

func (rpcc *RPCClient) Add(rpcdoc RPCDoc) (uint64, error) {
	var docid uint64
	err := rpcc.conn.Call("RPCServer.Add", rpcdoc, &docid)
	return docid, err
}

func (rpcc *RPCClient) Set(rpcdoc RPCDoc) (bool, error) {
	var ok bool
	err := rpcc.conn.Call("RPCServer.Set", rpcdoc, &ok)
	return ok, err
}

func (rpcc *RPCClient) Get(rpcdoc RPCDoc) (*Doc, error) {
	var doc *Doc
	err := rpcc.conn.Call("RPCServer.Get", rpcdoc, &doc)
	return doc, err
}

func (rpcc *RPCClient) GetAll(name string, id ...uint64) ([]*Doc, error) {
	rpcdoc := RPCDoc{
		Store:  name,
		DocIds: id,
	}
	var docs []*Doc
	err := rpcc.conn.Call("RPCServer.GetAll", rpcdoc, &docs)
	return docs, err
}

func (rpcc *RPCClient) Del(rpcdoc RPCDoc) (bool, error) {
	var ok bool
	err := rpcc.conn.Call("RPCServer.Del", rpcdoc, &ok)
	return ok, err
}

func (rpcc *RPCClient) Close() error {
	return rpcc.conn.Close()
}

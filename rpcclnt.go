package dbdb

import (
	"log"
	"net/rpc"
)

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

func (rpcc *RPCClient) AddStore(store string) (*bool, error) {
	var ok *bool
	err := rpcc.conn.Call("RPCServer.AddStore", store, &ok)
	return ok, err
}

func (rpcc *RPCClient) GetStore(store string) (*Store, error) {
	var st *Store
	err := rpcc.conn.Call("RPCServer.GetStore", store, &st)
	return st, err
}

func (rpcc *RPCClient) DelStore(store string) (*bool, error) {
	var ok *bool
	err := rpcc.conn.Call("RPCServer.DelStore", store, &ok)
	return ok, err
}

func (rpcc *RPCClient) Add(docmsg DocMsg) (*uint64, error) {
	var docid *uint64
	err := rpcc.conn.Call("RPCServer.Add", docmsg, &docid)
	return docid, err
}

func (rpcc *RPCClient) Set(docmsg DocMsg) (*bool, error) {
	var ok *bool
	err := rpcc.conn.Call("RPCServer.Set", docmsg, &ok)
	return ok, err
}

func (rpcc *RPCClient) Get(docmsg DocMsg) (*Doc, error) {
	var doc *Doc
	err := rpcc.conn.Call("RPCServer.Get", docmsg, &doc)
	return doc, err
}

func (rpcc *RPCClient) Del(docmsg DocMsg) (*bool, error) {
	var ok *bool
	err := rpcc.conn.Call("RPCServer.Del", docmsg, &ok)
	return ok, err
}

func (rpcc *RPCClient) Close() error {
	return rpcc.conn.Close()
}

package dbdb

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

type DocMsg struct {
	Store  string
	DocId  uint64
	DocVal interface{}
}

type RPCServer struct {
	ds *DataStore
}

func NewRPCServer(ds *DataStore) *RPCServer {
	return &RPCServer{
		ds: ds,
	}
}

func (rpcs *RPCServer) AddStore(store string, resp *bool) error {
	rpcs.ds.AddStore(store)
	*resp = true
	return nil
}

func (rpcs *RPCServer) GetStore(store string, resp *Store) error {
	st, ok := rpcs.ds.GetStore(store)
	if !ok {
		return fmt.Errorf("store (%s) not found\n", store)
	}
	*resp = *st
	return nil
}

func (rpcs *RPCServer) DelStore(store string, resp *bool) error {
	rpcs.ds.DelStore(store)
	*resp = true
	return nil
}

func (rpcs *RPCServer) Add(docmsg DocMsg, resp *uint64) error {
	docid := rpcs.ds.Add(docmsg.Store, docmsg.DocVal)
	if docid == 0 {
		return fmt.Errorf("error adding document (%v)\n", docmsg.DocVal)
	}
	*resp = docid
	return nil
}

func (rpcs *RPCServer) Set(docmsg DocMsg, resp *bool) error {
	rpcs.ds.Set(docmsg.Store, docmsg.DocId, docmsg.DocVal)
	*resp = true
	return nil
}

func (rpcs *RPCServer) Get(docmsg DocMsg, resp *Doc) error {
	doc := rpcs.ds.Get(docmsg.Store, docmsg.DocId)
	if doc == nil {
		return fmt.Errorf("error getting document (%v)\n", docmsg.DocId)
	}
	*resp = *doc
	return nil
}

func (rpcs *RPCServer) Del(docmsg DocMsg, resp *bool) error {
	rpcs.ds.Del(docmsg.Store, docmsg.DocId)
	*resp = true
	return nil
}

func (rpcs *RPCServer) ListenAndServe(host string) {
	rpc.Register(rpcs)
	addr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		log.Fatalf("RPCServer.Listen() -> net.ResolveTCPAddr() -> %v\n", err)
	}
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalf("RPCServer.Listen() -> net.ListenTCP() -> %v\n", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("error accepting conn: %v\n", err)
			continue
		}
		rpc.ServeConn(conn)
	}
}

package dbdb

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"net/rpc"
)

func Serve(ds *DataStore, port string) {

	// register some types for gob...
	gob.Register([]interface{}(nil))
	gob.Register(map[string]interface{}(nil))

	srv := NewRPCServer(ds)
	if err := rpc.Register(srv); err != nil {
		log.Fatal(err)
	}
	addr, err := net.ResolveTCPAddr("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening for connections on %s...", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("error accepting conn: %v\n", err)
			continue
		}
		rpc.ServeConn(conn)
	}
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

func (rpcs *RPCServer) HasStore(store string, resp *bool) error {
	ok := rpcs.ds.HasStore(store)
	if !ok {
		return fmt.Errorf("store (%s) not found\n", store)
	}
	*resp = ok
	return nil
}

/*
func (rpcs *RPCServer) GetStore(store string, resp *Store) error {
	st, ok := rpcs.ds.GetStore(store)
	if !ok {
		return fmt.Errorf("store (%s) not found\n", store)
	}
	*resp = *st
	return nil
}
*/

func (rpcs *RPCServer) DelStore(store string, resp *bool) error {
	rpcs.ds.DelStore(store)
	*resp = true
	return nil
}

func (rpcs *RPCServer) Add(rpcdoc RPCDoc, resp *uint64) error {
	docid := rpcs.ds.Add(rpcdoc.Store, rpcdoc.DocVal)
	if docid == 0 {
		return fmt.Errorf("error adding document (%v)\n", rpcdoc.DocVal)
	}
	*resp = docid
	return nil
}

func (rpcs *RPCServer) Set(rpcdoc RPCDoc, resp *bool) error {
	rpcs.ds.Set(rpcdoc.Store, rpcdoc.DocId, rpcdoc.DocVal)
	*resp = true
	return nil
}

func (rpcs *RPCServer) Get(rpcdoc RPCDoc, resp *Doc) error {
	doc := rpcs.ds.Get(rpcdoc.Store, rpcdoc.DocId)
	if doc == nil {
		return fmt.Errorf("error getting document (%v)\n", rpcdoc.DocId)
	}
	*resp = *doc
	return nil
}

func (rpcs *RPCServer) Del(rpcdoc RPCDoc, resp *bool) error {
	rpcs.ds.Del(rpcdoc.Store, rpcdoc.DocId)
	*resp = true
	return nil
}

/*
func (rpcs *RPCServer) ListenAndServe(host string) {
	if err := rpc.Register(rpcs); err != nil {
		log.Fatalf("RPCServer.ListenAndServe() -> rpc.Register() -> %v\n", err)
	}
	addr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		log.Fatalf("RPCServer.ListenAndServe() -> net.ResolveTCPAddr() -> %v\n", err)
	}
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalf("RPCServer.ListenAndServe() -> net.ListenTCP() -> %v\n", err)
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
*/

package dbdb

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"time"
)

func init() {
	gob.Register([]interface{}(nil))
}

var authtoken string

func Serve(ds *DataStore, port string, token string) {
	srv := NewServer(ds)
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
	authtoken = fmt.Sprintf("authtoken:%s\n", token)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("error accepting conn: %v\n", err)
			continue
		}
		if CanHandle(conn, authtoken) {
			go HandleConn(conn)
		}
	}
}

func CanHandle(conn net.Conn, token string) bool {
	auth, _ := bufio.NewReader(conn).ReadString('\n')
	if auth != token {
		conn.SetDeadline(time.Now())
		conn.Close()
		return false
	}
	if _, err := conn.Write([]byte{'1', '\n'}); err != nil {
		log.Fatal(err)
	}
	return true
}

func HandleConn(conn net.Conn) {
	// enter infinite loop and handle each request
	// on a tcp level so we can handle deadlines, etc...
	rpc.ServeConn(conn)
}

type Server struct {
	ds *DataStore
}

func NewServer(ds *DataStore) *Server {
	return &Server{
		ds: ds,
	}
}

//EXAMPLE...
//func (*Receiver) DoNothing(_, _ *struct{}), error

func (s *Server) GetAllStoreStats(_ struct{}, resp *[]*StoreStat) error {
	stats := s.ds.GetAllStoreStats()
	*resp = stats
	return nil
}

func (s *Server) GetStoreStat(store string, resp *StoreStat) error {
	stat := s.ds.GetStoreStat(store)
	*resp = *stat
	return nil
}

func (s *Server) AddStore(store string, resp *bool) error {
	s.ds.AddStore(store)
	*resp = true
	return nil
}

func (s *Server) HasStore(store string, resp *bool) error {
	ok := s.ds.HasStore(store)
	if !ok {
		return fmt.Errorf("store (%s) not found\n", store)
	}
	*resp = ok
	return nil
}

func (s *Server) DelStore(store string, resp *bool) error {
	s.ds.DelStore(store)
	*resp = true
	return nil
}

func (s *Server) Add(rpcdoc RPCDoc, resp *uint64) error {
	docid := s.ds.Add(rpcdoc.Store, rpcdoc.DocVal)
	if docid == 0 {
		return fmt.Errorf("error adding document (%+v)\n", rpcdoc.DocVal)
	}
	*resp = docid
	return nil
}

func (s *Server) Set(rpcdoc RPCDoc, resp *bool) error {
	s.ds.Set(rpcdoc.Store, rpcdoc.DocId, rpcdoc.DocVal)
	*resp = true
	return nil
}

func (s *Server) Has(rpcdoc RPCDoc, resp *bool) error {
	s.ds.Has(rpcdoc.Store, rpcdoc.DocId)
	*resp = true
	return nil
}

func (s *Server) Get(rpcdoc RPCDoc, resp *Doc) error {
	doc := s.ds.Get(rpcdoc.Store, rpcdoc.DocId)
	if doc == nil {
		return fmt.Errorf("error getting document (%d)\n", rpcdoc.DocId)
	}
	*resp = *doc
	return nil
}

func (s *Server) GetAll(rpcdoc RPCDoc, resp *[]*Doc) error {
	docs := s.ds.GetAll(rpcdoc.Store, rpcdoc.DocIds...)
	if docs == nil || len(docs) < 1 {
		return fmt.Errorf("error getting all documents (%d)\n", rpcdoc.DocIds)
	}
	*resp = docs
	return nil
}

func (s *Server) Del(rpcdoc RPCDoc, resp *bool) error {
	s.ds.Del(rpcdoc.Store, rpcdoc.DocId)
	*resp = true
	return nil
}

package dbdb

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"time"
)

type M map[string]interface{}

type Server struct {
	ds *DataStore
}

func NewServer(ds *DataStore) *Server {
	return &Server{
		ds: ds,
	}
}

func (s *Server) ListenAndServe(port string) {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go s.HandleConn(conn)
	}
}

func (s *Server) HandleConn(conn net.Conn) {
	dec, enc := json.NewDecoder(conn), json.NewEncoder(conn)
	for {
		var m M
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			conn.Close()
			return
		} else {
			conn.SetDeadline(time.Now().Add(time.Minute * 30))
		}
		switch m["cmd"].(string) {
		case "addstore":
			enc.Encode(s.AddStore(&m))
		case "getstore":
			enc.Encode(s.GetStore(&m))
		case "delstore":
			enc.Encode(s.DelStore(&m))
		case "add":
			enc.Encode(s.Add(&m))
		case "set":
			enc.Encode(s.Set(&m))
		case "get":
			enc.Encode(s.Get(&m))
		case "del":
			enc.Encode(s.Del(&m))
		case "exit":
			conn.SetDeadline(time.Now())
		default:
			enc.Encode(M{"error": "unknown command"})
		}
	}
}

func (s *Server) AddStore(m *M) *M {
	store, ok := (*m)["store"]
	(*m)["res"] = ok
	if ok {
		s.ds.AddStore(store.(string))
	}
	return m
}

func (s *Server) GetStore(m *M) *M {
	store, ok := (*m)["store"]
	(*m)["res"] = ok
	if ok {
		st, got := s.ds.GetStore(store.(string))
		if got {
			(*m)["store"] = st
		}
		(*m)["res"] = got
	}
	return m
}

func (s *Server) DelStore(m *M) *M {
	store, ok := (*m)["store"]
	(*m)["res"] = ok
	if ok {
		s.ds.DelStore(store.(string))
	}
	return m
}

func (s *Server) Add(m *M) *M {
	store, ok := (*m)["store"]
	(*m)["res"] = ok
	if ok {
		data, hasdat := (*m)["data"]
		if hasdat {
			id := s.ds.Add(store.(string), data)
			(*m)["id"] = uint64(id)
		}
		(*m)["res"] = hasdat
	}
	return m
}

func (s *Server) Set(m *M) *M {
	store, ok := (*m)["store"]
	(*m)["res"] = ok
	if ok {
		id, hasid := (*m)["id"]
		if hasid {
			data, hasdat := (*m)["data"]
			if hasdat {
				s.ds.Set(store.(string), uint64(id.(float64)), data)
			}
			(*m)["res"] = hasdat
		}
		(*m)["res"] = hasid
	}
	return m
}

func (s *Server) Get(m *M) *M {
	store, ok := (*m)["store"]
	(*m)["res"] = ok
	if ok {
		id, hasid := (*m)["id"]
		if hasid {
			doc := s.ds.Get(store.(string), uint64(id.(float64)))
			(*m)["doc"] = doc.Data
		}
		(*m)["res"] = hasid
	}
	return m
}

func (s *Server) Del(m *M) *M {
	store, ok := (*m)["store"]
	(*m)["res"] = ok
	if ok {
		id, hasid := (*m)["id"]
		if hasid {
			s.ds.Del(store.(string), uint64(id.(float64)))
		}
		(*m)["res"] = hasid
	}
	return m
}

package dbdb

import (
	"encoding/gob"
	"log"
	"net/rpc"
)

func init() {
	// might not work, maybe move to new client
	gob.Register([]interface{}(nil))
}

type Client struct {
	conn  *rpc.Client
	State bool
}

func NewClient(dsn string) *Client {
	return &Client{State: false}
}

func (c *Client) Connect(dsn string) error {
	conn, err := rpc.Dial("tcp", dsn)
	if err == nil {
		c.conn = conn
		c.State = true
	}
	return err
}

func (c *Client) Disconnect() error {
	c.State = false
	return c.conn.Close()
}

func (c *Client) GetAllStoreStats() []*StoreStat {
	var stats []*StoreStat
	Log(c.conn.Call(Server("GetAllStoreStats"), struct{}{}, &stats))
	return stats
}

func (c *Client) GetStoreStat(name string) *StoreStat {
	var stat *StoreStat
	Log(c.conn.Call(Server("GetStoreStat"), name, &stat))
	return stat
}

func (c *Client) HasStore(name string) bool {
	var ok bool
	Log(c.conn.Call(Server("HasStore"), name, &ok))
	return ok
}

func (c *Client) DelStore(name string) bool {
	var ok bool
	Log(c.conn.Call(Server("DelStore"), name, &ok))
	return ok
}

func (c *Client) AddStore(name string) bool {
	var ok bool
	Log(c.conn.Call(Server("AddStore"), name, &ok))
	return ok
}

func (c *Client) Add(name string, val interface{}) uint64 {
	var id uint64
	rpcdoc := RPCDoc{
		Store:  name,
		DocVal: ToMap(val),
	}
	Log(c.conn.Call(Server("Add"), rpcdoc, &id))
	return id
}

func (c *Client) Set(name string, id uint64, val interface{}) bool {
	var ok bool
	rpcdoc := RPCDoc{
		Store:  name,
		DocId:  id,
		DocVal: ToMap(val),
	}
	Log(c.conn.Call(Server("Set"), rpcdoc, &ok))
	return ok
}

func (c *Client) Get(name string, id uint64) *Doc {
	var doc *Doc
	rpcdoc := RPCDoc{
		Store: name,
		DocId: id,
	}
	Log(c.conn.Call(Server("Get"), rpcdoc, &doc))
	return doc
}

func (c *Client) GetAll(name string, id ...uint64) []*Doc {
	var docs []*Doc
	rpcdoc := RPCDoc{
		Store:  name,
		DocIds: id,
	}
	Log(c.conn.Call(Server("GetAll"), rpcdoc, &docs))
	return docs
}

func (c *Client) Del(name string, id uint64) bool {
	var ok bool
	rpcdoc := RPCDoc{
		Store: name,
		DocId: id,
	}
	Log(c.conn.Call(Server("Del"), rpcdoc, &ok))
	return ok
}

func Server(method string) string {
	return "RPCServer." + method
}

func Log(err error) {
	if err != nil {
		log.Println(err)
	}
}

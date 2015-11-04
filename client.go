package dbdb

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"net/rpc"
)

func init() {
	gob.Register([]interface{}(nil))
}

// helper used to wrap up the rpc caller method string
func RPC(method string) string {
	return "Server." + method
}

// helper used to wrap up Call, and log an errors
func Log(err error) {
	if err != nil {
		log.Println(err)
	}
}

type Client struct {
	conn *rpc.Client
}

func NewClient() *Client {
	return &Client{}
}

func CanConnect(conn net.Conn, token string) bool {
	_, err := conn.Write([]byte(fmt.Sprintf("authtoken:%s\n", token)))
	if err != nil {
		log.Printf("client: %v\n", err)
		return false
	}
	resp, _ := bufio.NewReader(conn).ReadString('\n')
	if resp != "1\n" {
		log.Println("client: failed to authenticate with server")
		return false
	}
	return true
}

func (c *Client) Connect(host string, token string) bool {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		log.Printf("client: %v\n", err)
		return false
	}
	if CanConnect(conn, token) {
		c.conn = rpc.NewClient(conn)
		return true
	}
	return false
}

func (c *Client) Disconnect() error {
	return c.conn.Close()
}

func (c *Client) Alive() bool {
	var nothing *struct{}
	err := c.conn.Call(RPC("Alive"), struct{}{}, &nothing)
	if err != nil {
		return false
	}
	return true
}

func (c *Client) GetAllStoreStats() []*StoreStat {
	var stats []*StoreStat
	Log(c.conn.Call(RPC("GetAllStoreStats"), struct{}{}, &stats))
	return stats
}

func (c *Client) GetStoreStat(name string) *StoreStat {
	var stat *StoreStat
	Log(c.conn.Call(RPC("GetStoreStat"), name, &stat))
	return stat
}

func (c *Client) HasStore(name string) bool {
	var ok bool
	Log(c.conn.Call(RPC("HasStore"), name, &ok))
	return ok
}

func (c *Client) DelStore(name string) bool {
	var ok bool
	Log(c.conn.Call(RPC("DelStore"), name, &ok))
	return ok
}

func (c *Client) AddStore(name string) bool {
	var ok bool
	Log(c.conn.Call(RPC("AddStore"), name, &ok))
	return ok
}

func (c *Client) Add(name string, val interface{}) uint64 {
	var id uint64
	rpcdoc := RPCDoc{
		Store:  name,
		DocVal: ToMap(val),
	}
	Log(c.conn.Call(RPC("Add"), rpcdoc, &id))
	return id
}

func (c *Client) Set(name string, id uint64, val interface{}) bool {
	var ok bool
	rpcdoc := RPCDoc{
		Store:  name,
		DocId:  id,
		DocVal: ToMap(val),
	}
	Log(c.conn.Call(RPC("Set"), rpcdoc, &ok))
	return ok
}

func (c *Client) Has(name string, id uint64) bool {
	var ok bool
	rpcdoc := RPCDoc{
		Store: name,
		DocId: id,
	}
	Log(c.conn.Call(RPC("Has"), rpcdoc, &ok))
	return ok
}

func (c *Client) Get(name string, id uint64) *Doc {
	var doc *Doc
	rpcdoc := RPCDoc{
		Store: name,
		DocId: id,
	}
	Log(c.conn.Call(RPC("Get"), rpcdoc, &doc))
	return doc
}

func (c *Client) GetAll(name string, id ...uint64) []*Doc {
	var docs []*Doc
	rpcdoc := RPCDoc{
		Store:  name,
		DocIds: id,
	}
	Log(c.conn.Call(RPC("GetAll"), rpcdoc, &docs))
	return docs
}

func (c *Client) Del(name string, id uint64) bool {
	var ok bool
	rpcdoc := RPCDoc{
		Store: name,
		DocId: id,
	}
	Log(c.conn.Call(RPC("Del"), rpcdoc, &ok))
	return ok
}

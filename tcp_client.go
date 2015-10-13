package dbdb

import (
	"encoding/json"
	"net"
	"sync"
)

type Client struct {
	conn net.Conn
	enc  *json.Encoder
	dec  *json.Decoder
	host string
	sync.RWMutex
}

func NewClient(host string) *Client {
	return &Client{
		host: host,
	}
}

func (c *Client) Open() bool {
	if conn, err := net.Dial("tcp", c.host); err == nil {
		c.enc, c.dec = json.NewEncoder(conn), json.NewDecoder(conn)
		c.conn = conn
		return true
	}
	return false
}

func (c *Client) Close() bool {
	if err := c.conn.Close(); err != nil {
		return false
	}
	return true
}

func (c *Client) Call(m M) M {
	c.enc.Encode(m)
	var res M
	c.dec.Decode(&res)
	return res
}

func (c *Client) AddStore(store string) bool {
	m := c.Call(M{"cmd": "addstore", "store": store})
	return m["res"].(bool)
}

func (c *Client) GetStore(store string) *Store {
	m := c.Call(M{"cmd": "getstore", "store": store})
	return m["store"].(*Store)
}

func (c *Client) DelStore(store string) bool {
	m := c.Call(M{"cmd": "delstore", "store": store})
	return m["res"].(bool)
}

func (c *Client) Add(store string, data interface{}) uint64 {
	m := c.Call(M{"cmd": "add", "store": store, "data": data})
	return uint64(m["id"].(float64))
}

func (c *Client) Set(store string, id uint64, data interface{}) bool {
	m := c.Call(M{"cmd": "set", "store": store, "id": id, "data": data})
	return m["res"].(bool)
}

func (c *Client) Get(store string, id uint64) map[string]interface{} {
	m := c.Call(M{"cmd": "get", "store": store, "id": id})
	return m["doc"].(map[string]interface{})
}

func (c *Client) GetAs(store string, id uint64, v interface{}) {
	m := c.Call(M{"cmd": "get", "store": store, "id": id})
	d := m["doc"].(map[string]interface{})
	ToStruct(d, v)
}

func (c *Client) Del(store string, id uint64) bool {
	m := c.Call(M{"cmd": "del", "store": store, "id": id})
	return m["res"].(bool)
}

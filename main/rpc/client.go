package main

import (
	"encoding/gob"
	"fmt"
	"log"

	"github.com/cagnosolutions/dbdb"
)

type User struct {
	Name   []string
	Age    int
	Active bool
}

var store = "users"

func main() {

	gob.Register([]interface{}(nil))

	// connect to server
	rpcc := dbdb.NewRPCClient("127.0.0.1:9999")

	// check if store exists...
	if ok, _ := rpcc.HasStore(store); !ok {
		// if not, create it
		if ok, err := rpcc.AddStore(store); !ok {
			log.Fatal("rpcc.AddStore() -> ", err)
		}
	}

	// add a document
	id, err := rpcc.Add(dbdb.RPCDoc{store, 0, map[string]interface{}{
		"Name":   []string{"John", "Doe"},
		"Age":    28.00,
		"Active": true,
	}})
	if err != nil {
		log.Fatal("Error adding document: ", err)
	}
	log.Printf("ID: %d\n", id)

	// get document previous to the one just added...
	if id > 1 {
		id = id - 1
	}
	doc, err := rpcc.Get(dbdb.RPCDoc{store, uint64(id), nil})
	if err != nil {
		log.Fatal("Error getting document: ", err)
	} else {
		fmt.Printf("Successfully got doc as map: %+#v\n", doc)
		var u User
		doc.As(&u)
		fmt.Printf("Successfully got doc as struct: %+#v\n", u)
	}

	// close connection
	if err := rpcc.Close(); err != nil {
		log.Fatal("rpcc.Close() -> ", err)
	}

}

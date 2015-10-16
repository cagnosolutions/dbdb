package main

import (
	"encoding/gob"
	"fmt"
	"log"

	"github.com/cagnosolutions/dbdb"
)

type Widget struct {
	Id        uint32
	Name      string
	Relations []uint32
	Reorder   bool
}

var store = "widgets"

func main() {

	gob.Register([]interface{}(nil))

	// connect to server
	rpcc := dbdb.NewRPCClient("127.0.0.1:9999")

	/*
		// check if store exists...
		if ok, _ := rpcc.HasStore(store); !ok {
			// if not, create it
			if ok, err := rpcc.AddStore(store); !ok {
				log.Fatal("rpcc.AddStore() -> ", err)
			}
		}

		// add a document
		id, err := rpcc.Add(dbdb.RPCDoc{store, 0, map[string]interface{}{
			"id":        uint32(12345),
			"name":      "My Awesome Widget",
			"relations": []uint32{23, 4382, 43, 329, 3124},
			"reorder":   true,
		}, nil})
		if err != nil {
			log.Fatal("Error adding document: ", err)
		}
		log.Printf("ID: %d\n", id)

		// get document previous to the one just added...
		if id > 1 {
			id = id - 1
		}
		doc, err := rpcc.Get(dbdb.RPCDoc{store, uint64(id), nil, nil})
		if err != nil {
			log.Fatal("Error getting document: ", err)
		} else {
			fmt.Printf("Successfully got doc as map: %+#v\n", doc)
			var w Widget
			doc.As(&w)
			fmt.Printf("Successfully got doc as struct: %+#v\n", w)
		}
	*/

	/*docs, err := rpcc.GetAll(store, uint64(3), uint64(1), uint64(4), uint64(2))
	if err != nil {
		log.Fatal("Error getting all documents: ", err)
	}
	for _, doc := range docs {
		fmt.Printf("Doc: %v\n", doc)
	}*/
	stats, err := rpcc.GetAllStoreStats()
	if err != nil {
		log.Fatal("Error getting all stats: ", err)
	}
	for _, stat := range stats {
		fmt.Printf("Stat: %v\n", stat)
	}
	// close connection
	if err := rpcc.Close(); err != nil {
		log.Fatal("rpcc.Close() -> ", err)
	}

	// new Client
	/*rpc := dbdb.NewClient("127.0.0.1:9999")
	if err := rpc.Connect("127.0.0.1:9999"); err != nil {
		log.Fatal(err)
	}

	stats := rpc.GetAllStoreStats()
	for _, stat := range stats {
		fmt.Printf("Stat: %v\n", stat)
	}

	rpc.Disconnect()*/
}

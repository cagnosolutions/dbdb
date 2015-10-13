package main

import (
	"fmt"
	"log"

	"github.com/cagnosolutions/dbdb"
)

func main() {
	log.Println("RPC Client Dialing (127.0.0.1:9999)...")
	rpcc := dbdb.NewRPCClient("127.0.0.1:9999")

	store := "rpcstore"

	if _, err := rpcc.GetStore(store); err != nil {
		if ok, err := rpcc.AddStore(store); !(*ok) {
			log.Fatal("rpcc.AddStore() -> ", err)
		}
	}

	rpcdoc := dbdb.RPCDoc{
		Store: store,
		DocId: 0,
		DocVal: map[string]interface{}{
			"name":   []string{"Scott", "Cagno"},
			"age":    28,
			"active": false,
		},
	}

	id, err := rpcc.Add(rpcdoc)
	if err != nil {
		log.Fatal("rpcc.Add(rpcdoc) -> ", err)
	}
	fmt.Printf("rpcc.Add(rpcdoc) -> SUCCESSFULLY ADDED DOCUMENT (%d)\n", id)

	if err := rpcc.Close(); err != nil {
		log.Fatal("rpcc.Close() -> ", err)
	}
}

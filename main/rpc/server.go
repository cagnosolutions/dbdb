package main

import (
	"log"

	"github.com/cagnosolutions/dbdb"
)

func main() {
	rpcserv := dbdb.NewRPCServer(dbdb.NewDataStore())
	log.Println("RPC Server Listening (0.0.0.0:9999)...")
	rpcserv.ListenAndServe(":9999")
}

package main

import (
	"github.com/cagnosolutions/dbdb"
)

var auth string = "foobar"

func main() {
	ds := dbdb.NewDataStore()
	dbdb.Serve(ds, "localhost:9999", auth)
}

package main

import (
	"github.com/cagnosolutions/dbdb"
)

var auth string = "foobar"

func main() {
	ds := dbdb.NewDataStore()
	dbdb.Serve(ds, "0.0.0.0:31337", auth)
}

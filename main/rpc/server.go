package main

import "github.com/cagnosolutions/dbdb"

func main() {

	dbdb.Serve(dbdb.NewDataStore(), ":9999")

}

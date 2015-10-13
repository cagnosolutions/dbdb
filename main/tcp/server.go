package main

import "github.com/cagnosolutions/dbdb"

func main() {
	dbdb.NewDataStore().Listen(":9999")
}

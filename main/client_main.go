package main

import (
	"fmt"
	"log"

	"github.com/cagnosolutions/dbdb"
)

type Widget struct {
	Code   int
	Name   string
	Unique bool
}

var ids []uint64

func main() {

	client := dbdb.NewClient("localhost:9999")

	if ok := client.Open(); !ok {
		log.Fatal("Error -> client.Open()")
	}

	store := "widgets"

	/*
		if ok := client.AddStore(store); !ok {
			log.Fatal("Error -> client.AddStore(...)")
		}
	*/

	/*
		if id := client.Add(store, Widget{123, "Widget 123", true}); id == 0 {
			log.Fatal("Error -> client.Add(...)")
		}
	*/

	var w Widget
	client.GetAs(store, uint64(3), &w)
	fmt.Printf("%+#v\n", w)

	/*
		doc := client.Get(store, uint64(3))
		fmt.Printf("%+#v\n", doc)
	*/

	if ok := client.Close(); !ok {
		log.Fatal("Error -> client.Close()")
	}

}

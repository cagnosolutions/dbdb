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

	// Create a new client connection
	client := dbdb.NewClient("localhost:9999")

	// Open the connection
	if ok := client.Open(); !ok {
		log.Fatal("Error -> client.Open()")
	}

	store := "widgets"

	// If there is no store "widgets", create it
	if st := client.GetStore(store); st == nil {
		if has := client.AddStore(store); !has {
			log.Fatal("Error -> client.AddStore(...)")
		}
	}

	// Add a widget to the "widgets" store
	if id := client.Add(store, Widget{123, "Widget 123", true}); id == 0 {
		log.Fatal("Error -> client.Add(...)")
	}

	// Get widget 1 from the "widgets" store as a widget
	var w Widget
	client.GetAs(store, uint64(1), &w)
	fmt.Printf("%+#v\n", w)

	// Get widget 1 from the "widgets" store as a map[string]interface{}
	doc := client.Get(store, uint64(1))
	fmt.Printf("%+#v\n", doc)

	// Close connection to server
	if ok := client.Close(); !ok {
		log.Fatal("Error -> client.Close()")
	}

}

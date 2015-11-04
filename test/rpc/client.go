package main

import (
	"log"

	"github.com/cagnosolutions/dbdb"
)

var auth string = "foobar"

func main() {
	client := dbdb.NewClient()
	if ok := client.Connect("192.168.0.81:9999", auth); !ok {
		log.Fatal("error connecting to host...")
	}
	if ok := client.HasStore("foobar"); !ok {
		client.AddStore("foobar")
	}
	client.Add("foobar", map[string]interface{}{
		"id":     1,
		"name":   []string{"Scott", "Cagno"},
		"email":  "scottiecagno@gmail.com",
		"age":    28,
		"active": true,
	})
	client.Disconnect()
}

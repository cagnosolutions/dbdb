package main

import (
	"log"

	"github.com/cagnosolutions/dbdb"
)

var auth string = "9999"

func main() {
	client := dbdb.NewClient()
	if err := client.Connect("192.168.0.81:9999", auth); err != nil {
		log.Fatal(err)
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

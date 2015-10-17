package main

import (
	"log"

	"github.com/cagnosolutions/dbdb"
)

var auth string = "foobar"

func main() {
	client := dbdb.NewClient()
	if err := client.Connect("localhost:9999", auth); err != nil {
		log.Fatal(err)
	}
	if ok := client.HasStore("users"); !ok {
		client.AddStore("users")
	}
	client.Add("users", map[string]interface{}{
		"id":     1,
		"name":   []string{"Scott", "Cagno"},
		"email":  "scottiecagno@gmail.com",
		"age":    28,
		"active": true,
	})
	client.Disconnect()
}

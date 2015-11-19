package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/cagnosolutions/dbdb"
)

var (
	queryStr string = "QUERY users WHERE id ^ `0`, email = `scottiecagno@gmail.com`, name = `Scott Cagno`"
	data     *Data  = NewData("users")
)

func main() {

	rs := data.Query(queryStr)
	fmt.Println(rs)
}

type Data struct {
	name string
	data []map[string]interface{}
}

func NewData(name string) *Data {
	d := &Data{
		name: name, // pretend store name
		data: make([]map[string]interface{}, 5),
	}
	// pretend that the index is the primary id
	d.data[0] = map[string]interface{}{
		"name":   "Scott Cagno",
		"email":  "scottiecagno@gmail.com",
		"age":    28,
		"active": true,
	}
	d.data[1] = map[string]interface{}{
		"name":   "Kayla Cagno",
		"email":  "kaylacagno@gmail.com",
		"age":    26,
		"active": false,
	}
	d.data[2] = map[string]interface{}{
		"name":   "Greg Pechiro",
		"email":  "gregpechiro@gmail.com",
		"age":    29,
		"active": true,
	}
	d.data[3] = map[string]interface{}{
		"name":   "Rosalie Pechiro",
		"email":  "rosaliepechiro@gmail.com",
		"age":    30,
		"active": true,
	}
	d.data[4] = map[string]interface{}{
		"name":   "Gabe Witmer",
		"email":  "gabenwitmer@gmail.com",
		"age":    26,
		"active": false,
	}
	return d
}

func (d *Data) Query(query string) []interface{} {
	parser := dbdb.NewParser(strings.NewReader(query))
	stmt, err := parser.Parse()
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	if d.name != stmt.Store {
		return nil
	}
	fmt.Println(stmt.Store)
	//var results []interface{}
	for _, set := range stmt.Set {
		fmt.Printf("%+v\n", set)
	}
	return nil
}

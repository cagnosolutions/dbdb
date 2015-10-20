package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/cagnosolutions/dbdb"
)

var queries = []string{
	"SELECT * FROM users",
	"SELECT id, name, email FROM users",
}

func main() {

	var ast []*dbdb.SelectStatement

	fmt.Println("Sample select statements...")
	for _, query := range queries {
		fmt.Println(query)
	}

	for _, query := range queries {
		parser := dbdb.NewParser(strings.NewReader(query))
		stmt, err := parser.Parse()
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		ast = append(ast, stmt)
	}

	for _, stmt := range ast {
		fmt.Printf("%+#v\n", stmt)
	}
}

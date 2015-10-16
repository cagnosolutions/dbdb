package dbdb

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func WriteStore(path string) {
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Fatalf("WriteStore() -> os.MkdirAll() -> %v\n", err)
	}
}

func WriteDoc(filepath string, doc interface{}) {
	data, err := json.Marshal(doc)
	if err != nil {
		log.Fatalf("WriteDoc() -> %v\n", err)
	}
	if err := ioutil.WriteFile(filepath, data, 0644); err != nil {
		log.Fatalf("WriteDoc() -> ioutil.WriteFile() -> %v\n", err)
	}
}

func DeleteStore(path string) {
	if err := os.Remove(path); err != nil {
		log.Fatalf("DeleteStore() -> os.Remove() -> %v\n", err)
	}
}

func DeleteDoc(filepath string) {
	if err := os.Remove(filepath); err != nil {
		log.Fatalf("DeleteDoc() -> os.Remove() -> %v\n", err)
	}
}

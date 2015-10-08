package dbdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync/atomic"
)

type Store struct {
	name string
	stid uint64
	docs *DocMap
}

func NewStore(name string) *Store {
	return &Store{
		name: name,
		docs: NewDocMap(),
	}
}

func (st *Store) Load(files []string) {
	var id uint64
	for _, file := range files {
		info := strings.Split(file, ".")
		id, err := strconv.ParseUint(info[0], 10, 64)
		if err != nil || len(info) != 2 {
			log.Fatalf("Store.Load() -> invalid file (%v), possible corruption?\n", file)
		}
		data, err := ioutil.ReadFile("db/" + st.name + "/" + file)
		if err != nil {
			log.Fatalf("Store.Load() -> invalid file (%v), possible corruption?\n", file)
		}
		var doc Doc
		if err := json.Unmarshal(data, &doc); err != nil {
			log.Fatalf("Store.Load() -> error unmarshaling data from file (%v), possible corruption?\n", file)
		}
		st.docs.Set(id, &doc)
	}
	st.stid = id
}

/*
func (st *Store) Load(files []string) {
	var id uint64
	for _, file := range files {
		info := strings.Split(file, ".")
		id, err := strconv.ParseUint(info[0], 10, 64)
		if err != nil || len(info) != 2 {
			log.Fatalf("Store.Load() -> invalid file (%v), possible corruption?\n", file)
		}
		data, err := ioutil.ReadFile("db/" + st.name + "/" + file)
		if err != nil {
			log.Fatalf("Store.Load() -> invalid file (%v), possible corruption?\n", file)
		}
		var val map[string]interface{}
		if err := json.Unmarshal(data, &val); err != nil {
			log.Fatalf("Store.Load() -> error unmarshaling data from file (%v), possible corruption?\n", file)
		}
		st.docs.Set(id, NewDoc(id, val))
	}
	st.stid = id
}
*/

func (st *Store) Add(val interface{}) uint64 {
	stid := atomic.AddUint64(&st.stid, uint64(1))
	doc := NewDoc(stid, val)
	st.docs.Set(stid, doc)
	func() {
		WriteDoc(fmt.Sprintf("db/%s/%d.json", st.name, stid), doc)
	}()
	return stid
}

func (st *Store) Set(id uint64, val interface{}) {
	if doc, ok := st.docs.Get(id); ok {
		doc.Update(val)
		st.docs.Set(id, doc)
		func() {
			WriteDoc(fmt.Sprintf("db/%s/%d.json", st.name, id), doc)
		}()
	}
}

func (st *Store) Get(id uint64) *Doc {
	if doc, ok := st.docs.Get(id); ok {
		return doc
	}
	return nil
}

func (st *Store) Del(id uint64) {
	st.docs.Del(id)
	func() {
		DeleteDoc(fmt.Sprintf("db/%s/%d.json", st.name, id))
	}()
}

package dbdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
)

type StoreStat struct {
	Name     string
	Id, Docs uint64
}

type StoreStatSorted []*StoreStat

func (sss StoreStatSorted) Len() int {
	return len(sss)
}

func (sss StoreStatSorted) Less(i, j int) bool {
	return sss[i].Name < sss[j].Name
}

func (sss StoreStatSorted) Swap(i, j int) {
	sss[i], sss[j] = sss[j], sss[i]
}

type Store struct {
	Name    string
	StoreId uint64
	Docs    *DocMap
}

func NewStore(Name string) *Store {
	return &Store{
		Name: Name,
		Docs: NewDocMap(),
	}
}

func (st *Store) Load(files []string) {
	var docid uint64
	for _, file := range files {
		info := strings.Split(file, ".")
		id, err := strconv.ParseUint(info[0], 10, 64)
		if err != nil || len(info) != 2 {
			log.Fatalf("Store.Load() -> invalid file (%v), possible corruption?\n", file)
		}
		docid = id
		data, err := ioutil.ReadFile("db/" + st.Name + "/" + file)
		if err != nil {
			log.Fatalf("Store.Load() -> invalid file (%v), possible corruption?\n", file)
		}
		var doc Doc
		if err := json.Unmarshal(data, &doc); err != nil {
			log.Fatalf("Store.Load() -> error unmarshaling data from file (%v), possible corruption?\n", file)
		}
		st.Docs.Set(id, &doc)
	}
	st.StoreId = docid
}

func (st *Store) Add(val interface{}) uint64 {
	StoreId := atomic.AddUint64(&st.StoreId, uint64(1))
	doc := NewDoc(StoreId, val)
	st.Docs.Set(StoreId, doc)
	func() {
		WriteDoc(fmt.Sprintf("db/%s/%d.json", st.Name, StoreId), doc)
	}()
	return StoreId
}

func (st *Store) Set(id uint64, val interface{}) {
	if doc, ok := st.Docs.Get(id); ok {
		doc.Update(val)
		st.Docs.Set(id, doc)
		func() {
			WriteDoc(fmt.Sprintf("db/%s/%d.json", st.Name, id), doc)
		}()
	}
}

func (st *Store) Has(id uint64) bool {
	_, ok := st.Docs.Get(id)
	return ok
}

func (st *Store) Get(id uint64) *Doc {
	if doc, ok := st.Docs.Get(id); ok {
		return doc
	}
	return nil
}

func (st *Store) GetAll(id ...uint64) []*Doc {
	if len(id) == 0 {
		return st.Docs.GetAll()
	}
	var docs DocSorted
	for _, docid := range id {
		if doc, ok := st.Docs.Get(docid); ok {
			docs = append(docs, doc)
		}
	}
	sort.Sort(docs)
	return docs
}

func (st *Store) Del(id uint64) {
	st.Docs.Del(id)
	func() {
		DeleteDoc(fmt.Sprintf("db/%s/%d.json", st.Name, id))
	}()
}

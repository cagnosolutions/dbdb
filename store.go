package dbdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync/atomic"
)

type StoreStat struct {
	Name     string
	Id, Docs uint64
	Size     float64
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

func (st *Store) Load(ids []int) {
	var docid uint64
	for _, id := range ids {
		docid = uint64(id)
		file := fmt.Sprintf("db/%s/%d.json", st.Name, docid)
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatalf("Store.Load() -> invalid file (%v), possible corruption?\n", file)
		}
		var doc Doc
		if err := json.Unmarshal(data, &doc); err != nil {
			log.Fatalf("Store.Load() -> error unmarshaling data from file (%v), possible corruption?\n", file)
		}
		st.Docs.Set(docid, &doc)
	}
	st.StoreId = docid
}

// size on disk, not document count
func (st *Store) Size() float64 {
	var size int64
	for i := 1; uint64(i) < st.StoreId; i++ {
		if info, err := os.Lstat(fmt.Sprintf("db/%s/%d.json", st.Name, i)); err == nil {
			size += info.Size()
		}
	}
	return toFixed(float64(size)/float64(1<<10), 2)
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

func (st *Store) Query(stmt map[string][]interface{}) []*Doc {
	// all docs, sorted
	docs := st.Docs.GetAll()
	return docs
}

package dbdb

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"sort"
	"sync"
)

type DataStore struct {
	Stores map[string]*Store
	sync.RWMutex
}

func NewDataStore() *DataStore {
	ds := &DataStore{
		Stores: make(map[string]*Store),
	}
	if data := Walk("db"); len(data) > 0 {
		ds.Lock()
		for store, ids := range data {
			st := NewStore(store)
			ds.Stores[store] = st
			st.Load(ids)
		}
		ds.Unlock()
		runtime.GC()
	}
	return ds
}

func (ds *DataStore) Import(data []byte) error {
	var dataStore map[string][]map[string]interface{}
	if err := json.Unmarshal(data, &dataStore); err != nil {
		log.Println("data import failed: ", err)
		return err
	}
	for store, docs := range dataStore {
		ds.AddStore(store)
		for _, doc := range docs {
			ds.Add(store, doc)
		}
	}
	runtime.GC()
	return nil
}

func (ds *DataStore) Export() ([]byte, error) {
	ds.Lock()
	dataStore := make(map[string][]map[string]interface{}, 0)
	for name, store := range ds.Stores {
		allDocs := store.GetAll()
		docData := make([]map[string]interface{}, 0)
		for _, doc := range allDocs {
			docData = append(docData, doc.Data)
		}
		dataStore[name] = docData
	}
	data, err := json.Marshal(dataStore)
	if err != nil {
		log.Println("data export failed: ", err)
		return nil, err
	}
	ds.Unlock()
	defer runtime.GC()
	return data, nil
}

func (ds *DataStore) ClearAll() {
	ds.Lock()
	ds.Stores = nil
	runtime.GC()
	ds.Stores = make(map[string]*Store, 0)
	if err := os.RemoveAll("db"); err != nil {
		log.Println("error clearing data: ", err)
	}
	ds.Unlock()
}

func (ds *DataStore) GetAllStoreStats() []*StoreStat {
	var stats StoreStatSorted
	ds.RLock()
	for name, store := range ds.Stores {
		stats = append(stats, &StoreStat{
			Name: name,
			Id:   store.StoreId,
			Docs: store.Docs.Size(),
			Size: store.Size(),
		})
	}
	ds.RUnlock()
	sort.Stable(stats)
	return stats
}

func (ds *DataStore) GetStoreStat(name string) *StoreStat {
	var stat *StoreStat
	ds.RLock()
	if st, ok := ds.Stores[name]; ok {
		stat = &StoreStat{
			Name: st.Name,
			Id:   st.StoreId,
			Docs: st.Docs.Size(),
			Size: st.Size(),
		}
	}
	ds.RUnlock()
	return stat
}

func (ds *DataStore) HasStore(name string) bool {
	ds.RLock()
	_, ok := ds.Stores[name]
	ds.RUnlock()
	return ok
}

func (ds *DataStore) AddStore(name string) {
	if _, ok := ds.GetStore(name); !ok {
		ds.Lock()
		ds.Stores[name] = NewStore(name)
		func() {
			WriteStore(fmt.Sprintf("db/%s/", name))
		}()
		ds.Unlock()
	}
}

// FOR INTERNAL USE ONLY
func (ds *DataStore) GetStore(name string) (*Store, bool) {
	ds.RLock()
	st, ok := ds.Stores[name]
	ds.RUnlock()
	return st, ok
}

func (ds *DataStore) DelStore(name string) {
	if _, ok := ds.GetStore(name); ok {
		ds.Lock()
		delete(ds.Stores, name)
		func() {
			DeleteStore(fmt.Sprintf("db/%s/", name))
		}()
		ds.Unlock()
	}
}

func (ds *DataStore) Add(name string, val interface{}) uint64 {
	if st, ok := ds.GetStore(name); ok {
		return st.Add(val)
	}
	return 0
}

func (ds *DataStore) Set(name string, id uint64, val interface{}) {
	if st, ok := ds.GetStore(name); ok {
		st.Set(id, val)
	}
}

func (ds *DataStore) Has(name string, id uint64) bool {
	if st, ok := ds.GetStore(name); ok {
		return st.Has(id)
	}
	return false
}

func (ds *DataStore) Get(name string, id uint64) *Doc {
	if st, ok := ds.GetStore(name); ok {
		return st.Get(id)
	}
	return nil
}

func (ds *DataStore) GetAll(name string, id ...uint64) DocSorted {
	if st, ok := ds.GetStore(name); ok {
		return st.GetAll(id...)
	}
	return nil
}

func (ds *DataStore) Del(name string, id uint64) {
	if st, ok := ds.GetStore(name); ok {
		st.Del(id)
	}
}

func (ds *DataStore) Query(name string, comps ...QueryComp) DocSorted {
	if st, ok := ds.GetStore(name); ok {
		if len(comps) > 0 {
			return st.Query(comps...)
		}
	}
	return nil
}

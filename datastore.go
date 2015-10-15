package dbdb

import (
	"fmt"
	"runtime"
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
		for store, files := range data {
			st := NewStore(store)
			ds.Stores[store] = st
			st.Load(files)
		}
		ds.Unlock()
		runtime.GC()
	}
	return ds
}

type StoreStat struct {
	Name     string // sort by this
	Id, Docs uint64
}

func (ds *DataStore) GetAllStoreStats() []*StoreStat {
	var stats []*StoreStat
	ds.RLock()
	for name, store := range ds.Stores {
		stats = append(stats, &StoreStat{
			Name: name,
			Id:   store.StoreId,
			Docs: store.Docs.Size(),
		})
	}
	ds.RUnlock()
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

func (ds *DataStore) Get(name string, id uint64) *Doc {
	if st, ok := ds.GetStore(name); ok {
		return st.Get(id)
	}
	return nil
}

func (ds *DataStore) Del(name string, id uint64) {
	if st, ok := ds.GetStore(name); ok {
		st.Del(id)
	}
}

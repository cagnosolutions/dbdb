package dbdb

import "sync/atomic"

type Store struct {
	stid uint64
	docs *DocMap
}

func NewStore() *Store {
	return &Store{
		docs: NewDocMap(),
	}
}

func (st *Store) Add(val interface{}) uint64 {
	stid := atomic.AddUint64(&st.stid, uint64(1))
	st.docs.Set(stid, &Doc{stid, val})
	return stid
}

func (st *Store) Set(id uint64, val interface{}) {
	// only set if document exists
	if doc, ok := st.docs.Get(id); ok {
		doc.Data = val
		st.docs.Set(id, doc)
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
}

package dbdb

import (
	"encoding/binary"
	"hash/fnv"
	"log"
	"sort"
	"sync"
)

const SHARD_COUNT = 64

type DocMap []*Shard

type Shard struct {
	Docs map[float64]*Doc
	sync.RWMutex
}

func NewDocMap() *DocMap {
	m := make(DocMap, SHARD_COUNT)
	for i := 0; uint32(i) < SHARD_COUNT; i++ {
		m[i] = &Shard{
			Docs: make(map[float64]*Doc),
		}
	}
	return &m
}

func getshard(n float64) uint64 {
	hasher := fnv.New64a()
	key := make([]byte, 10)
	binary.PutVarint(key, int64(n))
	_, err := hasher.Write(key)
	if err != nil {
		log.Fatal(err)
	}
	return hasher.Sum64() % SHARD_COUNT
}

func (m *DocMap) Size() uint64 {
	var count uint64
	for doc := range m.Iter() {
		if doc != nil {
			count++
		}
	}
	return count
}

func (m *DocMap) GetShard(key float64) *Shard {
	return (*m)[getshard(key)]
}

func (m *DocMap) Set(id float64, val *Doc) {
	shard := m.GetShard(id)
	shard.Lock()
	shard.Docs[id] = val
	shard.Unlock()
}

func (m *DocMap) Get(id float64) (*Doc, bool) {
	shard := m.GetShard(id)
	shard.RLock()
	val, ok := shard.Docs[id]
	shard.RUnlock()
	return val, ok
}

func (m *DocMap) GetAll() DocSorted {
	var docs DocSorted
	for doc := range m.Iter() {
		if doc != nil {
			docs = append(docs, doc)
		}
	}
	sort.Sort(docs)
	return docs
}

func (m *DocMap) Del(id float64) {
	if shard := m.GetShard(id); shard != nil {
		shard.Lock()
		delete(shard.Docs, id)
		shard.Unlock()
	}
}

func (m *DocMap) Iter() <-chan *Doc {
	ch := make(chan *Doc)
	go func() {
		for _, shard := range *m {
			shard.RLock()
			for _, doc := range shard.Docs {
				ch <- doc
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}

func (m *DocMap) Query(comps ...QueryComp) DocSorted {
	var docs DocSorted
	var match bool
	for doc := range m.Iter() {
		match = true
		for _, comp := range comps {
			if v, ok := doc.Data[comp.Field()]; !ok || !comp.Comp(v) {
				match = false // got mismatch
			}
		}
		if match { // no mismatches
			docs = append(docs, doc)
		}
	}
	sort.Sort(docs)
	return docs
}

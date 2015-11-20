package dbdb

import (
	"encoding/binary"
	"hash/fnv"
	"log"
	"sort"
	"sync"
)

const SHARD_COUNT uint32 = 64

type DocMap []*Shard

type Shard struct {
	Docs map[uint64]*Doc
	sync.RWMutex
}

func NewDocMap() *DocMap {
	m := make(DocMap, SHARD_COUNT)
	for i := 0; uint32(i) < SHARD_COUNT; i++ {
		m[i] = &Shard{
			Docs: make(map[uint64]*Doc),
		}
	}
	return &m
}

func getshard(n uint64) uint32 {
	hasher := fnv.New32a()
	key := make([]byte, 10)
	binary.PutUvarint(key, n)
	_, err := hasher.Write(key)
	if err != nil {
		log.Fatal(err)
	}
	return hasher.Sum32() % SHARD_COUNT
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

func (m *DocMap) GetShard(key uint64) *Shard {
	bucket := getshard(key)
	return (*m)[bucket]
}

func (m *DocMap) Set(id uint64, val *Doc) {
	shard := m.GetShard(id)
	shard.Lock()
	shard.Docs[id] = val
	shard.Unlock()
}

func (m *DocMap) Get(id uint64) (*Doc, bool) {
	shard := m.GetShard(id)
	shard.RLock()
	val, ok := shard.Docs[id]
	shard.RUnlock()
	return val, ok
}

func (m *DocMap) GetAll() []*Doc {
	var docs DocSorted
	for doc := range m.Iter() {
		if doc != nil {
			docs = append(docs, doc)
		}
	}
	sort.Sort(docs)
	return docs
}

func (m *DocMap) Del(id uint64) {
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

func (m *DocMap) Query(comps ...QueryComp) []*Doc {
	results := make([]*Doc, 0)
	for _, comp := range comps {
		for doc := range m.Iter() {
			if val, ok := doc.Data[comp.Field()]; ok {
				if comp.Comp(val) {
					results = append(results, doc)
				}
			}
		}
	}
	return results
}

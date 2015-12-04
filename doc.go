package dbdb

import "time"

type Doc struct {
	Id       uint64                 `json:"id"`
	Created  int64                  `json:"created"`
	Modified int64                  `json:"modified"`
	Data     map[string]interface{} `json:"data"`
}

func NewDoc(id uint64, data interface{}) *Doc {
	time := time.Now().Unix()
	doc := &Doc{
		Id:       id,
		Created:  time,
		Modified: time,
	}
	switch data.(type) {
	case map[string]interface{}:
		doc.Data = data.(map[string]interface{})
		SanitizeMapNums(doc.Data)
	default:
		doc.Data = ToMap(data)
	}
	return doc
}

func (d *Doc) Update(data interface{}) {
	switch data.(type) {
	case map[string]interface{}:
		d.Data = data.(map[string]interface{})
	default:
		d.Data = ToMap(data)
	}
	d.Modified = time.Now().Unix()
}

func (d *Doc) As(v interface{}) bool {
	if err := ToStruct(d.Data, v); err != nil {
		return false
	}
	return true
}

type DocSorted []*Doc

func (ds DocSorted) Len() int {
	return len(ds)
}

func (ds DocSorted) Less(i, j int) bool {
	return ds[i].Id < ds[j].Id
}

func (ds DocSorted) Swap(i, j int) {
	ds[i], ds[j] = ds[j], ds[i]
}

func (ds DocSorted) One() *Doc {
	if len(ds) >= 1 {
		return ds[0]
	}
	return &Doc{}
}

func (ds DocSorted) Limit(count int) DocSorted {
	if count >= 1 && count <= len(ds) {
		return ds[0:count]
	}
	return ds // 0 reutrns all results for now
}

func (ds DocSorted) Page(idx, count int) DocSorted {
	page := idx + count
	if page < len(ds) {
		return ds[idx:count]
	} else if idx > len(ds)-1 {
		return ds
	}
	return ds[idx:]
}

func (ds DocSorted) Ids() []uint64 {
	var ids []uint64
	for _, doc := range ds {
		ids = append(ids, doc.Id)
	}
	return ids
}

func (ds DocSorted) Fields(name string) []interface{} {
	var flds []interface{}
	for _, doc := range ds {
		if v, ok := doc.Data[name]; ok {
			flds = append(flds, v)
		}
	}
	return flds
}

type DocSet struct {
	Id   uint64
	Data map[string]interface{}
}

func (ds DocSorted) Data() <-chan DocSet {
	ch := make(chan DocSet)
	go func() {
		for _, doc := range ds {
			ch <- DocSet{doc.Id, doc.Data}
		}
		close(ch)
	}()
	return ch
}

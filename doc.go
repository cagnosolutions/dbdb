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

func (d *Doc) As(v interface{}) error {
	return ToStruct(d.Data, v)
}

// SORT
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

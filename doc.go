package dbdb

type Doc struct {
	Id   uint64                 `json:"id"`
	Data map[string]interface{} `json:"data"`
}

func NewDoc(id uint64, data interface{}) *Doc {
	doc := &Doc{Id: id}
	switch data.(type) {
	case map[string]interface{}:
		doc.Data = data.(map[string]interface{})
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
}

func (d *Doc) As(v interface{}) error {
	return ToStruct(d.Data, v)
}

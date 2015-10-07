package dbdb

type Doc struct {
	Id   uint64      `json:"id"`
	Data interface{} `json:"data"`
}

func (d *Doc) AsMap() map[string]interface{} {
	return ToMap(d.Data)
}

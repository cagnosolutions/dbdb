package dbdb

import "fmt"

type QueryComp interface {
	Comp(v interface{}) bool
	Field() string
}

type Eq struct {
	Fld string
	Val interface{}
}

func (c Eq) Comp(v interface{}) bool {
	return fmt.Sprint(v) == fmt.Sprint(c.Val)
}

func (c Eq) Field() string {
	return c.Fld
}

type Ne struct {
	Fld string
	Val interface{}
}

func (c Ne) Comp(v interface{}) bool {
	return fmt.Sprint(v) != fmt.Sprint(c.Val)
}

func (c Ne) Field() string {
	return c.Fld
}

type Gt struct {
	Fld string
	Val interface{}
}

func (c Gt) Comp(v interface{}) bool {
	var dbVal float64
	switch v.(type) {
	case float64:
		dbVal = v.(float64)
	default:
		return false
	}
	switch c.Val.(type) {
	case int:
		return dbVal > float64(c.Val.(int))
	case int8:
		return dbVal > float64(c.Val.(int8))
	case int16:
		return dbVal > float64(c.Val.(int16))
	case int32:
		return dbVal > float64(c.Val.(int32))
	case int64:
		return dbVal > float64(c.Val.(int64))
	case uint:
		return dbVal > float64(c.Val.(uint))
	case uint8:
		return dbVal > float64(c.Val.(uint8))
	case uint16:
		return dbVal > float64(c.Val.(uint16))
	case uint32:
		return dbVal > float64(c.Val.(uint32))
	case uint64:
		return dbVal > float64(c.Val.(uint64))
	case float32:
		return dbVal > float64(c.Val.(float32))
	default:
		return false
	}
}

func (c Gt) Field() string {
	return c.Fld
}

type Lt struct {
	Fld string
	Val interface{}
}

func (c Lt) Comp(v interface{}) bool {
	var dbVal float64
	switch v.(type) {
	case float64:
		dbVal = v.(float64)
	default:
		return false
	}
	switch c.Val.(type) {
	case int:
		return dbVal < float64(c.Val.(int))
	case int8:
		return dbVal < float64(c.Val.(int8))
	case int16:
		return dbVal < float64(c.Val.(int16))
	case int32:
		return dbVal < float64(c.Val.(int32))
	case int64:
		return dbVal < float64(c.Val.(int64))
	case uint:
		return dbVal < float64(c.Val.(uint))
	case uint8:
		return dbVal < float64(c.Val.(uint8))
	case uint16:
		return dbVal < float64(c.Val.(uint16))
	case uint32:
		return dbVal < float64(c.Val.(uint32))
	case uint64:
		return dbVal < float64(c.Val.(uint64))
	case float32:
		return dbVal < float64(c.Val.(float32))
	default:
		return false
	}
}

func (c Lt) Field() string {
	return c.Fld
}

package dbdb

import (
	"fmt"
	"reflect"
)

type QueryComp interface {
	Comp(v interface{}) bool
	Field() string
}

type Eq struct {
	Fld string
	Val interface{}
}

func (c Eq) Comp(v interface{}) bool {
	return reflect.DeepEqual(v, c.Val)
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
	return fmt.Sprint(v) > fmt.Sprint(c.Val)
}

func (c Gt) Field() string {
	return c.Fld
}

type Lt struct {
	Fld string
	Val interface{}
}

func (c Lt) Comp(v interface{}) bool {
	return fmt.Sprint(v) < fmt.Sprint(c.Val)
}

func (c Lt) Field() string {
	return c.Fld
}

/*
type QueryEngine struct {
	Queries []*Stmt
}

func NewQueryEngine() *QueryEngine {
	return &QueryEngine{
		Queries: make([]*Stmt, 0),
	}
}

func (qe *QueryEngine) Add(f, c string, v interface{}) *QueryEngine {
	qe.Queries = append(qe.Queries, &Stmt{f, c, v})
	return qe
}

func (qe *QueryEngine) HasMatch(doc *Doc) bool {
	if doc == nil {
		return false
	}
	for _, stmt := range qe.Queries {
		if docVal, ok := doc.Data[stmt.Fld]; ok {
			switch stmt.Comp {
			case "eq":
				if docVal == stmt.Val {
					return true
				}
			case "ne":
				if docVal != stmt.Val {
					return true
				}
			case "lt":
				return lt(docVal, stmt.Val)
			case "gt":
				return gt(docVal, stmt.Val)
			default:
				return false
			}
		}
	}
	return false
}

func lt(v1, v2 interface{}) bool {
	return fmt.Sprint(v1) < fmt.Sprint(v2)
}

func gt(v1, v2 interface{}) bool {
	return fmt.Sprint(v1) > fmt.Sprint(v2)
}

type Stmt struct {
	Fld, Comp string
	Val       interface{}
}
*/

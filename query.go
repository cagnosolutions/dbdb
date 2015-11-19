package dbdb

import "fmt"

type QueryComp interface {
	Comp(v interface{}) bool
	Field() string
}

type Eq struct {
	field string
	value interface{}
}

func (c Eq) Comp(v interface{}) bool {
	return v == c.value
}

func (c Eq) Field() string {
	return c.field
}

type Ne struct {
	field string
	value interface{}
}

func (c Ne) Comp(v interface{}) bool {
	return v != c.value
}

func (c Ne) Field() string {
	return c.field
}

type Gt struct {
	field string
	value interface{}
}

func (c Gt) Comp(v interface{}) bool {
	return fmt.Sprint(v) > fmt.Sprint(c.value)
}

func (c Gt) Field() string {
	return c.field
}

type Lt struct {
	field string
	value interface{}
}

func (c Lt) Comp(v interface{}) bool {
	return fmt.Sprint(v) < fmt.Sprint(c.value)
}

func (c Lt) Field() string {
	return c.field
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
		if docVal, ok := doc.Data[stmt.field]; ok {
			switch stmt.Comp {
			case "eq":
				if docVal == stmt.value {
					return true
				}
			case "ne":
				if docVal != stmt.value {
					return true
				}
			case "lt":
				return lt(docVal, stmt.value)
			case "gt":
				return gt(docVal, stmt.value)
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
	field, Comp string
	value       interface{}
}
*/

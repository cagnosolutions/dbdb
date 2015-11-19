package dbdb

import "fmt"

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
		if docVal, ok := doc.Data[stmt.Field]; ok {
			switch stmt.Comp {
			case "eq":
				if docVal == stmt.Value {
					return true
				}
			case "ne":
				if docVal != stmt.Value {
					return true
				}
			case "lt":
				return lt(docVal, stmt.Value)
			case "gt":
				return gt(docVal, stmt.Value)
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
	Field, Comp string
	Value       interface{}
}

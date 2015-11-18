package dbdb

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
			default:
				return false
			}
		}
	}
	return false
}

type Stmt struct {
	Field, Comp string
	Value       interface{}
}

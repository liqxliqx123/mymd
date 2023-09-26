package pariticiple_query

import "fmt"

type BoolQuery map[string]map[string][]map[string]interface{}

func (r *Root) ToESQuery() map[string]interface{} {
	query := r.Left.ToESQuery()
	if len(r.Right) > 0 {
		for _, c := range r.Right {
			query = c.ToESQuery(query)
		}
	}
	return map[string]interface{}{
		"query": query,
	}
}

func (c *Copilot) ToESQuery(existingQuery map[string]interface{}) map[string]interface{} {
	switch c.Op {
	case "AND":
		return map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{existingQuery, c.Right.ToESQuery()},
			},
		}
	case "OR":
		return map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{existingQuery, c.Right.ToESQuery()},
			},
		}
	case "NOT":
		return map[string]interface{}{
			"bool": map[string]interface{}{
				"must_not": c.Right.ToESQuery(),
			},
		}
	default:
		panic("unknown operator: " + c.Op)
	}
}

func (o *Operand) ToESQuery() map[string]interface{} {
	if o.SubOperand != nil {
		// Simplifying, assuming only one expression is present
		var boolSlice []map[string]interface{}
		query := o.SubOperand.Left.ToESQuery(o.Field)
		if o.SubOperand.Right != nil && len(o.SubOperand.Right) > 0 {
			if query != nil && len(o.SubOperand.Right) == 1 && o.SubOperand.Right[0].Op == "NOT" {
				boolSlice = append(boolSlice, query)
			}
			for _, e := range o.SubOperand.Right {
				boolSlice = append(boolSlice, e.ToESQuery(o.Field, query))
			}

		} else {
			boolSlice = append(boolSlice, query)
		}
		return map[string]interface{}{
			"bool": map[string]interface{}{
				"must": boolSlice,
			},
		}
	}
	panic("Operand ToESQuery Panic")
}

func (v ValueWithInterval) ToESQuery(field string) map[string]interface{} {
	query := map[string]interface{}{
		"wildcard": map[string]interface{}{
			field: map[string]interface{}{
				"value": fmt.Sprintf("*%s*", v.Left.Str),
			},
		},
	}
	return query
}

func (e *Expression) ToESQuery(field string, existingQuery map[string]interface{}) map[string]interface{} {
	query := e.Right.ToESQuery(field)
	switch e.Op {
	case "AND":
		return map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{existingQuery, query},
			},
		}
	case "OR":
		return map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{existingQuery, query},
			},
		}
	case "NOT":
		return map[string]interface{}{
			"bool": map[string]interface{}{
				"must_not": []map[string]interface{}{
					query,
				},
			},
		}
	default:
		panic("unknown operator: " + e.Op)
	}
}

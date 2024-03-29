// THIS FILE IS A GENERATED CODE. DO NOT EDIT
// generated version: 0.4.0
package todolist

import (
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/genproto/googleapis/type/latlng"
)

// QueryChainer ...
type QueryChainer struct {
	QueryGroup []*Query
}

// NewQueryChainer - constructor
func NewQueryChainer() *QueryChainer {
	return new(QueryChainer)
}

// Equal - change operator
func (q *QueryChainer) Equal(v interface{}) *QueryChainer {
	q.QueryGroup = append(q.QueryGroup, newQuery(v, OpTypeEqual))
	return q
}

// LessThan - change operator
func (q *QueryChainer) LessThan(v interface{}) *QueryChainer {
	q.QueryGroup = append(q.QueryGroup, newQuery(v, OpTypeLessThan))
	return q
}

// LessThanOrEqual - change operator
func (q *QueryChainer) LessThanOrEqual(v interface{}) *QueryChainer {
	q.QueryGroup = append(q.QueryGroup, newQuery(v, OpTypeLessThanOrEqual))
	return q
}

// GreaterThan - change operator
func (q *QueryChainer) GreaterThan(v interface{}) *QueryChainer {
	q.QueryGroup = append(q.QueryGroup, newQuery(v, OpTypeGreaterThan))
	return q
}

// GreaterThanOrEqual - change operator
func (q *QueryChainer) GreaterThanOrEqual(v interface{}) *QueryChainer {
	q.QueryGroup = append(q.QueryGroup, newQuery(v, OpTypeGreaterThanOrEqual))
	return q
}

// In - change operator
func (q *QueryChainer) In(v interface{}) *QueryChainer {
	q.QueryGroup = append(q.QueryGroup, newQuery(v, OpTypeIn))
	return q
}

// ArrayContains - change operator
func (q *QueryChainer) ArrayContains(v interface{}) *QueryChainer {
	q.QueryGroup = append(q.QueryGroup, newQuery(v, OpTypeArrayContains))
	return q
}

// ArrayContainsAny - change operator
func (q *QueryChainer) ArrayContainsAny(v interface{}) *QueryChainer {
	q.QueryGroup = append(q.QueryGroup, newQuery(v, OpTypeArrayContainsAny))
	return q
}

// Query ...
type Query struct {
	Operator OpType
	Value    interface{}
}

func newQuery(v interface{}, opType OpType) *Query {
	switch x := v.(type) {
	case bool, []bool,
		string, []string,
		int, []int,
		int64, []int64,
		float64, []float64,
		*latlng.LatLng, []*latlng.LatLng,
		*firestore.DocumentRef, []*firestore.DocumentRef,
		map[string]bool,
		map[string]string,
		map[string]int,
		map[string]int64,
		map[string]float64:
		// ok
	case time.Time:
		v = SetLastThreeToZero(x)
	case []time.Time:
		after := make([]time.Time, len(x), len(x))
		for n, t := range x {
			after[n] = SetLastThreeToZero(t)
		}
		v = after
	default:
		panic(fmt.Sprintf("unsupported types: %#v", v))
	}

	return &Query{
		Operator: opType,
		Value:    v,
	}
}

// IsSlice ...
func IsSlice(v interface{}) bool {
	switch v.(type) {
	case []bool, []string, []int, []int64, []float64,
		[]*latlng.LatLng, []*firestore.DocumentRef:
		return true
	}
	return false
}

// IsSlice ...
func (q *Query) IsSlice() bool {
	return IsSlice(q.Value)
}

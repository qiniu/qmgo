package qmgo

import (
	"context"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Query struct definition
type Query struct {
	ctx        context.Context
	collection *mongo.Collection
	filter     interface{}
	sort       interface{}
	project    interface{}
	limit      *int64
	skip       *int64
}

// Sort is Used to set the sorting rules for the returned results
// Format: "age" or "+age" means to sort the age field in ascending order, "-age" means in descending order
// When multiple sort fields are passed in at the same time, they are arranged in the order in which the fields are passed in.
// For example, {"age", "-name"}, first sort by age in ascending order, then sort by name in descending order
func (q *Query) Sort(fields ...string) QueryI {
	var sorts bson.D
	for _, field := range fields {
		key, n := SplitSortField(field)
		if key == "" {
			panic("Sort: empty field name")
		}
		sorts = append(sorts, bson.E{Key: key, Value: n})
	}

	return &Query{
		ctx:        q.ctx,
		collection: q.collection,
		filter:     q.filter,
		sort:       sorts,
		project:    q.project,
		limit:      q.limit,
		skip:       q.skip,
	}
}

// Select is used to determine which fields are displayed or not displayed in the returned results
// Format: bson.M{"age": 1} means that only the age field is displayed
// bson.M{"age": 0} means to display other fields except age
// When _id is not displayed and is set to 0, it will be returned to display
func (q *Query) Select(projection interface{}) QueryI {
	return &Query{
		ctx:        q.ctx,
		collection: q.collection,
		filter:     q.filter,
		sort:       q.sort,
		project:    projection,
		limit:      q.limit,
		skip:       q.skip,
	}
}

// Skip skip n records
func (q *Query) Skip(n int64) QueryI {
	return &Query{
		ctx:        q.ctx,
		collection: q.collection,
		filter:     q.filter,
		sort:       q.sort,
		project:    q.project,
		limit:      q.limit,
		skip:       &n,
	}
}

// Limit limits the maximum number of documents found to n
// The default value is 0, and 0  means no limit, and all matching results are returned
// When the limit value is less than 0, the negative limit is similar to the positive limit, but the cursor is closed after returning a single batch result.
// Reference https://docs.mongodb.com/manual/reference/method/cursor.limit/index.html
func (q *Query) Limit(n int64) QueryI {
	return &Query{
		ctx:        q.ctx,
		collection: q.collection,
		filter:     q.filter,
		sort:       q.sort,
		project:    q.project,
		limit:      &n,
		skip:       q.skip,
	}
}

// One query a record that meets the filter conditions
// If the search fails, an error will be returned
func (q *Query) One(result interface{}) error {
	opt := options.FindOne()

	if q.sort != nil {
		opt.SetSort(q.sort)
	}
	if q.project != nil {
		opt.SetProjection(q.project)
	}
	if q.skip != nil {
		opt.SetSkip(*q.skip)
	}

	err := q.collection.FindOne(q.ctx, q.filter, opt).Decode(result)

	if err != nil {
		return err
	}
	return err
}

// All query multiple records that meet the filter conditions
// The static type of result must be a slice pointer
func (q *Query) All(result interface{}) error {
	opt := options.Find()

	if q.sort != nil {
		opt.SetSort(q.sort)
	}
	if q.project != nil {
		opt.SetProjection(q.project)
	}
	if q.limit != nil {
		opt.SetLimit(*q.limit)
	}
	if q.skip != nil {
		opt.SetSkip(*q.skip)
	}

	var err error
	var cursor *mongo.Cursor

	cursor, err = q.collection.Find(q.ctx, q.filter, opt)

	c := Cursor{
		ctx:    q.ctx,
		cursor: cursor,
		err:    err,
	}
	return c.All(result)
}

// Count count the number of eligible entries
func (q *Query) Count() (n int64, err error) {
	opt := options.Count()

	if q.limit != nil {
		opt.SetLimit(*q.limit)
	}
	if q.skip != nil {
		opt.SetSkip(*q.skip)
	}

	return q.collection.CountDocuments(q.ctx, q.filter, opt)
}

// Distinct gets the unique value of the specified field in the collection and return it in the form of slice
// result should be passed a pointer to slice
// The function will verify whether the static type of the elements in the result slice is consistent with the data type obtained in mongodb
// reference https://docs.mongodb.com/manual/reference/command/distinct/
func (q *Query) Distinct(key string, result interface{}) error {
	resultVal := reflect.ValueOf(result)

	if resultVal.Kind() != reflect.Ptr {
		return ErrQueryNotSlicePointer
	}

	sliceVal := resultVal.Elem()
	if sliceVal.Kind() == reflect.Interface {
		sliceVal = sliceVal.Elem()
	}
	if sliceVal.Kind() != reflect.Slice {
		return ErrQueryNotSliceType
	}

	if !resultVal.Elem().CanSet() {
		return ErrQueryResultValCanNotChange
	}

	sliceVal = sliceVal.Slice(0, 0)
	elementType := sliceVal.Type().Elem()

	opt := options.Distinct()
	res, err := q.collection.Distinct(q.ctx, key, q.filter, opt)
	if err != nil {
		return err
	}

	for _, v := range res {
		vValue := reflect.ValueOf(v)
		vType := vValue.Type()

		if vType != elementType {
			fmt.Printf("mongo type: %s, result type: %s\n", vType.Name(), elementType.Name())
			return ErrQueryResultTypeInconsistent
		}
		sliceVal = reflect.Append(sliceVal, vValue)
	}

	resultVal.Elem().Set(sliceVal.Slice(0, len(res)))
	return nil
}

// Cursor gets a Cursor object, which can be used to traverse the query result set
// After obtaining the CursorI object, you should actively call the Close interface to close the cursor
func (q *Query) Cursor() CursorI {
	opt := options.Find()

	if q.sort != nil {
		opt.SetSort(q.sort)
	}
	if q.project != nil {
		opt.SetProjection(q.project)
	}
	if q.limit != nil {
		opt.SetLimit(*q.limit)
	}
	if q.skip != nil {
		opt.SetSkip(*q.skip)
	}

	var err error
	var cur *mongo.Cursor
	cur, err = q.collection.Find(q.ctx, q.filter, opt)
	return &Cursor{
		ctx:    q.ctx,
		cursor: cur,
		err:    err,
	}
}

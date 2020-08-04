package mongox

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ERR_QUERY_NOT_SLICE_POINTER         = errors.New("result argument must be a pointer to a slice")
	ERR_QUERY_NOT_SLICE_TYPE            = errors.New("result argument must be a slice address")
	ERR_QUERY_RESULT_TYPE_INCONSISTEN   = errors.New("result type is not equal mongodb value type")
	ERR_QUERY_RESULT_VAL_CAN_NOT_CHANGE = errors.New("the value of result can not be changed")
)

// Query
type Query struct {
	ctx        context.Context
	collection *mongo.Collection
	filter     interface{}
	sort       interface{}
	project    interface{}
	limit      *int64
	skip       *int64
}

// Sort
// Sort 用于设置对返回结果的排序规则
// 格式："age" 或者 "+age" 表示对age字段升序排列，"-age" 表示对age字段降序排列
// 当同时传入多个排序字段时，按照字段传入的顺序排列。例如，{"age", "-name"}，则先按age升序排列，再按name降序排列
func (q *Query) Sort(fields ...string) QueryI {
	var sorts bson.D
	for _, field := range fields {
		key, n := SplitSymbol(field)
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

// Select
// Select 用于决定哪些字段在返回结果中展示或者不展示
// 格式：bson.M{"age": 1} 表示只展示 age 字段
// bson.M{"age": 0} 表示展示除了 age 以外的其他字段
// 当 _id 没有显示的被置为 0 的情况下，都会被返回展示
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

// Skip
// Skip 用于跳过最初始的 n 个文档
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

// Limit
// Limit 将查找到的最大文档数量限制为 n
// 默认值为 0， 或者主动 limit 设置为 0，表示不限制，返回全部匹配结果
// 当 limit 值小于 0 时，负限制类似于正限制，但在返回单个批量结果后关闭光标。
// 参考 https://docs.mongodb.com/manual/reference/method/cursor.limit/index.html
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

// One
// 查询符合filter条件的一条记录
// 如果查找失败，将返回错误
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

	var err error
	data, err := q.collection.FindOne(q.ctx, q.filter, opt).DecodeBytes()

	if err != nil {
		return err
	}

	err = bson.Unmarshal(data, result)
	return err
}

// All
// 查询符合filter条件的多条记录
// result 静态类型必须是一个 slice 的指针
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

	if err != nil {
		return err
	}
	if err_ := cursor.Err(); err_ != nil {
		return err_
	}

	err = cursor.All(q.ctx, result)
	return err
}

// Count
// 统计符合条件的条目数
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

// Distinct
// 获取集合中指定字段的不重复值，并以 slice 的形式返回
// result 应传入一个 slice 的指针
// 函数会校验 result slice 中元素的静态类型是否与 mongodb 里获取的数据类型一致
// 语法参考 https://docs.mongodb.com/manual/reference/command/distinct/
func (q *Query) Distinct(key string, result interface{}) error {
	resultVal := reflect.ValueOf(result)

	if resultVal.Kind() != reflect.Ptr {
		return ERR_QUERY_NOT_SLICE_POINTER
	}

	sliceVal := resultVal.Elem()
	if sliceVal.Kind() == reflect.Interface {
		sliceVal = sliceVal.Elem()
	}
	if sliceVal.Kind() != reflect.Slice {
		return ERR_QUERY_NOT_SLICE_TYPE
	}

	if !resultVal.Elem().CanSet() {
		return ERR_QUERY_RESULT_VAL_CAN_NOT_CHANGE
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
			return ERR_QUERY_RESULT_TYPE_INCONSISTEN
		}
		sliceVal = reflect.Append(sliceVal, vValue)
	}

	resultVal.Elem().Set(sliceVal.Slice(0, len(res)))
	return nil
}

// Cursor
// 获取一个 Cursor 对象，可用于对查询结果集的遍历
// 获取 CursorI 对象后，应主动调用 Close 接口关闭游标
func (q *Query) Cursor() (CursorI, error) {
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
	if err != nil {
		return nil, err
	}

	return &Cursor{
		ctx:    q.ctx,
		cursor: cur,
	}, nil
}

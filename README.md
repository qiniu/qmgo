# Qmgo 

[![Build Status](https://travis-ci.org/qiniu/qmgo.png?branch=master)](https://travis-ci.org/qiniu/qmgo)
[![Coverage Status](https://codecov.io/gh/qiniu/qmgo/branch/master/graph/badge.svg)](https://codecov.io/gh/qiniu/qmgo)
[![Go Report Card](https://goreportcard.com/badge/github.com/qiniu/qmgo)](https://goreportcard.com/report/github.com/qiniu/qmgo)
[![GitHub release](https://img.shields.io/github/v/tag/qiniu/qmgo.svg?label=release)](https://github.com/qiniu/qmgo/releases)
[![GoDoc](https://pkg.go.dev/badge/github.com/qiniu/qmgo?status.svg)](https://pkg.go.dev/github.com/qiniu/qmgo?tab=doc)

[简体中文](README_ZH.md)

`Qmgo` is a `MongoDB` `driver` for `Go` . It is based on [MongoDB official driver](https://github.com/mongodb/mongo-go-driver), but easier to use like [mgo](https://github.com/go-mgo/mgo) (such as the chain call). 

- `Qmgo` can allow user to use the new features of `MongoDB` in a more elegant way.

- `Qmgo` is the first choice for migrating from `mgo` to the new `MongoDB driver` with minimal code changes.

## Requirements

-`Go 1.10` and above.

-`MongoDB 2.6` and above.

## Features
- CRUD to documents
- Options when create connection: connection pool、pool monitor、auth、readPreference
- Create indexes、Drop indexes
- Sort、limit、count、select、distinct
- Cursor
- Aggregate
- Transactions
- Predefine operator keys

## Installation

- Use `go mod` to automatically install dependencies by `import github.com/qiniu/qmgo`

Or 

- Use `go get github.com/qiniu/qmgo`

## Usage

- Start

`import` and create a new connection
```go
import (
    "context"
  
    "github.com/qiniu/qmgo"
)

ctx := context.Background()
client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: "mongodb://localhost:27017"})
db := client.Database("class")
coll := db.Collection("user")
```
If your connection points to a fixed database and collection, we recommend using the following more convenient way to initialize the connection. 
All operations are based on `cli` and no longer need to care about the database and collection.

```go
cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://localhost:27017", Database: "class", Coll: "user"})
```

***The following examples will be based on `cli`, if you use the first way for initialization, replace `cli` with `client`、`db` or `coll`***

After the initialization is successful, please `defer` to close the connection

```go
defer func() {
if err = cli.Close(ctx); err != nil {
        panic(err)
    }
}()
```

- Create index

Before doing the operation, we first initialize some data:

```go
type UserInfo struct {
	Name   string `bson:"name"`
	Age    uint16 `bson:"age"`
	Weight uint32 `bson:"weight"`
}

var oneUserInfo = UserInfo{
    Name: "xm",
    Age: 7,
    Weight: 40,
}
```

Create index

```go
cli.EnsureIndexes(ctx, []string{}, []string{"age", "name,weight"})
```

- Insert a document

```go
// insert one document
result, err := cli.Insert(ctx, oneUserInfo)
```

- Find a document

```go
// find one document
  one := UserInfo{}
  err = cli.Find(ctx, bson.M{"name": oneUserInfo.Name}).One(&one)
```

- Delete documents

```go
err = cli.Remove(ctx, bson.M{"age": 7})
```

- Insert multiple data

```go
// multiple insert
var batchUserInfoI = []interface{}{
	UserInfo{Name: "a1", Age: 6, Weight: 20},
	UserInfo{Name: "b2", Age: 6, Weight: 25},
	UserInfo{Name: "c3", Age: 6, Weight: 30},
	UserInfo{Name: "d4", Age: 6, Weight: 35},
	UserInfo{Name: "a1", Age: 7, Weight: 40},
	UserInfo{Name: "a1", Age: 8, Weight: 45},
}
result, err = cli.Collection.InsertMany(ctx, batchUserInfoI)
```

- Search all, sort and limit
```go
// find all, sort and limit
batch := []UserInfo{}
cli.Find(ctx, bson.M{"age": 6}).Sort("weight").Limit(7).All(&batch)
```
- Count
````go
count, err := cli.Find(ctx, bson.M{"age": 6}).Count()
````

- Update
````go
// UpdateOne one
err := cli.UpdateOne(ctx, bson.M{"name": "d4"}, bson.M{"$set": bson.M{"age": 7}})

// UpdateAll
result, err := cli.UpdateAll(ctx, bson.M{"age": 6}, bson.M{"$set": bson.M{"age": 10}})
````

- Select
````go
err := cli.Find(ctx, bson.M{"age": 10}).Select(bson.M{"age": 1}).One(&one)
````

- Aggregate
```go
matchStage := bson.D{{"$match", []bson.E{{"weight", bson.D{{"$gt", 30}}}}}}
groupStage := bson.D{{"$group", bson.D{{"_id", "$name"}, {"total", bson.D{{"$sum", "$age"}}}}}}
var showsWithInfo []bson.M
err = cli.Aggregate(context.Background(), Pipeline{matchStage, groupStage}).All(&showsWithInfo)
```

- Pool Monitor
````go
poolMonitor := &event.PoolMonitor{
	Event: func(evt *event.PoolEvent) {
		switch evt.Type {
		case event.GetSucceeded:
			fmt.Println("GetSucceeded")
		case event.ConnectionReturned:
			fmt.Println("ConnectionReturned")
		}
	},
}
cli, err := Open(ctx, &Config{Uri: URI, Database: DATABASE, Coll: COLL, PoolMonitor: poolMonitor})

````

- Transactions

The super simple and powerful transaction, with features like `timeout`、`retry`:
````go
callback := func(sessCtx context.Context) (interface{}, error) {
    // Important: make sure the sessCtx used in every operation in the whole transaction
    if _, err := cli.InsertOne(sessCtx, bson.D{{"abc", int32(1)}}); err != nil {
        return nil, err
    }
    if _, err := cli.InsertOne(sessCtx, bson.D{{"xyz", int32(999)}}); err != nil {
        return nil, err
    }
    return nil, nil
}
result, err = cli.DoTransaction(ctx, callback)
````
[More about transaction](https://github.com/qiniu/qmgo/wiki/Transactions)

- Predefine operator keys
````go
// aggregate
matchStage := bson.D{{operator.Match, []bson.E{{"weight", bson.D{{operator.Gt, 30}}}}}}
groupStage := bson.D{{operator.Group, bson.D{{"_id", "$name"}, {"total", bson.D{{operator.Sum, "$age"}}}}}}
var showsWithInfo []bson.M
err = cli.Aggregate(context.Background(), Pipeline{matchStage, groupStage}).All(&showsWithInfo)
````

## `Qmgo` vs `go.mongodb.org/mongo-driver`

Below we give an example of multi-file search、sort and limit to illustrate the similarities between `qmgo` and `mgo` and the improvement compare to `go.mongodb.org/mongo-driver`.
How do we do in`go.mongodb.org/mongo-driver`:

```go
// go.mongodb.org/mongo-driver
// find all, sort and limit
findOptions := options.Find()
findOptions.SetLimit(7) // set limit
var sorts D
sorts = append(sorts, E{Key: "weight", Value: 1})
findOptions.SetSort(sorts) // set sort

batch := []UserInfo{}
cur, err := coll.Find(ctx, bson.M{"age": 6}, findOptions)
cur.All(ctx, &batch)
```

How do we do in `Qmgo` and `mgo`:

```go
// qmgo
// find all, sort and limit
batch := []UserInfo{}
cli.Find(ctx, bson.M{"age": 6}).Sort("weight").Limit(7).All(&batch)

// mgo
// find all, sort and limit
coll.Find(bson.M{"age": 6}).Sort("weight").Limit(7).All(&batch)
```

## `Qmgo` vs `mgo`
[Differences between qmgo and mgo](https://github.com/qiniu/qmgo/wiki/Known-differences-between-Qmgo-and-Mgo)
 
## Who is using
If you are using qmgo, please feel free to add your project name or repository here！

- Qiniu QCDN management system
- Qiniu RTC quality monitoring system
- Jesselivermore huanshoulv stock real-time quotes system


## Contributing

The Qmgo project welcomes all contributors. We appreciate your help! 

## Join qmgo wechat group:

![avatar](http://pgo8q04yu.bkt.clouddn.com/qmgoG-2)



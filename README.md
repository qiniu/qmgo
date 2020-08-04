# Qmgo 

[简体中文](README_ZH.md)

`Qmgo` is a `MongoDB` `dirver` for `Go` . It is based on [Mongo official driver](https://github.com/mongodb/mongo-go-driver) with better interface design which mainly refer to [mgo](https://github.com/go-mgo/ mgo) (such as the chain call). 

Because the interface design of  [Mongo official driver](https://github.com/mongodb/mongo-go-driver) is not friendly enough，and the well-designed driver `mgo` is no longer maintained.

So we hope `qmgo` can allow users to use the new features of `MongoDB` in a more elegant way.

At the same time, we Hope that `qmgo` is the first choice for migrating from `mgo` to the new `MongoDB driver`.

## Requirements

-`Go 1.10` and above.

-`MongoDB 2.6` and above.

## Installation

The recommended way is to use `go mod` to automatically install dependencies by `import github.com/qiniu/qmgo` and `build` .

Of course, the following methods are also feasible:

```
go get github.com/qiniu/qmgo
```

## Usage

- Start, `import` and create a new connection
```
import(
    "context"
  
    "github.com/qiniu/qmgo"
)	
    ctx := context.Background()
    client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: "mongodb://localhost:27017"})
    db := client.Database("class")
    coll := db.Collection("user")
```
If your connection points to a fixed database and collection, we recommend using the following more convenient method to initialize the connection. The subsequent operations are based on `cli` and no longer need to care about the database and collection.

```go
cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://localhost:27017", Database: "class", Coll: "user"})
```

***The following examples will be based on `cli`, if you use the first method for initialization, replace `cli` with `coll`***

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
type BsonT map[string]interface{}

type UserInfo struct {
    Name string `bson:"name"`
    Age uint16 `bson:"age"`
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
cli.EnsureIndexes(ctx, []string{"name"}, []string{"age", "name,weight"})
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
  err = cli.Find(ctx, BsonT{"name": oneUserInfo.Name}).One(&one)
```

- Delete documents

```go
err = cli.Remove(ctx, BsonT{"age": 7})
```

- Insert multiple data

```go
// batch insert
var batchUserInfoI = []interface{}{
    UserInfo{Name: "wxy", Age: 6, Weight: 20},
    UserInfo{Name: "jZ", Age: 6, Weight: 25},
    UserInfo{Name: "zp", Age: 6, Weight: 30},
    UserInfo{Name: "yxw", Age: 6, Weight: 35},
}
result, err = cli.Collection.InsertMany(ctx, batchUserInfoI)
```

- Search all, sort and limit

```go
// find all, sort and limit
batch := []UserInfo{}
cli.Find(ctx, BsonT{"age": 6}).Sort("weight").Limit(7).All(&batch)
```

## Feature

- Supported
  - CRUD to documents
  - Create indexes
  - Sort、limit、count
- TODO
  - Transaction
  - Aggregate
  - Options for every operation



## `qmgo` vs `mgo` vs `go.mongodb.org/mongo-driver`

Below we give an example of multi-file search、sort and limit to illustrate the similarities between `qmgo` and `mgo` and the improvement compare to `go.mongodb.org/mongo-driver`.

How do we do in`go.mongodb.org/mongo-driver`:

```go
// go.mongodb.org/mongo-driver
// find all, sort and limit
findOptions := options.Find()
findOptions.SetLimit(7) // set limit
var sorts bson.D
sorts = append(sorts, bson.E{Key: "weight", Value: 1})
findOptions.SetSort(sorts) // set sort

batch := []UserInfo{}
cur, err := coll.Find(ctx, BsonT{"age": 6}, findOptions)
cur.All(ctx, &batch)
```

How do we do in `Qmgo` and `mgo`:

```go
// qmgo
// find all, sort and limit
batch := []UserInfo{}
cli.Find(ctx, BsonT{"age": 6}).Sort("weight").Limit(7).All(&batch)

// mgo
// find all, sort and limit
coll.Find(BsonT{"age": 6}).Sort("weight").Limit(7).All(&batch)
```



## contributing

The Qmgo project welcomes all contributors. We appreciate your help! 



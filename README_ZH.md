# Qmgo

`Qmgo` 是一款`Go`语言的`MongoDB` `dirver`，它基于[MongoDB 官方driver](https://github.com/mongodb/mongo-go-driver) 开发实现，同时使用更易用的接口设计，主要参考[mgo](https://github.com/go-mgo/mgo) （比如`mgo`的链式调用）。

- `Qmgo`能让用户以更优雅的姿势使用`MongoDB`的新特性。

- `Qmgo`是从`mgo`迁移到新`MongoDB driver`的第一选择，对代码的改动影响最小。

## 要求

- `Go 1.10` 及以上。
- `MongoDB 2.6` 及以上。

## 安装

推荐方式是使用`go mod`，通过在源码中`import github.com/qiniu/qmgo` 并`build` 来自动安装依赖。

当然，通过下面方式同样可行：

```
go get github.com/qiniu/qmgo
```

## Usage

- 开始

`import`并新建连接
```go  
import(
    "context"
    
    "github.com/qiniu/qmgo"
)

ctx := context.Background()
client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: "mongodb://localhost:27017"})
db := client.Database("class")
coll := db.Collection("user")
      
```

如果你的连接是指向固定的database和collection，我们推荐使用下面的更方便的方法初始化连接，后续操作都基于`cli`而不用再关心database和collection

```go
cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://localhost:27017", Database: "class", Coll: "user"})
```

***后面都会基于`cli`来举例，如果你使用第一种传统的方式进行初始化，根据上下文，将`cli`替换成`client`、`db` 或 `coll`即可***

在初始化成功后，请`defer`来关闭连接 

```go
defer func() {
    if err = cli.Close(ctx); err != nil {
        panic(err)
    }
}()
```

- 创建索引

做操作前，我们先初始化一些数据：

```go

type UserInfo struct {
	Name   string `bson:"name"`
	Age    uint16 `bson:"age"`
	Weight uint32 `bson:"weight"`
}

var oneUserInfo = UserInfo{
	Name:   "xm",
	Age:    7,
	Weight: 40,
}	
```

创建索引

```go
cli.EnsureIndexes(ctx, []string{}, []string{"age", "name,weight"})
```

- 插入一个文档

```go
// insert one document
result, err := cli.Insert(ctx, oneUserInfo)
```

- 查找一个文档

```go
	// find one document
one := UserInfo{}
err = cli.Find(ctx, bson.M{"name": oneUserInfo.Name}).One(&one)
```

- 删除文档

```go
err = cli.Remove(ctx, bson.M{"age": 7})
```

- 插入多条数据

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

- 批量查找、`Sort`和`Limit`

```go
// find all 、sort and limit
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

## 功能
- 文档的增删改查
- 索引配置
- `Sort`、`Limit`、`Count`、`Select`
- `Cursor`
- 聚合`Aggregate`
- Pool Monitor


## `qmgo` & `mgo` vs `go.mongodb.org/mongo-driver`

下面我们举一个多文件查找、`sort`和`limit`的例子, 说明`qmgo`和`mgo`的相似，以及对`go.mongodb.org/mongo-driver`的改进

官方`Driver`需要这样实现

```go
// go.mongodb.org/mongo-driver
// find all 、sort and limit
findOptions := options.Find()
findOptions.SetLimit(7)  // set limit
var sorts D
sorts = append(sorts, E{Key: "weight", Value: 1})
findOptions.SetSort(sorts) // set sort

batch := []UserInfo{}
cur, err := coll.Find(ctx, bson.M{"age": 6}, findOptions)
cur.All(ctx, &batch)
```

`Qmgo`和`mgo`更简单，而且实现相似：

```go
// qmgo
// find all 、sort and limit
batch := []UserInfo{}
cli.Find(ctx, bson.M{"age": 6}).Sort("weight").Limit(7).All(&batch)

// mgo
// find all 、sort and limit
coll.Find(bson.M{"age": 6}).Sort("weight").Limit(7).All(&batch)
```

## `Qmgo` vs `mgo`
[Qmgo 和 Mgo 的差异](https://github.com/qiniu/qmgo/wiki/Known-differences-between-Qmgo-and-Mgo)
 
 
## contributing

非常欢迎您对`Qmgo`的任何贡献，非常感谢您的帮助！

## 加入 qmgo 微信群:

![avatar](http://pgo8q04yu.bkt.clouddn.com/qmgoG)

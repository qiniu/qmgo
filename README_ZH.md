# Qmgo

`Qmgo` 是一款`Go`语言的`MongoDB` `driver`，它基于[MongoDB 官方 driver](https://github.com/mongodb/mongo-go-driver) 开发实现，同时使用更易用的接口设计，比如参考[mgo](https://github.com/go-mgo/mgo) （比如`mgo`的链式调用）。

- `Qmgo`让您以更优雅的姿势使用`MongoDB`的新特性。

- `Qmgo`是从`mgo`迁移到新`MongoDB driver`的第一选择，对代码的改动影响最小。

## 要求

- `Go 1.10` 及以上。
- `MongoDB 2.6` 及以上。

## 功能

- 文档的增删改查, 均支持官方driver支持的所有options
- `Sort`、`limit`、`count`、`select`、`distinct`
- 事务
- `Hooks`
- 自动化更新的默认和定制fields
- 预定义操作符
- 聚合`Aggregate`、索引操作、`cursor`
- `validation tags` 基于tag的字段验证
- 可自定义插件化编程

## 安装

推荐方式是使用`go mod`，通过在源码中`import github.com/qiniu/qmgo` 来自动安装依赖。

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
    
    如果你的连接是指向固定的 database 和 collection，我们推荐使用下面的更方便的方法初始化连接，后续操作都基于`cli`而不用再关心 database 和 collection
    
    ```go
    cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://localhost:27017", Database: "class", Coll: "user"})
    ```
    
    **_后面都会基于`cli`来举例，如果你使用第一种传统的方式进行初始化，根据上下文，将`cli`替换成`client`、`db` 或 `coll`即可_**
    
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
    
    var userInfo = UserInfo{
        Name:   "xm",
        Age:    7,
        Weight: 40,
    }
    ```

    创建索引

    ```go
    cli.CreateOneIndex(context.Background(), options.IndexModel{Key: []string{"name"}})
    cli.CreateIndexes(context.Background(), []options.IndexModel{{Key: []string{"id2", "id3"}}})
    ```

- 插入一个文档

    ```go
    // insert one document
    result, err := cli.InsertOne(ctx, userInfo)
    ```

- 查找一个文档

    ```go
        // find one document
    one := UserInfo{}
    err = cli.Find(ctx, bson.M{"name": userInfo.Name}).One(&one)
    ```

- 删除文档

    ```go
    err = cli.Remove(ctx, bson.M{"age": 7})
    ```

- 插入多条数据

    ```go
    // multiple insert
    var userInfos = []UserInfo{
        UserInfo{Name: "a1", Age: 6, Weight: 20},
        UserInfo{Name: "b2", Age: 6, Weight: 25},
        UserInfo{Name: "c3", Age: 6, Weight: 30},
        UserInfo{Name: "d4", Age: 6, Weight: 35},
        UserInfo{Name: "a1", Age: 7, Weight: 40},
        UserInfo{Name: "a1", Age: 8, Weight: 45},
    }
    result, err = cli.Collection.InsertMany(ctx, userInfos)
    ```

- 批量查找、`Sort`和`Limit`

    ```go
    // find all 、sort and limit
    batch := []UserInfo{}
    cli.Find(ctx, bson.M{"age": 6}).Sort("weight").Limit(7).All(&batch)
    ```

- Count

    ```go
    count, err := cli.Find(ctx, bson.M{"age": 6}).Count()
    ```

- Update

    ```go
    // UpdateOne one
    err := cli.UpdateOne(ctx, bson.M{"name": "d4"}, bson.M{"$set": bson.M{"age": 7}})
    
    // UpdateAll
    result, err := cli.UpdateAll(ctx, bson.M{"age": 6}, bson.M{"$set": bson.M{"age": 10}})
    ```

- Select

    ```go
    err := cli.Find(ctx, bson.M{"age": 10}).Select(bson.M{"age": 1}).One(&one)
    ```

- Aggregate

    ```go
    matchStage := bson.D{{"$match", []bson.E{{"weight", bson.D{{"$gt", 30}}}}}}
    groupStage := bson.D{{"$group", bson.D{{"_id", "$name"}, {"total", bson.D{{"$sum", "$age"}}}}}}
    var showsWithInfo []bson.M
    err = cli.Aggregate(context.Background(), Pipeline{matchStage, groupStage}).All(&showsWithInfo)
    ```

- 建立连接时支持所有 mongoDB 的`Options`

    ```go
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
    
    opt := options.Client().SetPoolMonitor(poolMonitor)  // more options use the chain options.
    cli, err := Open(ctx, &Config{Uri: URI, Database: DATABASE, Coll: COLL}, opt)
    
    ```

- 事务

    有史以来最简单和强大的事务, 同时还有超时和重试等功能:
    
    ```go
    callback := func(sessCtx context.Context) (interface{}, error) {
        // 重要：确保事务中的每一个操作，都使用传入的sessCtx参数
        if _, err := cli.InsertOne(sessCtx, bson.D{{"abc", int32(1)}}); err != nil {
            return nil, err
        }
        if _, err := cli.InsertOne(sessCtx, bson.D{{"xyz", int32(999)}}); err != nil {
            return nil, err
        }
        return nil, nil
    }
    result, err = cli.DoTransaction(ctx, callback)
    ```
    
    [关于事务的更多内容](https://github.com/qiniu/qmgo/wiki/Transactions)

- 预定义操作符

    ```go
    // aggregate
    matchStage := bson.D{{operator.Match, []bson.E{{"weight", bson.D{{operator.Gt, 30}}}}}}
    groupStage := bson.D{{operator.Group, bson.D{{"_id", "$name"}, {"total", bson.D{{operator.Sum, "$age"}}}}}}
    var showsWithInfo []bson.M
    err = cli.Aggregate(context.Background(), Pipeline{matchStage, groupStage}).All(&showsWithInfo)
    ```

- Hooks

    Qmgo 灵活的 hooks:

    ```go
    type User struct {
        Name         string    `bson:"name"`
        Age          int       `bson:"age"`
    }
    func (u *User) BeforeInsert(ctx context.Context) error {
        fmt.Println("before insert called")
        return nil
    }
    func (u *User) AfterInsert(ctx context.Context) error {
        fmt.Println("after insert called")
        return nil
    }
    
    u := &User{Name: "Alice", Age: 7}
    _, err := cli.InsertOne(context.Background(), u)
    ```
    
    [Hooks 详情介绍](<https://github.com/qiniu/qmgo/wiki/Hooks--(%E7%AE%80%E4%BD%93%E4%B8%AD%E6%96%87)>)


- 自动化更新fields

    Qmgo支持2种方式来自动化更新特定的字段

    - 默认 fields
    
    在文档结构体里注入 `field.DefaultField`, `Qmgo` 会自动在更新和插入操作时更新 `createAt`、`updateAt` and `_id` field的值.
    
    ````go
    type User struct {
        field.DefaultField `bson:",inline"`
        
        Name string `bson:"name"`
        Age  int    `bson:"age"`
    }
  
  	u := &User{Name: "Lucas", Age: 7}
  	_, err := cli.InsertOne(context.Background(), u)
    // tag为createAt、updateAt 和 _id 的字段会自动更新插入
    ```` 

    - Custom fields
    
    可以自定义field名, `Qmgo` 会自动在更新和插入操作时更新他们.

    ```go
    type User struct {
        Name string `bson:"name"`
        Age  int    `bson:"age"`
    
        MyId         string    `bson:"myId"`
        CreateTimeAt time.Time `bson:"createTimeAt"`
        UpdateTimeAt int64     `bson:"updateTimeAt"`
    }
    // 指定自定义field的field名
    func (u *User) CustomFields() field.CustomFieldsBuilder {
        return field.NewCustom().SetCreateAt("CreateTimeAt").SetUpdateAt("UpdateTimeAt").SetId("MyId")
    }
  
    u := &User{Name: "Lucas", Age: 7}
    _, err := cli.InsertOne(context.Background(), u)
    // CreateTimeAt、UpdateTimeAt and MyId 会自动更新并插入DB 
  
    // 假设Id和ui已经初始化
  	err = cli.ReplaceOne(context.Background(), bson.M{"_id": Id}, &ui)
    // UpdateTimeAt 会被自动更新
    ```
  
    [例子介绍](https://github.com/qiniu/qmgo/blob/master/field_test.go)

    [自动化 fields 详情介绍](https://github.com/qiniu/qmgo/wiki/Automatically-update-fields)
  
- `Validation tags` 基于tag的字段验证
    
    功能基于[go-playground/validator](https://github.com/go-playground/validator)实现。
    
    所以`Qmgo`支持所有[go-playground/validator 的struct验证规则](https://github.com/go-playground/validator#usage-and-documentation)，比如：
    ```go
    type User struct {
        FirstName string            `bson:"fname"`
        LastName  string            `bson:"lname"`
        Age       uint8             `bson:"age" validate:"gte=0,lte=130" `    // Age must in [0,130]
        Email     string            `bson:"e-mail" validate:"required,email"` //  Email can't be empty string, and must has email format
        CreateAt  time.Time         `bson:"createAt" validate:"lte"`          // CreateAt must lte than current time
        Relations map[string]string `bson:"relations" validate:"max=2"`       // Relations can't has more than 2 elements
    }
    ```
  
    本功能只对以下API有效：
    ` InsertOne、InsertyMany、Upsert、UpsertId、ReplaceOne `

- 插件化编程
    
    - 实现以下方法
    ```go
    func Do(ctx context.Context, doc interface{}, opType operator.OpType, opts ...interface{}) error{
      // do anything
    }
    ```
    
    - 调用middleware包的Register方法，注入`Do`
      Qmgo会在支持的[操作](operator/operate_type.go)执行前后调用`Do`
    ```go
    middleware.Register(Do)
    ```
    [Example](middleware/middleware_test.go)
    
    Qmgo的hook、自动更新field和validation tags都基于plugin的方式实现
  
## `qmgo` vs `go.mongodb.org/mongo-driver`

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

[Qmgo 和 Mgo 的差异](https://github.com/qiniu/qmgo/wiki/Differences-between-Qmgo-and-Mgo)

## Contributing

非常欢迎您对`Qmgo`的任何贡献，非常感谢您的帮助！


## 沟通交流:

- 加入 [qmgo discussions](https://github.com/qiniu/qmgo/discussions)

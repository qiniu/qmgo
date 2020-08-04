package mongox

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Config 是 mongodb 客户端配置
// Uri 格式 [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
// 参考 https://docs.mongodb.com/manual/reference/connection-string/
type Config struct {
	Uri         string `json:"uri"`
	DB          string `json:"db"`
	Coll        string `json:"coll"`
	Timeout     int64  `json:"timeout"`     // Timeout 参数是以毫秒为单位
	MaxPollSize uint64 `json:"maxpollsize"` // MaxPollSize 参数用于设定 Client 底层可以建立的最大连接数；不设置情况下默认大小为100；如果设置为0，则表示设置大小为math.MaxInt64.
}

// Client
// mongodb 客户端定义
type Client struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

// 根据 Config 获取一个 Client
func Open(conf *Config) (cli *Client, err error) {
	var opts *options.ClientOptions
	opts = new(options.ClientOptions)

	if conf.Timeout > 0 {
		timeoutDur := time.Duration(conf.Timeout) * time.Second
		opts.SetConnectTimeout(timeoutDur)
		opts.SetServerSelectionTimeout(timeoutDur)
		opts.SetSocketTimeout(timeoutDur)
	}
	if conf.MaxPollSize > 0 {
		opts.SetMaxPoolSize(conf.MaxPollSize)
	}
	//opts.SetMonitor() // for trace
	opts.ApplyURI(conf.Uri)

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		fmt.Println(err)
		return
	}

	db := client.Database(conf.DB)
	coll := db.Collection(conf.Coll)

	cli = &Client{
		client:     client,
		database:   db,
		collection: coll,
	}

	return
}

// Close
// 关闭 Client
func (c *Client) Close(ctx context.Context) error {
	err := c.client.Disconnect(ctx)
	return err
}

// GetCollection
// 获取一个 Collection 对象
// 不一定会建立新的tcp连接，mongo.Collection是并发安全的
// 底层tcp链接建立，由mongo.Client控制，最大连接数为 Config 参数 MaxPollSize 决定
func (c *Client) GetCollection(ctx context.Context) *Collection {
	var err error
	var cp *mongo.Collection

	cp, err = c.collection.Clone()
	if err != nil {
		return nil
	}

	if ctx == nil {
		ctx = context.TODO()
	}

	return &Collection{
		Ctx:        ctx,
		Collection: cp,
	}
}

// GetDatabaseName
// 获取 Database 名称
func (c *Client) GetDatabaseName() string {
	return c.database.Name()
}

// GetCollectionName
// 获取 Collection 名称
func (c *Client) GetCollectionName() string {
	return c.collection.Name()
}

// Ping
func (c *Client) Ping(timeout int64) error {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	if err = c.client.Ping(ctx, readpref.Primary()); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

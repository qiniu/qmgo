package qmgo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Client specifies the instance to operate mongoDB
type QmgoClient struct {
	*Collection
	*Database
	*Client
}

// Open creates client instance according to config
// QmgoClient can operates all mongo.client 、mongo.database and mongo.collection
func Open(ctx context.Context, conf *Config) (cli *QmgoClient, err error) {
	client, err := NewClient(ctx, conf)
	if err != nil {
		fmt.Println("new client fail", err)
		return
	}

	db := client.Database(conf.Database)
	if err != nil {
		fmt.Println("new database fail", err)
		return
	}

	coll := db.Collection(conf.Coll)

	cli = &QmgoClient{
		Client:     client,
		Database:   db,
		Collection: coll,
	}

	return
}

// Client creates client to mongo
type Client struct {
	client *mongo.Client
}

// NewClient creates mongo.client
func NewClient(ctx context.Context, conf *Config) (cli *Client, err error) {
	client, err := client(ctx, conf)
	if err != nil {
		fmt.Println("new client fail", err)
		return
	}
	cli = &Client{
		client: client,
	}
	return
}

// client creates connection to mongo
func client(ctx context.Context, conf *Config) (client *mongo.Client, err error) {
	var opts *options.ClientOptions
	opts = new(options.ClientOptions)

	if conf.ConnectTimeoutMS != nil {
		timeoutDur := time.Duration(*conf.ConnectTimeoutMS) * time.Millisecond
		opts.SetConnectTimeout(timeoutDur)

	}
	if conf.SocketTimeoutMS != nil {
		timeoutDur := time.Duration(*conf.SocketTimeoutMS) * time.Millisecond
		opts.SetSocketTimeout(timeoutDur)
	} else {
		opts.SetSocketTimeout(300 * time.Second)
	}
	if conf.MaxPoolSize != nil {
		opts.SetMaxPoolSize(*conf.MaxPoolSize)
	}
	opts.ApplyURI(conf.Uri)

	client, err = mongo.Connect(ctx, opts)
	if err != nil {
		fmt.Println(err)
		return
	}
	pCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err = client.Ping(pCtx, readpref.Primary()); err != nil {
		fmt.Println(err)
		return
	}
	return
}

// Close closes sockets to the topology referenced by this Client.
func (c *Client) Close(ctx context.Context) error {
	err := c.client.Disconnect(ctx)
	return err
}

// Ping confirm connection is alive
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

// Database create connection to database
func (c *Client) Database(name string) *Database {
	return &Database{database: c.client.Database(name)}
}

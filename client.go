package qmgo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Config for initial mongodb instance
type Config struct {
	// URI example: [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
	// URI Reference: https://docs.mongodb.com/manual/reference/connection-string/
	Uri      string `json:"uri"`
	Database string `json:"database"`
	Coll     string `json:"coll"`
	// ConnectTimeoutMS specifies a timeout that is used for creating connections to the server.
	//	If set to 0, no timeout will be used.
	//	The default is 30 seconds.
	ConnectTimeoutMS *int64 `json:"connectTimeoutMS"`
	// MaxPoolSize specifies that maximum number of connections allowed in the driver's connection pool to each server.
	// If this is 0, it will be set to math.MaxInt64,
	// The default is 100.
	MaxPoolSize *uint64 `json:"maxPoolSize"`
	// SocketTimeoutMS specifies how long the driver will wait for a socket read or write to return before returning a
	// network error. If this is 0 meaning no timeout is used and socket operations can block indefinitely.
	// The default is 300,000 ms.
	SocketTimeoutMS *int64 `json:"socketTimeoutMS"`
	// ReadPreference determines which servers are considered suitable for read operations.
	// default is PrimaryMode
	ReadPreference *ReadPref `json:"readPreference"`
	// PoolMonitor to receive connection pool events
	PoolMonitor *event.PoolMonitor
}

// ReadPref determines which servers are considered suitable for read operations.
type ReadPref struct {
	// MaxStaleness is the maximum amount of time to allow a server to be considered eligible for selection.
	// Supported from version 3.4.
	MaxStalenessMS int64 `json:"maxStalenessMS"`
	// indicates the user's preference on reads.
	// PrimaryMode as default
	Mode readpref.Mode `json:"mode"`
}

// QmgoClient specifies the instance to operate mongoDB
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
	if conf.PoolMonitor != nil {
		opts.SetPoolMonitor(conf.PoolMonitor)
	}
	if conf.ReadPreference != nil {
		readPreference, err := newReadPref(*conf.ReadPreference)
		if err != nil {
			return nil, err
		}
		opts.SetReadPreference(readPreference)
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

// newReadPref create readpref.ReadPref from config
func newReadPref(pref ReadPref) (*readpref.ReadPref, error) {
	readPrefOpts := make([]readpref.Option, 0, 1)
	if pref.MaxStalenessMS != 0 {
		readPrefOpts = append(readPrefOpts, readpref.WithMaxStaleness(time.Duration(pref.MaxStalenessMS)*time.Millisecond))
	}
	mode := readpref.PrimaryMode
	if pref.Mode != 0 {
		mode = pref.Mode
	}
	readPreference, err := readpref.New(mode, readPrefOpts...)
	return readPreference, err
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
		return err
	}
	return nil
}

// Database create connection to database
func (c *Client) Database(name string) *Database {
	return &Database{database: c.client.Database(name)}
}

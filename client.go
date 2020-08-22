package qmgo

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"go.mongodb.org/mongo-driver/x/mongo/driver/description"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"
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
	// MinPoolSize specifies the minimum number of connections allowed in the driver's connection pool to each server. If
	// this is non-zero, each server's pool will be maintained in the background to ensure that the size does not fall below
	// the minimum. This can also be set through the "minPoolSize" URI option (e.g. "minPoolSize=100"). The default is 0.
	MinPoolSize *uint64 `json:"minPoolSize"`
	// SocketTimeoutMS specifies how long the driver will wait for a socket read or write to return before returning a
	// network error. If this is 0 meaning no timeout is used and socket operations can block indefinitely.
	// The default is 300,000 ms.
	SocketTimeoutMS *int64 `json:"socketTimeoutMS"`
	// ReadPreference determines which servers are considered suitable for read operations.
	// default is PrimaryMode
	ReadPreference *ReadPref `json:"readPreference"`
	// can be used to provide authentication options when configuring a Client.
	Auth *Credential `json:"auth"`
	// PoolMonitor to receive connection pool events
	PoolMonitor *event.PoolMonitor
}

// Credential can be used to provide authentication options when configuring a Client.
//
// AuthMechanism: the mechanism to use for authentication. Supported values include "SCRAM-SHA-256", "SCRAM-SHA-1",
// "MONGODB-CR", "PLAIN", "GSSAPI", "MONGODB-X509", and "MONGODB-AWS". This can also be set through the "authMechanism"
// URI option. (e.g. "authMechanism=PLAIN"). For more information, see
// https://docs.mongodb.com/manual/core/authentication-mechanisms/.
// AuthSource: the name of the database to use for authentication. This defaults to "$external" for MONGODB-X509,
// GSSAPI, and PLAIN and "admin" for all other mechanisms. This can also be set through the "authSource" URI option
// (e.g. "authSource=otherDb").
//
// Username: the username for authentication. This can also be set through the URI as a username:password pair before
// the first @ character. For example, a URI for user "user", password "pwd", and host "localhost:27017" would be
// "mongodb://user:pwd@localhost:27017". This is optional for X509 authentication and will be extracted from the
// client certificate if not specified.
//
// Password: the password for authentication. This must not be specified for X509 and is optional for GSSAPI
// authentication.
//
// PasswordSet: For GSSAPI, this must be true if a password is specified, even if the password is the empty string, and
// false if no password is specified, indicating that the password should be taken from the context of the running
// process. For other mechanisms, this field is ignored.
type Credential struct {
	AuthMechanism string `json:"authMechanism"`
	AuthSource    string `json:"authSource"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	PasswordSet   bool   `json:"passwordSet"`
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
// QmgoClient can operates all qmgo.client 、qmgo.database and qmgo.collection
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
	conf   Config
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
		conf:   *conf,
	}
	return
}

// client creates connection to mongo
func client(ctx context.Context, conf *Config) (client *mongo.Client, err error) {
	opts, err := newConnectOpts(conf)
	if err != nil {
		return nil, err
	}
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

// newConnectOpts creates client options from conf
// Qmgo will follow this way official mongodb driver do：
// - the configuration in uri takes precedence over the configuration in the setter
// - Check the validity of the configuration in the uri, while the configuration in the setter is basically not checked
func newConnectOpts(conf *Config) (*options.ClientOptions, error) {
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
	if conf.MinPoolSize != nil {
		opts.SetMinPoolSize(*conf.MinPoolSize)
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
	if conf.Auth != nil {
		auth, err := newAuth(*conf.Auth)
		if err != nil {
			return nil, err
		}
		opts.SetAuth(auth)
	}
	opts.ApplyURI(conf.Uri)
	return opts, nil
}

// newAuth create options.Credential from conf.Auth
func newAuth(auth Credential) (credential options.Credential, err error) {
	if auth.AuthMechanism != "" {
		credential.AuthMechanism = auth.AuthMechanism
	}
	if auth.AuthSource != "" {
		credential.AuthSource = auth.AuthSource
	}
	if auth.Username != "" {
		// Validate and process the username.
		if strings.Contains(auth.Username, "/") {
			err = ErrNotSupportedUsername
			return
		}
		credential.Username, err = url.QueryUnescape(auth.Username)
		if err != nil {
			err = ErrNotSupportedUsername
			return
		}
	}
	credential.PasswordSet = auth.PasswordSet
	if auth.Password != "" {
		if strings.Contains(auth.Password, ":") {
			err = ErrNotSupportedPassword
			return
		}
		if strings.Contains(auth.Password, "/") {
			err = ErrNotSupportedPassword
			return
		}
		credential.Password, err = url.QueryUnescape(auth.Password)
		if err != nil {
			err = ErrNotSupportedPassword
			return
		}
		credential.Password = auth.Password
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

// Session create one session on client
// Watch out, close session after operation done
func (c *Client) Session() (*Session, error) {
	s, err := c.client.StartSession()
	return &Session{session: s}, err
}

// DoTransaction do whole transaction in one function
// precondition：
// - version of mongoDB server >= v4.0
// - Topology of mongoDB server is not Single
// At the same time, please pay attention to the following
// - make sure all operations in callback use the sessCtx as context parameter
// - if operations in callback takes more than(include equal) 120s, the operations will not take effect,
// - if operation in callback return qmgo.ErrTransactionRetry,
//   the whole transaction will retry, so this transaction must be idempotent
// - if operations in callback return qmgo.ErrTransactionNotSupported,
// - If the ctx parameter already has a Session attached to it, it will be replaced by this session.
func (c *Client) DoTransaction(ctx context.Context, callback func(sessCtx context.Context) (interface{}, error)) (interface{}, error) {
	if !c.transactionAllowed() {
		return nil, ErrTransactionNotSupported
	}
	s, err := c.Session()
	if err != nil {
		return nil, err
	}
	defer s.EndSession(ctx)
	return s.StartTransaction(ctx, callback)
}

// ServerVersion get the version of mongoDB server, like 4.4.0
func (c *Client) ServerVersion() string {
	var serverStatus bson.Raw
	err := c.client.Database("admin").RunCommand(
		context.Background(),
		bson.D{{"serverStatus", 1}},
	).Decode(&serverStatus)
	if err != nil {
		fmt.Println("run command err", err)
		return ""
	}
	v, err := serverStatus.LookupErr("version")
	if err != nil {
		fmt.Println("look up err", err)
		return ""
	}
	return v.StringValue()
}

// topology get topology from mongoDB server, like single、ReplicaSet
// TODO dont know why need to do `cli, err := Open(ctx, &c.conf)` to get topo,
// Before figure it out, we only use this function in UT
func (c *Client) topology() (description.TopologyKind, error) {
	ConnString, err := connstring.ParseAndValidate(c.conf.Uri)
	if err != nil {
		return 0, err
	}
	topo, err := topology.New(
		topology.WithConnString(func(connstring.ConnString) connstring.ConnString {
			return ConnString
		}))
	err = topo.Connect()

	ctx := context.Background()
	conf := c.conf
	var maxPoolSize uint64 = 2
	var minPoolSize uint64 = 0
	conf.MaxPoolSize = &maxPoolSize
	conf.MinPoolSize = &minPoolSize
	// I dont get why need on more connect to make topo valid....
	cli, err := Open(ctx, &c.conf)
	defer cli.Close(ctx)
	cli.ServerVersion()
	return topo.Kind(), nil
}

// transactionAllowed check if transaction is allowed
func (c *Client) transactionAllowed() bool {
	vr, err := CompareVersions("4.0", c.ServerVersion())
	if err != nil {
		return false
	}
	if vr > 0 {
		fmt.Println("transaction is not supported because mongo server version is below 4.0")
		return false
	}
	// TODO dont know why need to do `cli, err := Open(ctx, &c.conf)` in topology() to get topo,
	// Before figure it out, we only use this function in UT
	//topo, err := c.topology()
	//if topo == description.Single {
	//	fmt.Println("transaction is not supported because mongo server topology is single")
	//	return false
	//}
	return true
}

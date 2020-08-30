package options

import (
	"crypto/tls"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// ClientOption operation function to set the configure a Client instance.
type ClientOption func(*options.ClientOptions)

// SetAppName specifies an application name that is sent to the server when creating new connections. It is used by the
// server to log connection and profiling information (e.g. slow query logs). This can also be set through the "appName"
// URI option (e.g "appName=example_application"). The default is empty, meaning no app name will be sent.
func SetAppName(s string) ClientOption {
	return func(opt *options.ClientOptions) {
		opt.SetAppName(s)
	}
}

// SetAuth specifies a Credential containing options for configuring authentication. See the options.Credential
// documentation for more information about Credential fields. The default is an empty Credential, meaning no
// authentication will be configured.
func SetAuth(auth options.Credential) ClientOption {
	return func(opt *options.ClientOptions) {
		opt.SetAuth(auth)
	}
}

// SetCompressors sets the compressors that can be used when communicating with a server. Valid values are:
//
// 1. "snappy" - requires server version >= 3.4
//
// 2. "zlib" - requires server version >= 3.6
//
// 3. "zstd" - requires server version >= 4.2, and driver version >= 1.2.0 with cgo support enabled or driver version >= 1.3.0
//    without cgo
//
// To use compression, it must be enabled on the server as well. If this option is specified, the driver will perform a
// negotiation with the server to determine a common list of of compressors and will use the first one in that list when
// performing operations. See
// https://docs.mongodb.com/manual/reference/program/mongod/#cmdoption-mongod-networkmessagecompressors for more
// information about how to enable this feature on the server.
//
// This can also be set through the "compressors" URI option (e.g. "compressors=zstd,zlib,snappy"). The default is
// an empty slice, meaning no compression will be enabled.
func SetCompressors(comps []string) ClientOption {
	return func(opt *options.ClientOptions) {
		opt.SetCompressors(comps)
	}
}

// SetConnectTimeout specifies a timeout that is used for creating connections to the server. If a custom Dialer is
// specified through SetDialer, this option must not be used. This can be set through ApplyURI with the
// "connectTimeoutMS" (e.g "connectTimeoutMS=30") option. If set to 0, no timeout will be used. The default is 30
// seconds.
func SetConnectTimeout(d time.Duration) ClientOption {
	return func(opt *options.ClientOptions) {
		opt.SetConnectTimeout(d)
	}
}

// SetDialer specifies a custom ContextDialer to be used to create new connections to the server. The default is a
// net.Dialer instance with a 300 second keepalive time.
func SetDialer(d options.ContextDialer) ClientOption {
	return func(opt *options.ClientOptions) {
		opt.SetDialer(d)
	}
}

// SetDirect specifies whether or not a direct connect should be made. To use this option, a URI with a single host must
// be specified through ApplyURI. If set to true, the driver will only connect to the host provided in the URI and will
// not discover other hosts in the cluster. This can also be set through the "connect" URI option with the following
// values:
//
// 1. "connect=direct" for direct connections
//
// 2. "connect=automatic" for automatic discovery.
//
// The default is false ("automatic" in the connection string).
func SetDirect(b bool) ClientOption {
	return func(opt *options.ClientOptions) {
		opt.SetDirect(b)
	}
}

// SetHeartbeatInterval specifies the amount of time to wait between periodic background server checks. This can also be
// set through the "heartbeatIntervalMS" URI option (e.g. "heartbeatIntervalMS=10000"). The default is 10 seconds.
func SetHeartbeatInterval(d time.Duration) ClientOption {
	return func(opt *options.ClientOptions) {
		opt.SetHeartbeatInterval(d)
	}
}

// SetHosts specifies a list of host names or IP addresses for servers in a cluster. Both IPv4 and IPv6 addresses are
// supported. IPv6 literals must be enclosed in '[]' following RFC-2732 syntax.
//
// Hosts can also be specified as a comma-separated list in a URI. For example, to include "localhost:27017" and
// "localhost:27018", a URI could be "mongodb://localhost:27017,localhost:27018". The default is ["localhost:27017"]
func SetHosts(s []string) ClientOption {
	return func(opt *options.ClientOptions) {
		opt.SetHosts(s)
	}
}

// SetLocalThreshold specifies the width of the 'latency window': when choosing between multiple suitable servers for an
// operation, this is the acceptable non-negative delta between shortest and longest average round-trip times. A server
// within the latency window is selected randomly. This can also be set through the "localThresholdMS" URI option (e.g.
// "localThresholdMS=15000"). The default is 15 milliseconds.
func SetLocalThreshold(d time.Duration) ClientOption {
	return func(opt *options.ClientOptions) {
		opt.SetLocalThreshold(d)
	}
}

// SetMaxConnIdleTime specifies the maximum amount of time that a connection will remain idle in a connection pool
// before it is removed from the pool and closed. This can also be set through the "maxIdleTimeMS" URI option (e.g.
// "maxIdleTimeMS=10000"). The default is 0, meaning a connection can remain unused indefinitely.
func SetMaxConnIdleTime(d time.Duration) ClientOption {
	return func(opt *options.ClientOptions) {
		opt.SetMaxConnIdleTime(d)
	}
}

// SetMaxPoolSize specifies that maximum number of connections allowed in the driver's connection pool to each server.
// Requests to a server will block if this maximum is reached. This can also be set through the "maxPoolSize" URI option
// (e.g. "maxPoolSize=100"). The default is 100. If this is 0, it will be set to math.MaxInt64.
func SetMaxPoolSize(u uint64) ClientOption {
	return func(opt *options.ClientOptions) {
		opt.SetMaxPoolSize(u)
	}
}

// SetMinPoolSize specifies the minimum number of connections allowed in the driver's connection pool to each server. If
// this is non-zero, each server's pool will be maintained in the background to ensure that the size does not fall below
// the minimum. This can also be set through the "minPoolSize" URI option (e.g. "minPoolSize=100"). The default is 0.
func SetMinPoolSize(u uint64) ClientOption {
	return func(opt *options.ClientOptions) {
		opt.SetMinPoolSize(u)
	}
}

// SetPoolMonitor specifies a PoolMonitor to receive connection pool events. See the event.PoolMonitor documentation
// for more information about the structure of the monitor and events that can be received.
func SetPoolMonitor(m *event.PoolMonitor) ClientOption {
	return func(opt *options.ClientOptions) {
		opt.SetPoolMonitor(m)
	}
}

// SetMonitor specifies a CommandMonitor to receive command events. See the event.CommandMonitor documentation for more
// information about the structure of the monitor and events that can be received.
func SetMonitor(m *event.CommandMonitor) ClientOption {
	return func(opt *options.ClientOptions) {
		opt.SetMonitor(m)
	}
}

// SetReadConcern specifies the read concern to use for read operations. A read concern level can also be set through
// the "readConcernLevel" URI option (e.g. "readConcernLevel=majority"). The default is nil, meaning the server will use
// its configured default.
func SetReadConcern(rc *readconcern.ReadConcern) ClientOption {
	return func(opt *options.ClientOptions) {
		opt.SetReadConcern(rc)
	}
}

// SetReadPreference specifies the read preference to use for read operations. This can also be set through the
// following URI options:
//
// 1. "readPreference" - Specifiy the read preference mode (e.g. "readPreference=primary").
//
// 2. "readPreferenceTags": Specify one or more read preference tags
// (e.g. "readPreferenceTags=region:south,datacenter:A").
//
// 3. "maxStalenessSeconds" (or "maxStaleness"): Specify a maximum replication lag for reads from secondaries in a
// replica set (e.g. "maxStalenessSeconds=10").
//
// The default is readpref.Primary(). See https://docs.mongodb.com/manual/core/read-preference/#read-preference for
// more information about read preferences.
func SetReadPreference(rp *readpref.ReadPref) ClientOption {
	return func(opt *options.ClientOptions) {
		opt.SetReadPreference(rp)
	}
}

// SetRegistry specifies the BSON registry to use for BSON marshalling/unmarshalling operations. The default is
// bson.DefaultRegistry.
func SetRegistry(registry *bsoncodec.Registry) ClientOption {
	return func(opt *options.ClientOptions) {
		opt.SetRegistry(registry)
	}
}

// SetReplicaSet specifies the replica set name for the cluster. If specified, the cluster will be treated as a replica
// set and the driver will automatically discover all servers in the set, starting with the nodes specified through
// ApplyURI or SetHosts. All nodes in the replica set must have the same replica set name, or they will not be
// considered as part of the set by the Client. This can also be set through the "replicaSet" URI option (e.g.
// "replicaSet=replset"). The default is empty.
func SetReplicaSet(s string) ClientOption {
	return func(opt *options.ClientOptions) {
		opt.SetReplicaSet(s)
	}
}

// SetRetryWrites specifies whether supported write operations should be retried once on certain errors, such as network
// errors.
//
// Supported operations are InsertOne, UpdateOne, ReplaceOne, DeleteOne, FindOneAndDelete, FindOneAndReplace,
// FindOneAndDelete, InsertMany, and BulkWrite. Note that BulkWrite requests must not include UpdateManyModel or
// DeleteManyModel instances to be considered retryable. Unacknowledged writes will not be retried, even if this option
// is set to true.
//
// This option requires server version >= 3.6 and a replica set or sharded cluster and will be ignored for any other
// cluster type. This can also be set through the "retryWrites" URI option (e.g. "retryWrites=true"). The default is
// true.
func SetRetryWrites(b bool) ClientOption {
	return func(opt *options.ClientOptions) {
		opt.SetRetryWrites(b)
	}
}

// SetServerSelectionTimeout specifies how long the driver will wait to find an available, suitable server to execute an
// operation. This can also be set through the "serverSelectionTimeoutMS" URI option (e.g.
// "serverSelectionTimeoutMS=30000"). The default value is 30 seconds.
func SetServerSelectionTimeout(d time.Duration) ClientOption {
	return func(opt *options.ClientOptions) {
		opt.SetServerSelectionTimeout(d)
	}
}

// SetSocketTimeout specifies how long the driver will wait for a socket read or write to return before returning a
// network error. This can also be set through the "socketTimeoutMS" URI option (e.g. "socketTimeoutMS=1000"). The
// default value is 0, meaning no timeout is used and socket operations can block indefinitely.
func SetSocketTimeout(d time.Duration) ClientOption {
	return func(opt *options.ClientOptions) {
		opt.SetSocketTimeout(d)
	}
}

// SetTLSConfig specifies a tls.Config instance to use use to configure TLS on all connections created to the cluster.
// This can also be set through the following URI options:
//
// 1. "tls" (or "ssl"): Specify if TLS should be used (e.g. "tls=true").
//
// 2. Either "tlsCertificateKeyFile" (or "sslClientCertificateKeyFile") or a combination of "tlsCertificateFile" and
// "tlsPrivateKeyFile". The "tlsCertificateKeyFile" option specifies a path to the client certificate and private key,
// which must be concatenated into one file. The "tlsCertificateFile" and "tlsPrivateKey" combination specifies separate
// paths to the client certificate and private key, respectively. Note that if "tlsCertificateKeyFile" is used, the
// other two options must not be specified.
//
// 3. "tlsCertificateKeyFilePassword" (or "sslClientCertificateKeyPassword"): Specify the password to decrypt the client
// private key file (e.g. "tlsCertificateKeyFilePassword=password").
//
// 4. "tlsCaFile" (or "sslCertificateAuthorityFile"): Specify the path to a single or bundle of certificate authorities
// to be considered trusted when making a TLS connection (e.g. "tlsCaFile=/path/to/caFile").
//
// 5. "tlsInsecure" (or "sslInsecure"): Specifies whether or not certificates and hostnames received from the server
// should be validated. If true (e.g. "tlsInsecure=true"), the TLS library will accept any certificate presented by the
// server and any host name in that certificate. Note that setting this to true makes TLS susceptible to
// man-in-the-middle attacks and should only be done for testing.
//
// The default is nil, meaning no TLS will be enabled.
func SetTLSConfig(cfg *tls.Config) ClientOption {
	return func(opt *options.ClientOptions) {
		opt.SetTLSConfig(cfg)
	}
}

// SetWriteConcern specifies the write concern to use to for write operations. This can also be set through the following
// URI options:
//
// 1. "w": Specify the number of nodes in the cluster that must acknowledge write operations before the operation
// returns or "majority" to specify that a majority of the nodes must acknowledge writes. This can either be an integer
// (e.g. "w=10") or the string "majority" (e.g. "w=majority").
//
// 2. "wTimeoutMS": Specify how long write operations should wait for the correct number of nodes to acknowledge the
// operation (e.g. "wTimeoutMS=1000").
//
// 3. "journal": Specifies whether or not write operations should be written to an on-disk journal on the server before
// returning (e.g. "journal=true").
//
// The default is nil, meaning the server will use its configured default.
func SetWriteConcern(wc *writeconcern.WriteConcern) ClientOption {
	return func(opt *options.ClientOptions) {
		opt.SetWriteConcern(wc)
	}
}

// SetZlibLevel specifies the level for the zlib compressor. This option is ignored if zlib is not specified as a
// compressor through ApplyURI or SetCompressors. Supported values are -1 through 9, inclusive. -1 tells the zlib
// library to use its default, 0 means no compression, 1 means best speed, and 9 means best compression.
// This can also be set through the "zlibCompressionLevel" URI option (e.g. "zlibCompressionLevel=-1"). Defaults to -1.
func SetZlibLevel(level int) ClientOption {
	return func(opt *options.ClientOptions) {
		opt.SetZlibLevel(level)
	}
}

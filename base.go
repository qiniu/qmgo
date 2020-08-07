package qmgo

import (
	"strings"
)

// D is an ordered representation of a BSON document. This type should be used when the order of the elements matters,
// such as MongoDB command documents. If the order of the elements does not matter, an M should be used instead.
//
// Example usage:
// D{{"foo", "bar"}, {"hello", "world"}, {"pi", 3.14159}}
type D []E

// E represents a BSON element for a D. It is usually used inside a D.
type E struct {
	Key   string
	Value interface{}
}

// M is an unordered representation of a BSON document. This type should be used when the order of the elements does not
// matter. This type is handled as a regular map[string]interface{} when encoding and decoding. Elements will be
// serialized in an undefined, random order. If the order of the elements matters, a D should be used instead.
//
// Example usage:
//
// M{"foo": "bar", "hello": "world", "pi": 3.14159}.
type M map[string]interface{}

// An A is an ordered representation of a BSON array.
// Example usage:
// A{"bar", "world", 3.14159, D{{"qux", 12345}}}
type A []interface{}

// QmgoConfig for initial mongodb instance
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
}

// IsDup check if err is mongo E11000 (duplicate err)。
func IsDup(err error) bool {
	return strings.Contains(err.Error(), "E11000")
}

//// Now return Millisecond current time
//func Now() time.Time {
//	return time.Unix(0, time.Now().UnixNano()/1e6*1e6)
//}

// SplitSortField handle sort symbol: "+"/"-" in front of field
// if "+"， return sort as 1
// if "-"， return sort as -1
func SplitSortField(field string) (key string, sort int32) {
	sort = 1
	key = field

	if len(field) != 0 {
		switch field[0] {
		case '+':
			key = strings.TrimPrefix(field, "+")
			sort = 1
		case '-':
			key = strings.TrimPrefix(field, "-")
			sort = -1
		}
	}

	return key, sort
}

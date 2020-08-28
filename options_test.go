package qmgo

import (
	"context"
	"crypto/tls"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

var tClientOptions = reflect.TypeOf(&ClientOptions{})

func TestClientOptions(t *testing.T) {
	ast := require.New(t)
	testCases := []struct {
		name        string
		fn          interface{} // method to be run
		arg         interface{} // argument for method
		field       string      // field to be set
		dereference bool        // Should we compare a pointer or the field
	}{
		{"AppName", SetAppName, "example-application", "AppName", true},
		{"Auth", SetAuth, options.Credential{Username: "foo", Password: "bar"}, "Auth", true},
		{"Compressors", SetCompressors, []string{"zstd", "snappy", "zlib"}, "Compressors", true},
		{"ConnectTimeout", SetConnectTimeout, 5 * time.Second, "ConnectTimeout", true},
		{"Dialer", SetDialer, testDialer{Num: 12345}, "Dialer", true},
		{"HeartbeatInterval", SetHeartbeatInterval, 5 * time.Second, "HeartbeatInterval", true},
		{"Hosts", SetHosts, []string{"localhost:27017", "localhost:27018", "localhost:27019"}, "Hosts", true},
		{"LocalThreshold", SetLocalThreshold, 5 * time.Second, "LocalThreshold", true},
		{"MaxConnIdleTime", SetMaxConnIdleTime, 5 * time.Second, "MaxConnIdleTime", true},
		{"MaxPoolSize", SetMaxPoolSize, uint64(250), "MaxPoolSize", true},
		{"MinPoolSize", SetMinPoolSize, uint64(10), "MinPoolSize", true},
		{"PoolMonitor", SetPoolMonitor, &event.PoolMonitor{}, "PoolMonitor", false},
		{"Monitor", SetMonitor, &event.CommandMonitor{}, "Monitor", false},
		{"ReadConcern", SetReadConcern, readconcern.Majority(), "ReadConcern", false},
		{"ReadPreference", SetReadPreference, readpref.SecondaryPreferred(), "ReadPreference", false},
		{"Registry", SetRegistry, bson.NewRegistryBuilder().Build(), "Registry", false},
		{"ReplicaSet", SetReplicaSet, "example-replicaset", "ReplicaSet", true},
		{"RetryWrites", SetRetryWrites, true, "RetryWrites", true},
		{"ServerSelectionTimeout", SetServerSelectionTimeout, 5 * time.Second, "ServerSelectionTimeout", true},
		{"Direct", SetDirect, true, "Direct", true},
		{"SocketTimeout", SetSocketTimeout, 5 * time.Second, "SocketTimeout", true},
		{"TLSConfig", SetTLSConfig, &tls.Config{}, "TLSConfig", false},
		{"WriteConcern", SetWriteConcern, writeconcern.New(writeconcern.WMajority()), "WriteConcern", false},
		{"ZlibLevel", SetZlibLevel, 6, "ZlibLevel", true},
	}

	var opt = &ClientOptions{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fn := reflect.ValueOf(tc.fn)
			if fn.Kind() != reflect.Func {
				t.Fatal("fn argument must be a function")
			}

			args := make([]reflect.Value, 1)
			want := reflect.ValueOf(tc.arg)
			args[0] = want

			if !want.IsValid() || !want.CanInterface() {
				t.Fatal("arg property of test case must be valid")
			}

			if _, exists := tClientOptions.Elem().FieldByName(tc.field); !exists {
				t.Fatalf("field (%s) does not exist in ClientOptions", tc.field)
			}

			client := reflect.New(tClientOptions.Elem())
			got := client.Elem().FieldByName(tc.field)
			if !got.IsValid() || !got.CanInterface() {
				t.Fatal("cannot create concrete instance from retrieved field")
			}

			result := fn.Call(args)
			ast.Equal(1, len(result))
			optionFunc := result[0]
			var a OptionFunc = func(opt *ClientOptions) {}
			ast.Equal(reflect.TypeOf(a), optionFunc.Type())
			wantT := want.Type()
			gotT := got.Type()
			if wantT.Kind() == reflect.Ptr {
				wantT = wantT.Elem()
			}
			if gotT.Kind() == reflect.Ptr {
				gotT = gotT.Elem()
			}
			if gotT.Kind() == reflect.Interface {
				ast.Equal(true, wantT.ConvertibleTo(gotT))
			} else {
				ast.Equal(wantT, gotT)
			}

			in := []reflect.Value{}
			optVal := reflect.ValueOf(opt)
			in = append(in, optVal)
			optionFunc.Call(in)
			field := optVal.Elem().FieldByName(tc.name)
			if field.Kind() == reflect.Ptr {
				field = field.Elem()
			}
			argVal := reflect.ValueOf(tc.arg)
			if argVal.Kind() == reflect.Ptr {
				argVal = argVal.Elem()
			}
			ast.Equal(argVal.Interface(), field.Interface())
		})
	}
}

type testDialer struct {
	Num int
}

func (testDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return nil, nil
}

package qmgo

import "go.mongodb.org/mongo-driver/bson"

// alias mongo drive bson primitives
// thus user don't need to import go.mongodb.org/mongo-driver/mongo, it's all in qmgo
type (
	// M is an alias of bson.M
	M = bson.M
	// A is an alias of bson.A
	A = bson.A
	// D is an alias of bson.D
	D = bson.D
	// E is an alias of bson.E
	E = bson.E
)

package validator

import (
	"github.com/qiniu/qmgo/operator"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

// User contains user information
type User struct {
	FirstName      string     `bson:"fname"`
	LastName       string     `bson:"lname"`
	Age            uint8      `bson:"age" validate:"gte=0,lte=130"`
	Email          string     `bson:"e-mail" validate:"required,email"`
	FavouriteColor string     `bson:"favouriteColor" validate:"hexcolor|rgb|rgba"`
	Addresses      []*Address `bson:"addresses" validate:"required,dive,required"` // a person can have a home and cottage...
}

// Address houses a users address information
type Address struct {
	Street string `validate:"required"`
	City   string `validate:"required"`
	Planet string `validate:"required"`
	Phone  string `validate:"required"`
}

func TestValidator(t *testing.T) {
	ast := require.New(t)

	user := &User{}
	// not need validator op
	ast.NoError(Do(user, operator.BeforeRemove))
	ast.NoError(Do(user, operator.AfterInsert))
	// check success
	address := &Address{
		Street: "Eavesdown Docks",
		Planet: "Persphone",
		Phone:  "none",
		City:   "Unknown",
	}

	user = &User{
		FirstName:      "",
		LastName:       "",
		Age:            45,
		Email:          "1234@gmail.com",
		FavouriteColor: "#000",
		Addresses:      []*Address{address, address},
	}
	ast.NoError(Do(user, operator.BeforeInsert))
	ast.NoError(Do(user, operator.BeforeUpsert))
	ast.NoError(Do(*user, operator.BeforeUpsert))

	users := []*User{user, user, user}
	ast.NoError(Do(users, operator.BeforeInsert))

	// check failure
	user.Age = 150
	ast.Error(Do(user, operator.BeforeInsert))
	user.Age = 22
	user.Email = "1234@gmail" // invalid email
	ast.Error(Do(user, operator.BeforeInsert))
	user.Email = "1234@gmail.com"
	user.Addresses[0].City = "" // string tag use default value
	ast.Error(Do(user, operator.BeforeInsert))

	// input slice
	users = []*User{user, user, user}
	ast.Error(Do(users, operator.BeforeInsert))

	useris := []interface{}{user, user, user}
	ast.Error(Do(useris, operator.BeforeInsert))

	user.Addresses[0].City = "shanghai"
	users = []*User{user, user, user}
	ast.NoError(Do(users, operator.BeforeInsert))

	us := []User{*user, *user, *user}
	ast.NoError(Do(us, operator.BeforeInsert))
	ast.NoError(Do(&us, operator.BeforeInsert))

	// all bson type
	mdoc := []interface{}{bson.M{"name": "", "age": 12}, bson.M{"name": "", "age": 12}}
	ast.NoError(Do(mdoc, operator.BeforeInsert))
	adoc := bson.A{"Alex", "12"}
	ast.NoError(Do(adoc, operator.BeforeInsert))
	edoc := bson.E{"Alex", "12"}
	ast.NoError(Do(edoc, operator.BeforeInsert))
	ddoc := bson.D{{"foo", "bar"}, {"hello", "world"}, {"pi", 3.14159}}
	ast.NoError(Do(ddoc, operator.BeforeInsert))

	// nil ptr
	user = nil
	ast.NoError(Do(user, operator.BeforeInsert))
	ast.NoError(Do(nil, operator.BeforeInsert))
}

package validator

import (
	"context"
	"github.com/go-playground/validator/v10"
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

// CustomRule use custom rule
type CustomRule struct {
	Name string `validate:"required,foo"`
}

func TestValidator(t *testing.T) {
	ast := require.New(t)
	ctx := context.Background()

	user := &User{}
	// not need validator op
	ast.NoError(Do(ctx, user, operator.BeforeRemove))
	ast.NoError(Do(ctx, user, operator.AfterInsert))
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
	ast.NoError(Do(ctx, user, operator.BeforeInsert))
	ast.NoError(Do(ctx, user, operator.BeforeUpsert))
	ast.NoError(Do(ctx, *user, operator.BeforeUpsert))

	users := []*User{user, user, user}
	ast.NoError(Do(ctx, users, operator.BeforeInsert))

	// check failure
	user.Age = 150
	ast.Error(Do(ctx, user, operator.BeforeInsert))
	user.Age = 22
	user.Email = "1234@gmail" // invalid email
	ast.Error(Do(ctx, user, operator.BeforeInsert))
	user.Email = "1234@gmail.com"
	user.Addresses[0].City = "" // string tag use default value
	ast.Error(Do(ctx, user, operator.BeforeInsert))

	// input slice
	users = []*User{user, user, user}
	ast.Error(Do(ctx, users, operator.BeforeInsert))

	useris := []interface{}{user, user, user}
	ast.Error(Do(ctx, useris, operator.BeforeInsert))

	user.Addresses[0].City = "shanghai"
	users = []*User{user, user, user}
	ast.NoError(Do(ctx, users, operator.BeforeInsert))

	us := []User{*user, *user, *user}
	ast.NoError(Do(ctx, us, operator.BeforeInsert))
	ast.NoError(Do(ctx, &us, operator.BeforeInsert))

	// all bson type
	mdoc := []interface{}{bson.M{"name": "", "age": 12}, bson.M{"name": "", "age": 12}}
	ast.NoError(Do(ctx, mdoc, operator.BeforeInsert))
	adoc := bson.A{"Alex", "12"}
	ast.NoError(Do(ctx, adoc, operator.BeforeInsert))
	edoc := bson.E{"Alex", "12"}
	ast.NoError(Do(ctx, edoc, operator.BeforeInsert))
	ddoc := bson.D{{"foo", "bar"}, {"hello", "world"}, {"pi", 3.14159}}
	ast.NoError(Do(ctx, ddoc, operator.BeforeInsert))

	// nil ptr
	user = nil
	ast.NoError(Do(ctx, user, operator.BeforeInsert))
	ast.NoError(Do(ctx, nil, operator.BeforeInsert))

	// use custom rules
	customRule := &CustomRule{Name: "bar"}
	v := validator.New()
	_ = v.RegisterValidation("foo", func(fl validator.FieldLevel) bool {
		return fl.Field().String() == "bar"
	})
	SetValidate(v)
	ast.NoError(Do(ctx, customRule, operator.BeforeInsert))
}

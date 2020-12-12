package validator

import (
	"github.com/qiniu/qmgo/operator"
	"github.com/stretchr/testify/require"
	"testing"
)

// User contains user information
type User struct {
	FirstName      string     `json:"fname"`
	LastName       string     `json:"lname"`
	Age            uint8      `validate:"gte=0,lte=130"`
	Email          string     `json:"e-mail" validate:"required,email"`
	FavouriteColor string     `validate:"hexcolor|rgb|rgba"`
	Addresses      []*Address `validate:"required,dive,required"` // a person can have a home and cottage...
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

	users := []*User{user, user, user}
	ast.NoError(Do(users, operator.BeforeInsert))

	userss := [][]*User{[]*User{user, user, user}, []*User{user, user}}
	ast.NoError(Do(userss, operator.BeforeInsert))

	// check failure
	user.Age = 150
	ast.Error(Do(user, operator.BeforeInsert))
	user.Age = 22
	user.Email = "1234@gmail" // invalid email
	ast.Error(Do(user, operator.BeforeInsert))
	user.Email = "1234@gmail.com"
	user.Addresses[0].City = "" // string tag use default value
	ast.Error(Do(user, operator.BeforeInsert))

	users = []*User{user, user, user}
	ast.Error(Do(user, operator.BeforeInsert))

}

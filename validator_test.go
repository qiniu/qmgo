package qmgo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"testing"

	"github.com/stretchr/testify/require"
)

// User contains user information
type User struct {
	FirstName string `json:"fname"`
	LastName  string `json:"lname"`
	Age       uint8  `validate:"gte=0,lte=130"`
	Email     string `json:"e-mail" validate:"required,email"`
}

func TestValidator(t *testing.T) {
	ast := require.New(t)
	cli := initClient("test")
	ctx := context.Background()
	defer cli.Close(ctx)
	defer cli.DropCollection(ctx)

	user := &User{
		FirstName: "",
		LastName:  "",
		Age:       45,
		Email:     "1234@gmail.com",
	}
	_, err := cli.InsertOne(ctx, user)
	ast.NoError(err)

	user.Age = 200 // invalid age
	_, err = cli.InsertOne(ctx, user)
	ast.Error(err)

	users := []*User{user, user, user}
	_, err = cli.InsertMany(ctx, users)
	ast.Error(err)

	user.Age = 20
	user.Email = "1234@gmail" // invalid email
	err = cli.ReplaceOne(ctx, bson.M{"age": 45}, user)
	ast.Error(err)

	user.Email = "" // invalid empty email
	_, err = cli.Upsert(ctx, bson.M{"age": 45}, user)
	ast.Error(err)

}

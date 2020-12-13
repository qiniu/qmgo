package qmgo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// User contains user information
type User struct {
	FirstName string `json:"fname"`
	LastName  string `json:"lname"`
	// Age must in [0,130]
	Age uint8 `validate:"gte=0,lte=130"`
	//  Email can't be empty string, and must has email format
	Email string `json:"e-mail" validate:"required,email"`
	// CreateAt must lte than current time
	CreateAt time.Time `json:"createAt" validate:"lte"`
	// Relations can't has more than 2 elements
	Relations map[string]string `json:"relations" validate:"max=2"`
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
	user.Email = "1234@gmail" // email tag, invalid email
	err = cli.ReplaceOne(ctx, bson.M{"age": 45}, user)
	ast.Error(err)

	user.Email = "" // required tag, invalid empty string
	_, err = cli.Upsert(ctx, bson.M{"age": 45}, user)
	ast.Error(err)

	user.Email = "1234@gmail.com"
	user.CreateAt = time.Now().Add(1 * time.Hour) // lte tag for time, time must lte current time
	_, err = cli.Upsert(ctx, bson.M{"age": 45}, user)
	ast.Error(err)

	user.CreateAt = time.Now()
	user.Relations = map[string]string{"Alex": "friend", "Joe": "friend"}
	_, err = cli.Upsert(ctx, bson.M{"age": 45}, user)
	ast.NoError(err)

	user.Relations = map[string]string{"Alex": "friend", "Joe": "friend", "Bob": "sister"} // max tag, numbers of map
	_, err = cli.Upsert(ctx, bson.M{"age": 45}, user)
	ast.Error(err)
}

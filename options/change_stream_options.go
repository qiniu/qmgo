package options

import "go.mongodb.org/mongo-driver/mongo/options"

type ChangeStreamOptions struct {
	*options.ChangeStreamOptions
}
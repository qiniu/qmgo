package options

import "go.mongodb.org/mongo-driver/mongo/options"

type ReplaceOptions struct {
	UpdateHook interface{}
	*options.ReplaceOptions
}

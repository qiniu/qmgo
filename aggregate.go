package qmgo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Pipeline []D

type Aggregate struct {
	ctx        context.Context
	pipeline   interface{}
	collection *mongo.Collection
}

func (a *Aggregate) All(results interface{}) error {
	c, err := a.collection.Aggregate(a.ctx, a.pipeline)
	if err != nil {
		return err
	}
	return c.All(a.ctx, results)
}

func (a *Aggregate) One(result interface{}) error {
	c, err := a.collection.Aggregate(a.ctx, a.pipeline)
	if err != nil {
		return err
	}
	cr := Cursor{
		ctx:    a.ctx,
		cursor: c,
		err:    err,
	}
	defer cr.Close()
	if !cr.Next(result) {
		return ERR_NO_SUCH_RECORD
	}
	return err
}

func (a *Aggregate) Iter() CursorI {
	c, err := a.collection.Aggregate(a.ctx, a.pipeline)
	return &Cursor{
		ctx:    a.ctx,
		cursor: c,
		err:    err,
	}
}

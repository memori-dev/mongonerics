package mongonerics

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Client[T any] struct {
	Collection *mongo.Collection
}

type Options struct {
	GetCtx func() (context.Context, context.CancelFunc)
}

func defaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func defaultOptions() *Options {
	return &Options{
		GetCtx: defaultContext,
	}
}

func (this *Client[T]) Create(doc *T, opts ...func(*Options)) (*mongo.InsertOneResult, error) {
	options := defaultOptions()
	for i := 0; i < len(opts); i++ {
		opts[i](options)
	}

	ctx, cancel := options.GetCtx()
	defer cancel()

	return this.Collection.InsertOne(ctx, doc)
}

func (this *Client[T]) Read(filter interface{}, opts ...func(*Options)) ([]*T, error) {
	options := defaultOptions()
	for i := 0; i < len(opts); i++ {
		opts[i](options)
	}

	ctx, cancel := options.GetCtx()
	defer cancel()

	cur, err := this.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var docs []*T
	if err := cur.All(ctx, &docs); err != nil {
		return nil, err
	}

	return docs, nil
}

func (this *Client[T]) ReadOne(filter interface{}, opts ...func(*Options)) (*T, error) {
	options := defaultOptions()
	for i := 0; i < len(opts); i++ {
		opts[i](options)
	}

	ctx, cancel := options.GetCtx()
	defer cancel()

	res := this.Collection.FindOne(ctx, filter)
	if res.Err() != nil {
		return nil, res.Err()
	}

	doc := new(T)
	if err := res.Decode(doc); err != nil {
		return nil, err
	}

	return doc, nil
}

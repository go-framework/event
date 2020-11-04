package inapp

import (
	"context"
)

type subOptionCtxKey struct{}

// Set SubscribeOption into context.
func NewSubscribeOptionContext(ctx context.Context, opt ...SubscribeOption) context.Context {
	return context.WithValue(ctx, subOptionCtxKey{}, opt)
}

// Get SubscribeOption from context.
func GetSubscribeOptionFromContext(ctx context.Context) ([]SubscribeOption, bool) {
	opts, ok := ctx.Value(subOptionCtxKey{}).([]SubscribeOption)
	return opts, ok
}

// Get SubscribeOptions from context, when not exist return nil.
func GetSubscribeOptionsFromContext(ctx context.Context) *SubscribeOptions {
	opts, ok := ctx.Value(subOptionCtxKey{}).([]SubscribeOption)
	if !ok {
		return nil
	}
	options := GetDefaultSubscribeOptions()
	for _, opt := range opts {
		opt(options)
	}
	return options
}

type pubOptionCtxKey struct{}

// Set PublishOption into context.
func NewPublishOptionContext(ctx context.Context, opt ...PublishOption) context.Context {
	return context.WithValue(ctx, pubOptionCtxKey{}, opt)
}

// Get PublishOption from context.
func GetPublishOptionFromContext(ctx context.Context) ([]PublishOption, bool) {
	opts, ok := ctx.Value(pubOptionCtxKey{}).([]PublishOption)
	return opts, ok
}

// Get PublishOptions from context, when not exist return default value.
func GetPublishOptionsFromContext(ctx context.Context) *PublishOptions {
	options := GetDefaultPublishOptions()
	opts, ok := ctx.Value(pubOptionCtxKey{}).([]PublishOption)
	if !ok {
		return options
	}
	for _, opt := range opts {
		opt(options)
	}
	return options
}

type dataCtxKey struct{}

// Set data into context, it's use for passed data out by context.
func NewDataContext(ctx context.Context, data interface{}) context.Context {
	return context.WithValue(ctx, dataCtxKey{}, data)
}

// Get data from context, it's use for got data out with context.
func GetDataFromContext(ctx context.Context) (interface{}, bool) {
	data, ok := ctx.Value(dataCtxKey{}).(interface{})
	return data, ok
}

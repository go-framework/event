package inapp

import (
	"context"
)

// Default Event.
var DefaultEvent = NewEvent()

func Subscribe(ctx context.Context, event string, callback func(context.Context, ...interface{}) error) {
	DefaultEvent.Subscribe(ctx, event, callback)
}

func Publish(ctx context.Context, event string, args ...interface{}) error {
	return DefaultEvent.Publish(ctx, event, args...)
}

func Unsubscribe(event string, callback ...func(context.Context, ...interface{}) error) {
	DefaultEvent.Unsubscribe(event, callback...)
}

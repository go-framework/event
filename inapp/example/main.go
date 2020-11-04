package main

import (
	"context"
	"fmt"

	"github.com/go-framework/event/inapp"
)

func main() {
	// new inapp Event
	var event = inapp.NewEvent()

	// callback
	var f1 = func(ctx context.Context, args ...interface{}) error {
		fmt.Sprintf("got args %v\n", args)
		return nil
	}

	// Subscribe test event
	event.Subscribe(context.TODO(), "test", f1)

	// Subscribe test event with Once option
	event.Subscribe(inapp.NewSubscribeOptionContext(context.TODO(), inapp.WithOnceOption(true)), "test", func(ctx context.Context, args ...interface{}) error {
		fmt.Sprintf("got args %v\n", args)
		return nil
	})

	// Publish to test event
	event.Publish(context.TODO(), "test", "i'am a arg")

	var errCh = make(chan error)
	// Publish to test event with callback done
	event.Publish(inapp.NewPublishOptionContext(context.TODO(), inapp.WithErrorOption(errCh)), "test", "i'am a arg")
	select {
	case err := <-errCh:
		if err != nil {
			fmt.Printf("got error = %v\n", err)
			return
		}
		fmt.Printf("succeed\n")
	}

	// Publish to test event with callback done in Strict mode
	event.Publish(inapp.NewPublishOptionContext(context.TODO(), inapp.WithErrorOption(errCh), inapp.WithStrictModeOption(true)), "test", "i'am a arg")
	select {
	case err := <-errCh:
		if err != nil {
			fmt.Printf("got error = %v\n", err)
			return
		}
		fmt.Printf("succeed\n")
	}

	var count = 0
	// Publish to test event with context value and callback done in Strict mode
	event.Publish(inapp.NewPublishOptionContext(inapp.NewDataContext(context.TODO(), &count), inapp.WithErrorOption(errCh), inapp.WithStrictModeOption(true)), "test", "i'am a arg")
	select {
	case err := <-errCh:
		if err != nil {
			fmt.Printf("got error = %v\n", err)
			return
		}
		fmt.Printf("update count = %d\n", count)
	}

	// Unsubscribe test even by callback func f1
	event.Unsubscribe("test", f1)

	// Unsubscribe all test even
	event.Unsubscribe("test")
}

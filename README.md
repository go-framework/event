# Event

Event is an event notify written in Go (Golang). support Subscribe/Publish/Unsubscribe.

## Installation

To install Event package, you need to install Go and set your Go workspace first.

The first need Go installed (version 1.12+ is required), then you can use the below Go command to install Event.

```shell script
$ go get -u github.com/go-framework/event
```

## Support

### [InApp](https://github.com/go-framework/event/tree/master/inapp)

InApp event is an in application notify, support Subscribe/Publish/Unsubscribe.

Import it in your code:

```go
import "github.com/go-framework/event/inapp"
```

1. New InApp event.
    ```go
    // new inapp Event
    var event = inapp.NewEvent()
    ```

2. Subscribe an event. 
    - Normal
    - OnceOption: callback will be removed when pass with Once option after done. 

    ```go
    // callback func
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
    ```
   
3. Publish to event.
    - Normal
    - ErrorOption: use an error chan got the callback done returns
    - StrictMode: it will interrupt and return when callback return error
    - ContextValue
    
    ```go
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
    ```

4. Unsubscribe
    - Unsubscribe a specified callback with param
    - Unsubscribe all event callback with nil param
    
    ```go
    // Unsubscribe test even by callback func f1
    event.Unsubscribe("test", f1)
    
    // Unsubscribe all test even
    event.Unsubscribe("test")
    ```

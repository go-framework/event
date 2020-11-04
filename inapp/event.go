package inapp

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var (
	ErrNotExistEvent = errors.New("event not exist")
)

// Event is a inapp name. subscribe name into inbox, when publish added to list.
type Event struct {
	list sync.Map // the active event list. map[string]*event
}

// New Event.
func NewEvent() *Event {
	return new(Event)
}

// Subscribe event with name and callback func f, passed option by context.
func (e *Event) Subscribe(ctx context.Context, name string, f func(context.Context, ...interface{}) error) {
	if f == nil {
		return
	}
	cb := &callback{
		f:                f,
		subscribeOptions: GetSubscribeOptionsFromContext(ctx),
	}

	actual, ok := e.list.LoadOrStore(name, &event{
		doneLock:  make(chan struct{}, 1),
		callbacks: callbacks{cb},
	})

	var event = actual.(*event)

	if !ok {
		event.doneLock <- struct{}{}
		return
	}

	event.mu.Lock()
	// mutex with Unsubscribe
	// if exist the same f then replace the latest callback.
	if !event.callbacks.replace(cb) {
		event.callbacks = append(event.callbacks, cb) // append
	}
	if event.doneLock == nil {
		event.doneLock = make(chan struct{}, 1)
		event.doneLock <- struct{}{}
	}
	e.list.LoadOrStore(name, event)
	event.mu.Unlock()
}

// Publish event with args and publish option by context to async done callbacks, will be remove Once subscribed.
func (e *Event) Publish(ctx context.Context, name string, args ...interface{}) error {
	actual, ok := e.list.Load(name)
	if !ok {
		return ErrNotExistEvent
	}

	var event = actual.(*event)

	// callback done
	done := func(ctx context.Context, args ...interface{}) {
		var publishOptions = GetPublishOptionsFromContext(ctx)
		var err error

		defer func() {
			if e := recover(); e != nil {
				switch v := e.(type) {
				case error:
					err = v
				default:
					err = fmt.Errorf("%v", e)
				}
			}

			if publishOptions.Err != nil {
				publishOptions.Err <- err
			}
		}()

		defer func() {
			event.mu.Lock()
			// mutex with Subscribe
			event.callbacks = event.callbacks.clearRemoveFlags()
			if len(event.callbacks) == 0 {
				close(event.doneLock)
				event.doneLock = nil
				e.list.Delete(name)
			}
			event.mu.Unlock()
		}()

		var errs = make(Errors, 0)
		<-event.doneLock
		for i := 0; i < len(event.callbacks); i++ {
			// once subscribe set remove flag
			if event.callbacks[i].subscribeOptions != nil && event.callbacks[i].subscribeOptions.Once {
				event.callbacks[i].remove = true
			}
			if event.callbacks[i].f == nil {
				continue
			}
			// exec f
			err = func() (_err error) {
				defer func() {
					if e := recover(); e != nil {
						switch v := e.(type) {
						case error:
							_err = v
						default:
							_err = fmt.Errorf("%v", e)
						}
					}
				}()
				return event.callbacks[i].f(ctx, args...)
			}()
			if err != nil {
				// strict mode
				if publishOptions.Strict {
					return
				}
				errs = append(errs, err)
			}
		}
		event.doneLock <- struct{}{}
		err = errs.Nil()
	}

	// done
	go done(ctx, args...)

	return nil
}

// Unsubscribe event with callback func list, remove all event when func list is ignore.
func (e *Event) Unsubscribe(name string, f ...func(context.Context, ...interface{}) error) {
	actual, ok := e.list.Load(name)
	if !ok {
		return
	}

	var event = actual.(*event)

	select {
	case <-event.doneLock: // not in Publish progress
		event.mu.Lock()
		// mutex with Subscribe
		event.callbacks = event.callbacks.remove(f...)
		if len(event.callbacks) == 0 {
			close(event.doneLock)
			event.doneLock = nil
			e.list.Delete(name)
		} else {
			event.doneLock <- struct{}{}
		}
		event.mu.Unlock()
	default:
		if len(f) == 0 {
			event.mu.Lock()
			// mutex with Subscribe
			event.callbacks = event.callbacks.markRemoveAll()
			event.mu.Unlock()
		} else {
			event.callbacks = event.callbacks.markRemove(f...)
		}
	}
}

// event case.
type event struct {
	callbacks callbacks     // name callback list
	mu        sync.Mutex    // mu protects callback list.
	doneLock  chan struct{} // doneLock has a one-element buffer and is empty when held, it protects at callbacks reduce.
}

// callback list.
type callbacks []*callback

// event callback.
type callback struct {
	f                func(context.Context, ...interface{}) error
	remove           bool // remove flag for remove when publish.
	subscribeOptions *SubscribeOptions
}

func (list *callbacks) replace(cb *callback) bool {
	for i := 0; i < len(*list); i++ {
		if reflect.ValueOf(cb.f).Pointer() == reflect.ValueOf((*list)[i].f).Pointer() {
			(*list)[i] = cb
			return true
		}
	}
	return false
}

func (list *callbacks) remove(f ...func(context.Context, ...interface{}) error) callbacks {
	if len(f) == 0 {
		return (*list)[:0]
	}
	for i := 0; i < len(*list); i++ {
		for _, item := range f {
			if reflect.ValueOf((*list)[i].f).Pointer() == reflect.ValueOf(item).Pointer() {
				*list = append((*list)[:i], (*list)[i+1:]...)
				i--
				break
			}
		}
	}
	return *list
}

func (list *callbacks) markRemove(f ...func(context.Context, ...interface{}) error) callbacks {
	for i := 0; i < len(*list); i++ {
		for _, item := range f {
			if reflect.ValueOf((*list)[i].f).Pointer() == reflect.ValueOf(item).Pointer() {
				(*list)[i].remove = true
				break
			}
		}
	}
	return *list
}

func (list *callbacks) markRemoveAll() callbacks {
	for i := 0; i < len(*list); i++ {
		(*list)[i].remove = true
	}
	return *list
}

func (list *callbacks) clearRemoveFlags() callbacks {
	for i := 0; i < len(*list); i++ {
		if (*list)[i].remove {
			*list = append((*list)[:i], (*list)[i+1:]...)
			i--
		}
	}
	return *list
}

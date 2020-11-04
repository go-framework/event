package inapp

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
	"time"
)

var (
	ErrTest       = errors.New("test error")
	ErrPanic      = errors.New("test panic")
	ErrTimeout    = errors.New("timeout")
	ErrUnexpected = errors.New("unexpected")
)
var (
	f1 = func(ctx context.Context, args ...interface{}) error {
		if len(args) > 0 {
			if count, ok := args[0].(*int); ok {
				*count++
			}
		}
		return nil
	}
	f2 = func(ctx context.Context, args ...interface{}) error {
		if len(args) > 0 {
			if count, ok := args[0].(*int); ok {
				*count++
			}
		}
		return nil
	}

	f3 = func(ctx context.Context, args ...interface{}) error {
		if len(args) > 0 {
			if count, ok := args[0].(*int); ok {
				*count++
			}
		}
		return nil
	}

	f4 = func(ctx context.Context, args ...interface{}) error {
		if len(args) > 0 {
			if count, ok := args[0].(*int); ok {
				*count += 1
			}
		}
		return nil
	}

	fError = func(ctx context.Context, args ...interface{}) error {
		return ErrTest
	}

	fPanic = func(ctx context.Context, args ...interface{}) error {
		panic(ErrPanic)
	}
)

var (
	cb1 = &callback{
		f:                f1,
		subscribeOptions: nil,
		remove:           true,
	}
	cb2 = &callback{
		f: f2,
		subscribeOptions: &SubscribeOptions{
			Once: true,
		},
	}
	cb3 = &callback{
		f: f3,
		subscribeOptions: &SubscribeOptions{
			Once: false,
		},
	}
	cb4 = &callback{
		f: f4,
		subscribeOptions: &SubscribeOptions{
			Once: true,
		},
		remove: true,
	}
)

func Test_callbacks_replace(t *testing.T) {
	type args struct {
		cb *callback
	}
	tests := []struct {
		name string
		list callbacks
		args args
		want bool
	}{
		{
			name: "exist",
			list: callbacks{cb1, cb2, cb3, cb4},
			args: args{
				cb: &callback{
					f: f1,
					subscribeOptions: &SubscribeOptions{
						Once: true,
					},
				},
			},
			want: true,
		},
		{
			name: "not exist",
			list: callbacks{cb1, cb2, cb3},
			args: args{
				cb: cb4,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.list.replace(tt.args.cb); got != tt.want {
				t.Errorf("replace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_callbacks_remove(t *testing.T) {
	type args struct {
		f []func(context.Context, ...interface{}) error
	}
	tests := []struct {
		name      string
		list      callbacks
		args      args
		want_list callbacks
	}{
		{
			name: "remove nil",
			list: callbacks{cb1, cb2, cb3, cb4},
			args: args{
			},
			want_list: make(callbacks, 0),
		},
		{
			name: "remove all",
			list: callbacks{cb1, cb2, cb3, cb4},
			args: args{
				f: []func(context.Context, ...interface{}) error{f1, f4, f2, f3},
			},
			want_list: make(callbacks, 0),
		},
		{
			name: "remove all",
			list: callbacks{cb1, cb2, cb3, cb4},
			args: args{
				f: []func(context.Context, ...interface{}) error{f4, f2, f3, f1},
			},
			want_list: make(callbacks, 0),
		},
		{
			name: "remove spec",
			list: callbacks{cb1, cb2, cb3, cb4},
			args: args{
				f: []func(context.Context, ...interface{}) error{f1, f4},
			},
			want_list: callbacks{cb2, cb3},
		},
		{
			name: "remove spec",
			list: callbacks{cb1, cb2, cb3, cb4},
			args: args{
				f: []func(context.Context, ...interface{}) error{f4, f2},
			},
			want_list: callbacks{cb1, cb3},
		},
		{
			name: "remove spec",
			list: callbacks{cb1, cb2, cb3, cb4},
			args: args{
				f: []func(context.Context, ...interface{}) error{f4, f3, f2},
			},
			want_list: callbacks{cb1},
		},
		{
			name: "remove first",
			list: callbacks{cb1, cb2, cb3, cb4},
			args: args{
				f: []func(context.Context, ...interface{}) error{f1},
			},
			want_list: callbacks{cb2, cb3, cb4},
		},
		{
			name: "remove last",
			list: callbacks{cb1, cb2, cb3, cb4},
			args: args{
				f: []func(context.Context, ...interface{}) error{f4},
			},
			want_list: callbacks{cb1, cb2, cb3},
		},
		{
			name: "remove empty",
			args: args{
				f: []func(context.Context, ...interface{}) error{f4},
			},
		},
		{
			name: "remove not exist",
			list: callbacks{cb1, cb2, cb3},
			args: args{
				f: []func(context.Context, ...interface{}) error{f4},
			},
			want_list: callbacks{cb1, cb2, cb3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got_list := tt.list.remove(tt.args.f...); !reflect.DeepEqual(got_list, tt.want_list) {
				t.Errorf("remove() = %v, want %v", got_list, tt.want_list)
			}
		})
	}
}

func Test_callbacks_markRemove(t *testing.T) {
	type args struct {
		f []func(context.Context, ...interface{}) error
	}
	tests := []struct {
		name   string
		list   callbacks
		args   args
		want   callbacks
		result func(callbacks) error
	}{
		{
			name: "spec",
			list: callbacks{
				&callback{
					f:                f1,
					subscribeOptions: nil,
				},
				&callback{
					f: f2,
					subscribeOptions: &SubscribeOptions{
						Once: true,
					},
				},
				&callback{
					f: f3,
					subscribeOptions: &SubscribeOptions{
						Once: false,
					},
				},
				&callback{
					f: f4,
					subscribeOptions: &SubscribeOptions{
						Once: true,
					},
				},
			},
			args: args{
				f: []func(context.Context, ...interface{}) error{f2},
			},
			result: func(list callbacks) error {
				if len(list) != 4 {
					return fmt.Errorf("want length 4, got %d", len(list))
				}
				if !list[1].remove {
					return fmt.Errorf("want remove true, got %t", list[1].remove)
				}
				return nil
			},
		},
		{
			name: "empty",
			list: callbacks{
				&callback{
					f:                f1,
					subscribeOptions: nil,
				},
				&callback{
					f: f2,
					subscribeOptions: &SubscribeOptions{
						Once: true,
					},
				},
				&callback{
					f: f3,
					subscribeOptions: &SubscribeOptions{
						Once: false,
					},
				},
				&callback{
					f: f4,
					subscribeOptions: &SubscribeOptions{
						Once: true,
					},
				},
			},
			args: args{
				f: []func(context.Context, ...interface{}) error{},
			},
			result: func(list callbacks) error {
				if len(list) != 4 {
					return fmt.Errorf("want length 4, got %d", len(list))
				}
				for idx, item := range list {
					if item.remove {
						return fmt.Errorf("want remove false, got %t @%d", item.remove, idx)
					}
				}
				return nil
			},
		},
		{
			name: "more",
			list: callbacks{
				&callback{
					f:                f1,
					subscribeOptions: nil,
				},
				&callback{
					f: f2,
					subscribeOptions: &SubscribeOptions{
						Once: true,
					},
				},
				&callback{
					f: f3,
					subscribeOptions: &SubscribeOptions{
						Once: false,
					},
				},
			},
			args: args{
				f: []func(context.Context, ...interface{}) error{f1, f2, f3, f4},
			},
			result: func(list callbacks) error {
				if len(list) != 3 {
					return fmt.Errorf("want length 3, got %d", len(list))
				}
				for idx, item := range list {
					if !item.remove {
						return fmt.Errorf("want remove true, got %t @%d", item.remove, idx)
					}
				}
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.list.markRemove(tt.args.f...)
			if err := tt.result(got); err != nil {
				t.Errorf("markRemove() error %v", err)
			}
		})
	}
}

func Test_callbacks_clearRemoveFlags(t *testing.T) {
	tests := []struct {
		name string
		list callbacks
		want callbacks
	}{
		{
			name: "nil",
			list: callbacks{},
			want: callbacks{},
		},
		{
			name: "all",
			list: callbacks{cb1, cb4},
			want: callbacks{},
		},
		{
			name: "spec",
			list: callbacks{cb1, cb2, cb3, cb4},
			want: callbacks{cb2, cb3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.list.clearRemoveFlags(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("clearRemoveFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvent_Subscribe(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
		f    func(context.Context, ...interface{}) error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "new test f1",
			args: args{
				ctx:  context.TODO(),
				name: "test",
				f:    f1,
			},
		},
		{
			name: "new test f1 with once",
			args: args{
				ctx:  NewSubscribeOptionContext(context.TODO(), WithOnceOption(true)),
				name: "test",
				f:    f1,
			},
		},
		{
			name: "new test f2",
			args: args{
				ctx:  context.TODO(),
				name: "test",
				f:    f2,
			},
		},
		{
			name: "new test f3 with once",
			args: args{
				ctx:  NewSubscribeOptionContext(context.TODO(), WithOnceOption(true)),
				name: "test",
				f:    f3,
			},
		},
		{
			name: "new test f3",
			args: args{
				ctx:  context.TODO(),
				name: "test",
				f:    f3,
			},
		},
	}

	e := &Event{}
	for _, tt := range tests {
		e.Subscribe(tt.args.ctx, tt.args.name, tt.args.f)
	}

	// test event
	value, ok := e.list.Load("test")
	if !ok {
		t.Fatalf("should have test name")
	}
	eventCase, ok := value.(*event)
	if !ok {
		t.Fatalf("should be callbacks type")
	}
	if len(eventCase.callbacks) != 3 {
		t.Fatalf("want test length 3, got %d", len(eventCase.callbacks))
	}
	for _, item := range eventCase.callbacks {
		if reflect.ValueOf(item.f).Pointer() == reflect.ValueOf(f1).Pointer() {
			if item.subscribeOptions.Once != true {
				t.Fatalf("want once true, got %t", item.subscribeOptions.Once)
			}
		} else if reflect.ValueOf(item.f).Pointer() == reflect.ValueOf(f2).Pointer() {
			if item.subscribeOptions != nil {
				t.Fatalf("want subscribeOptions nil, got %v", item.subscribeOptions)
			}
		} else if reflect.ValueOf(item.f).Pointer() == reflect.ValueOf(f3).Pointer() {
			if item.subscribeOptions != nil {
				t.Fatalf("want subscribeOptions nil, got %v", item.subscribeOptions)
			}
		}
	}
}

func TestEvent_Publish(t *testing.T) {
	var (
		count int
		errCh      = make(chan error)
		i8    int8 = 1
	)

	var (
		fData = func(ctx context.Context, args ...interface{}) error {
			data, ok := GetDataFromContext(ctx)
			if !ok {
				return ErrUnexpected
			}
			switch v := data.(type) {
			case int8:
				if v != 1 {
					return ErrUnexpected
				}
			case *int8:
				if *v != 1 {
					return ErrUnexpected
				}
				*v += 1
			}
			return nil
		}
	)

	type args struct {
		ctx  context.Context
		name string
		args []interface{}
	}
	tests := []struct {
		name    string
		args    args
		init    func(*Event, *args)
		result  func(*Event) error
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				ctx:  context.TODO(),
				name: "test",
				args: []interface{}{&count},
			},
			init: func(e *Event, args *args) {
				e.Subscribe(args.ctx, args.name, f1)
				e.Subscribe(args.ctx, args.name, f2)
				e.Subscribe(args.ctx, args.name, f3)
				e.Subscribe(args.ctx, args.name, f4)
			},
			result: func(e *Event) error {
				return nil
			},
			wantErr: false,
		},
		{
			name: "not exist",
			args: args{
				ctx:  context.TODO(),
				name: "test",
				args: []interface{}{&count},
			},
			wantErr: true,
		},
		{
			name: "return error with error option",
			args: args{
				ctx:  NewPublishOptionContext(context.TODO(), WithErrorOption(errCh)),
				name: "test",
				args: []interface{}{&count},
			},
			init: func(e *Event, args *args) {
				e.Subscribe(args.ctx, args.name, f1)
				e.Subscribe(args.ctx, args.name, f2)
				e.Subscribe(args.ctx, args.name, f3)
				e.Subscribe(args.ctx, args.name, f4)
				e.Subscribe(args.ctx, args.name, fError)
			},
			result: func(e *Event) error {
				timer := time.NewTimer(time.Second * 3)
				defer timer.Stop()
				select {
				case err := <-errCh:
					switch e := err.(type) {
					case Errors:
						for _, item := range e {
							if item == ErrTest {
								return nil
							}
						}
					case error:
						if e == ErrTest {
							return nil
						}
					}
					return err
				case <-timer.C:
					return ErrTimeout
				}
			},
			wantErr: false,
		},
		{
			name: "return nil with strict and error option",
			args: args{
				ctx:  NewPublishOptionContext(context.TODO(), WithErrorOption(errCh)),
				name: "test",
				args: []interface{}{&count},
			},
			init: func(e *Event, args *args) {
				e.Subscribe(args.ctx, args.name, f1)
				e.Subscribe(args.ctx, args.name, f2)
				e.Subscribe(args.ctx, args.name, f3)
				e.Subscribe(args.ctx, args.name, f4)
			},
			result: func(e *Event) error {
				timer := time.NewTimer(time.Second * 3)
				defer timer.Stop()
				select {
				case err := <-errCh:
					return err
				case <-timer.C:
					return ErrTimeout
				}
			},
			wantErr: false,
		},
		{
			name: "return error with strict and error option",
			args: args{
				ctx:  NewPublishOptionContext(context.TODO(), WithErrorOption(errCh), WithStrictModeOption(true)),
				name: "test",
				args: []interface{}{&count},
			},
			init: func(e *Event, args *args) {
				e.Subscribe(args.ctx, args.name, f1)
				e.Subscribe(args.ctx, args.name, f2)
				e.Subscribe(args.ctx, args.name, f3)
				e.Subscribe(args.ctx, args.name, f4)
				e.Subscribe(args.ctx, args.name, fError)
			},
			result: func(e *Event) error {
				timer := time.NewTimer(time.Second * 3)
				defer timer.Stop()
				select {
				case err := <-errCh:
					if err == ErrTest {
						return nil
					}
					return err
				case <-timer.C:
					return ErrTimeout
				}
			},
			wantErr: false,
		},
		{
			name: "panic with error option",
			args: args{
				ctx:  NewPublishOptionContext(context.TODO(), WithErrorOption(errCh)),
				name: "test",
				args: []interface{}{&count},
			},
			init: func(e *Event, args *args) {
				e.Subscribe(args.ctx, args.name, f1)
				e.Subscribe(args.ctx, args.name, f2)
				e.Subscribe(args.ctx, args.name, f3)
				e.Subscribe(args.ctx, args.name, f4)
				e.Subscribe(args.ctx, args.name, fPanic)
			},
			result: func(e *Event) error {
				timer := time.NewTimer(time.Second * 3)
				defer timer.Stop()
				select {
				case err := <-errCh:
					if err == ErrPanic {
						return nil
					}
					return err
				case <-timer.C:
					return ErrTimeout
				}
			},
			wantErr: false,
		},
		{
			name: "with data context",
			args: args{
				ctx:  NewPublishOptionContext(NewDataContext(context.TODO(), i8), WithErrorOption(errCh)),
				name: "test",
				args: []interface{}{&count},
			},
			init: func(e *Event, args *args) {
				e.Subscribe(args.ctx, args.name, fData)
			},
			result: func(e *Event) error {
				timer := time.NewTimer(time.Second * 3)
				defer timer.Stop()
				select {
				case err := <-errCh:
					if err == nil {
						return nil
					}
					return err
				case <-timer.C:
					return ErrTimeout
				}
			},
			wantErr: false,
		},
		{
			name: "with point data context",
			args: args{
				ctx:  NewPublishOptionContext(NewDataContext(context.TODO(), &i8), WithErrorOption(errCh)),
				name: "test",
				args: []interface{}{&count},
			},
			init: func(e *Event, args *args) {
				e.Subscribe(args.ctx, args.name, fData)
			},
			result: func(e *Event) error {
				timer := time.NewTimer(time.Second * 3)
				defer timer.Stop()
				select {
				case err := <-errCh:
					if i8 != 2 {
						return ErrUnexpected
					}
					if err == nil {
						return nil
					}
					return err
				case <-timer.C:
					return ErrTimeout
				}
			},
			wantErr: false,
		},
		{
			name: "with point data context once",
			args: args{
				ctx:  NewPublishOptionContext(NewDataContext(context.TODO(), &i8), WithErrorOption(errCh)),
				name: "test",
				args: []interface{}{&count},
			},
			init: func(e *Event, args *args) {
				e.Subscribe(NewSubscribeOptionContext(args.ctx, WithOnceOption(true)), args.name, fData)
			},
			result: func(e *Event) error {
				timer := time.NewTimer(time.Second * 3)
				defer timer.Stop()
				select {
				case err := <-errCh:
					if i8 != 2 {
						return ErrUnexpected
					}
					if _, ok := e.list.Load("test"); ok {
						return ErrUnexpected
					} else {
						return nil
					}
					if err == nil {
						return nil
					}
					return err
				case <-timer.C:
					return ErrTimeout
				}
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Event{}
			if tt.init != nil {
				tt.init(e, &tt.args)
			}
			if err := e.Publish(tt.args.ctx, tt.args.name, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.result != nil {
				if err := tt.result(e); err != nil {
					t.Errorf("Publish() error %v", err)
				}
			}
		})
	}
}

func TestEvent_Unsubscribe(t *testing.T) {
	type args struct {
		name string
		f    []func(context.Context, ...interface{}) error
	}
	tests := []struct {
		name   string
		init   func(e *Event, args *args)
		args   args
		result func(*Event, *args) error
	}{
		{
			name: "nil with doneLock busy",
			init: func(e *Event, args *args) {
				e.Subscribe(context.TODO(), args.name, f1)
				e.Subscribe(context.TODO(), args.name, f2)
				e.Subscribe(context.TODO(), args.name, f3)
				value, _ := e.list.Load(args.name)
				var event = value.(*event)
				<-event.doneLock
			},
			args: args{
				name: "test",
				f:    nil,
			},
			result: func(e *Event, args *args) error {
				value, ok := e.list.Load(args.name)
				if !ok {
					return fmt.Errorf("should be have the %s event", args.name)
				}
				var event = value.(*event)
				for _, item := range event.callbacks {
					if !item.remove {
						return fmt.Errorf("want remove true, got false")
					}
				}
				return nil
			},
		},
		{
			name: "all with doneLock busy",
			init: func(e *Event, args *args) {
				e.Subscribe(context.TODO(), args.name, f1)
				e.Subscribe(context.TODO(), args.name, f2)
				e.Subscribe(context.TODO(), args.name, f3)
				value, _ := e.list.Load(args.name)
				var event = value.(*event)
				<-event.doneLock
			},
			args: args{
				name: "test",
				f:    []func(context.Context, ...interface{}) error{f1, f2, f3, f4},
			},
			result: func(e *Event, args *args) error {
				value, ok := e.list.Load(args.name)
				if !ok {
					return fmt.Errorf("should be have the %s event", args.name)
				}
				var event = value.(*event)
				for _, item := range event.callbacks {
					if !item.remove {
						return fmt.Errorf("want remove true, got false")
					}
				}
				return nil
			},
		},
		{
			name: "spec with doneLock busy",
			init: func(e *Event, args *args) {
				e.Subscribe(context.TODO(), args.name, f1)
				e.Subscribe(context.TODO(), args.name, f2)
				e.Subscribe(context.TODO(), args.name, f3)
				value, _ := e.list.Load(args.name)
				var event = value.(*event)
				<-event.doneLock
			},
			args: args{
				name: "test",
				f:    []func(context.Context, ...interface{}) error{f1, f3},
			},
			result: func(e *Event, args *args) error {
				value, ok := e.list.Load(args.name)
				if !ok {
					return fmt.Errorf("should be have the %s event", args.name)
				}
				var event = value.(*event)
				for _, item := range event.callbacks {
					if !item.remove &&
						(reflect.ValueOf(item.f).Pointer() == reflect.ValueOf(f1).Pointer() ||
							reflect.ValueOf(item.f).Pointer() == reflect.ValueOf(f3).Pointer()) {
						return fmt.Errorf("want remove true, got false")
					}
				}
				return nil
			},
		},
		{
			name: "nil",
			init: func(e *Event, args *args) {
				e.Subscribe(context.TODO(), args.name, f1)
				e.Subscribe(context.TODO(), args.name, f2)
				e.Subscribe(context.TODO(), args.name, f3)
			},
			args: args{
				name: "test",
				f:    nil,
			},
			result: func(e *Event, args *args) error {
				_, ok := e.list.Load(args.name)
				if ok {
					return fmt.Errorf("should be remove the %s event", args.name)
				}
				return nil
			},
		},
		{
			name: "multiple",
			init: func(e *Event, args *args) {
				e.Subscribe(context.TODO(), args.name, f1)
				e.Subscribe(context.TODO(), args.name, f2)
				e.Subscribe(context.TODO(), args.name+"1", f3)
			},
			args: args{
				name: "test",
				f:    nil,
			},
			result: func(e *Event, args *args) error {
				_, ok := e.list.Load(args.name)
				if ok {
					return fmt.Errorf("should be remove the %s event", args.name)
				}
				_, ok = e.list.Load(args.name + "1")
				if !ok {
					return fmt.Errorf("should be have the %s event", args.name+"1")
				}
				return nil
			},
		},
		{
			name: "first",
			init: func(e *Event, args *args) {
				e.Subscribe(context.TODO(), args.name, f1)
				e.Subscribe(context.TODO(), args.name, f2)
				e.Subscribe(context.TODO(), args.name, f3)
				e.Subscribe(context.TODO(), args.name, f4)
			},
			args: args{
				name: "test",
				f:    []func(context.Context, ...interface{}) error{f1},
			},
			result: func(e *Event, args *args) error {
				value, ok := e.list.Load(args.name)
				if !ok {
					return fmt.Errorf("should be have the %s event", args.name)
				}
				var event = value.(*event)
				if len(event.callbacks) != 3 {
					return fmt.Errorf("want length 3, got %d", len(event.callbacks))
				}
				for _, item := range event.callbacks {
					if reflect.ValueOf(item.f).Pointer() == reflect.ValueOf(f1).Pointer() {
						return fmt.Errorf("should be not have the %v func", reflect.TypeOf(f1))
					}
				}
				return nil
			},
		},
		{
			name: "all",
			init: func(e *Event, args *args) {
				e.Subscribe(context.TODO(), args.name, f1)
				e.Subscribe(context.TODO(), args.name, f2)
				e.Subscribe(context.TODO(), args.name, f3)
				e.Subscribe(context.TODO(), args.name, f4)
			},
			args: args{
				name: "test",
				f:    []func(context.Context, ...interface{}) error{f1, f2, f3, f4},
			},
			result: func(e *Event, args *args) error {
				_, ok := e.list.Load(args.name)
				if ok {
					return fmt.Errorf("should be remove the %s event", args.name)
				}
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Event{}
			if tt.init != nil {
				tt.init(e, &tt.args)
			}
			e.Unsubscribe(tt.args.name, tt.args.f...)
			if tt.result != nil {
				if err := tt.result(e, &tt.args); err != nil {
					t.Errorf("Unsubscribe() error = %v", err)
				}
			}
		})
	}
}

func TestEventMutex(t *testing.T) {
	var (
		name = "test"

		normal1 = func(ctx context.Context, args ...interface{}) error {
			if len(args) > 0 {
				switch v := args[0].(type) {
				case int:
					log.Printf("---------------- normal1 got count %d ----------------", v)
				case int64:
					log.Printf("---------------- normal1 got timestamp %d ----------------", v)
				case string:
					log.Printf("---------------- normal1 got string %s ----------------", v)
				}
			}
			return nil
		}

		normal2 = func(ctx context.Context, args ...interface{}) error {
			if len(args) > 0 {
				switch v := args[0].(type) {
				case int:
					log.Printf("---------------- normal2 got count %d ----------------", v)
				case int64:
					log.Printf("---------------- normal2 got timestamp %d ----------------", v)
				case string:
					log.Printf("---------------- normal2 got string %s ----------------", v)
				}
			}
			return nil
		}

		normal3 = func(ctx context.Context, args ...interface{}) error {
			if len(args) > 0 {
				switch v := args[0].(type) {
				case int:
					log.Printf("---------------- normal3 got count %d ----------------", v)
				case int64:
					log.Printf("---------------- normal3 got timestamp %d ----------------", v)
				case string:
					log.Printf("---------------- normal3 got string %s ----------------", v)
				}
			}
			return nil
		}
	)

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	e := NewEvent()
	ctx, cancel := context.WithTimeout(context.TODO(), time.Minute*10)
	rand.Seed(time.Now().UnixNano())

	log.Printf("start ...")
	defer func() {
		log.Printf("exit ...")
	}()

	go func() {
		var n int
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			n = rand.Int()
			if n%3 == 0 {
				e.Subscribe(NewSubscribeOptionContext(context.TODO(), WithOnceOption(true)), name, normal1)
				log.Printf("Subscribe normal1 with once")
			} else if n%5 == 0 {
				e.Subscribe(NewSubscribeOptionContext(context.TODO(), WithOnceOption(false)), name, normal2)
				log.Printf("Subscribe normal2")
			} else if n%7 == 0 {
				e.Subscribe(NewSubscribeOptionContext(context.TODO(), WithOnceOption(true)), name, normal3)
				log.Printf("Subscribe normal3 with once")
				e.Subscribe(NewSubscribeOptionContext(context.TODO(), WithOnceOption(false)), name, normal1)
				log.Printf("Subscribe normal1")
				e.Subscribe(NewSubscribeOptionContext(context.TODO(), WithOnceOption(false)), name, normal2)
				log.Printf("Subscribe normal2")
			}
			time.Sleep(time.Duration(100+rand.Intn(500)) * time.Millisecond)
		}
	}()

	go func() {
		var n int
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			n = rand.Int()
			timestamp := time.Now().UnixNano()
			if n%3 == 0 {
				err := e.Publish(context.TODO(), name, timestamp)
				if err != nil {
					log.Printf("Publish with timestamp: %d error: %v", timestamp, err)
				} else {
					log.Printf("Publish with timestamp: %d", timestamp)
				}
			} else if n%5 == 0 {
				err := e.Publish(context.TODO(), name, strconv.FormatInt(timestamp, 10))
				if err != nil {
					log.Printf("Publish with timestamp: %d error: %v", timestamp, err)
				} else {
					log.Printf("Publish with string timestamp: %d", timestamp)
				}
			}
			time.Sleep(time.Duration(100+rand.Intn(500)) * time.Millisecond)
		}
	}()

	go func() {
		var n int

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			n = rand.Int()
			if n%3 == 0 {
				e.Unsubscribe(name, normal3)
				log.Printf("Unsubscribe normal3")
			} else if n%5 == 0 {
				e.Unsubscribe(name, normal1)
				log.Printf("Unsubscribe normal1")
			} else if n%7 == 0 {
				e.Unsubscribe(name)
				log.Printf("Unsubscribe all")
			}
			time.Sleep(time.Duration(100+rand.Intn(500)) * time.Millisecond)
		}
	}()

	select {
	case <-ctx.Done():
		cancel()
		return
	}
}

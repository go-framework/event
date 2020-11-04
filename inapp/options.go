package inapp

// Subscribe option func.
type SubscribeOption func(options *SubscribeOptions)

// Subscribe options.
type SubscribeOptions struct {
	Once bool // Listen for a Event, but only once. The listener will be removed once it triggers for the first time.
}

// Get default SubscribeOptions value.
func GetDefaultSubscribeOptions() *SubscribeOptions {
	opts := &SubscribeOptions{}
	return opts
}

// Once subscribe option.
func WithOnceOption(once bool) SubscribeOption {
	return func(options *SubscribeOptions) {
		options.Once = once
	}
}

// Publish option func.
type PublishOption func(options *PublishOptions)

// Publish options.
type PublishOptions struct {
	Strict bool       // Strict mode, when done callback error strict is true will be stop and return.
	Err    chan error // Err is finished signal, value is publish callback return.
}

// Get default PublishOptions value.
func GetDefaultPublishOptions() *PublishOptions {
	opts := &PublishOptions{}
	return opts
}

// WithStrictModeOption will be stop publish when callback func got error.
func WithStrictModeOption(strict bool) PublishOption {
	return func(options *PublishOptions) {
		options.Strict = strict
	}
}

// WithErrorOption will got the callback finished signal.
func WithErrorOption(ch chan error) PublishOption {
	return func(options *PublishOptions) {
		options.Err = ch
	}
}

package resolver

import (
	"context"
	"time"
)

type Options struct {
	interval time.Duration
	ctx      context.Context
	instant  bool
}

type Option func(options *Options)

func OptionInterval(v time.Duration) Option {
	return func(options *Options) {
		options.interval = v
	}
}

func OptionContext(v context.Context) Option {
	return func(options *Options) {
		options.ctx = v
	}
}

func OptionInstant(v bool) Option {
	return func(options *Options) {
		options.instant = v
	}
}

func NewOptions() Options {
	return Options{
		ctx: context.Background(),
	}
}

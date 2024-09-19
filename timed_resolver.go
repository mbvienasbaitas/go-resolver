package resolver

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type TimedResolver[T any] struct {
	lock       *sync.RWMutex
	opts       Options
	resolver   Resolver[T]
	isResolved bool
	resolved   T
}

func (receiver *TimedResolver[T]) Resolve(ctx context.Context) (T, error) {
	receiver.lock.RLock()

	if receiver.isResolved == true {
		receiver.lock.RUnlock()

		return receiver.resolved, nil
	}

	receiver.lock.RUnlock()

	return receiver.resolveAndBind(ctx)
}

func (receiver *TimedResolver[T]) resolveAndBind(ctx context.Context) (T, error) {
	receiver.lock.Lock()

	defer receiver.lock.Unlock()

	resolved, err := receiver.resolver.Resolve(ctx)

	if err != nil {
		return receiver.resolved, err
	}

	receiver.resolved = resolved
	receiver.isResolved = true

	return receiver.resolved, nil
}

func NewTimedResolver[T any](resolver Resolver[T], opts ...Option) (Resolver[T], error) {
	options := NewOptions()

	for _, o := range opts {
		o(&options)
	}

	if resolver == nil {
		return nil, fmt.Errorf("resolver must be set")
	}

	if options.ctx == nil {
		return nil, fmt.Errorf("context must be set")
	}

	svc := &TimedResolver[T]{
		lock:       &sync.RWMutex{},
		opts:       options,
		isResolved: false,
		resolver:   resolver,
	}

	if options.instant {
		_, _ = svc.resolveAndBind(options.ctx)
	}

	// launch auto refresher
	go func(d time.Duration) {
		ticker := time.NewTicker(d)

		for {
			select {
			case <-ticker.C:
				_, _ = svc.resolveAndBind(options.ctx)
			}
		}
	}(options.interval)

	return svc, nil
}

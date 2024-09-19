package resolver

import (
	"context"
)

type FuncResolver[T any] func(ctx context.Context) (T, error)

func (receiver FuncResolver[T]) Resolve(ctx context.Context) (T, error) {
	return receiver(ctx)
}

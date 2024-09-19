package resolver

import (
	"context"
)

type Resolver[T any] interface {
	Resolve(ctx context.Context) (T, error)
}

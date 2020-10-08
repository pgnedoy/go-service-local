package types

import "context"

type TaskFunc func(ctx context.Context, input string) (string, error)

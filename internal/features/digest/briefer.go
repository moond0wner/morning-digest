package digest

import "context"

type Briefer interface {
	Name() string
	Brief(ctx context.Context) (string, error)
}

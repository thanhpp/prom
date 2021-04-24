package booting

import (
	"context"
)

// Daemon start func & stop func
type Daemon func(ctx context.Context) (start func() error, stop func())

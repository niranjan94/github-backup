package internal

import (
	"context"
	"time"
)

// DelayWithContext returns nil after the specified duration or error if interrupted.
func DelayWithContext(ctx context.Context, d time.Duration) {
	t := time.NewTimer(d)
	select {
	case <-ctx.Done():
		t.Stop()
	case <-t.C:
	}
}

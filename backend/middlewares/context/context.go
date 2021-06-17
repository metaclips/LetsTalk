package context

import "context"

// Create a sub context so that we can deny Admin operation access to database
// on a failed operation.
// Note: cancelling child context should not affect parent context.
var ctx, cancel = context.WithCancel(context.Background())

func CancelContext() error {
	if err := ctx.Err(); err == context.Canceled {
		return err
	}

	cancel()
	return nil
}

// Creates a child context from main context.
func CreateChildCancelContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(ctx)
}

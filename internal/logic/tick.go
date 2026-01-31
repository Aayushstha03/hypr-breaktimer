package logic

import (
	"context"
)

func Tick(ctx context.Context) error {
	return show(ctx, showDueOnly)
}

package schema

import (
	"context"

	"github.com/chromedp/chromedp"
)

func WithBrowserContext(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := chromedp.NewContext(ctx)
	if err := chromedp.Run(ctx); err != nil {
		panic(err)
	}
	return ctx, cancel
}

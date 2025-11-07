package main

import (
	"context"

	"github.com/chistyakoviv/logbot/internal/app"
)

func main() {
	ctx := context.Background()
	a := app.NewApp(ctx)
	a.Run(ctx)
}

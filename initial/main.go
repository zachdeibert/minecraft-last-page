package main

import (
	"context"

	"github.com/thecodeteam/goodbye"

	"./logic"
)

func main() {
	ctx := context.Background()
	defer func() {
		if e := recover(); e != nil {
			panic(e)
		} else {
			goodbye.Exit(ctx, -1)
		}
	}()
	goodbye.Notify(ctx)
	logic.Run()
}

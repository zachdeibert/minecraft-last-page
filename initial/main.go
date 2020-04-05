package main

import (
	"context"
	"fmt"

	"github.com/thecodeteam/goodbye"

	"./gpu"
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
	sieve, e := gpu.CreateSieve()
	if e != nil {
		panic(e)
	}
	defer sieve.Close()
	if e = sieve.Configure(1, 3, 2, 10); e != nil {
		panic(e)
	}
	res, e := sieve.Run("n")
	if e != nil {
		panic(e)
	}
	for _, s := range res {
		fmt.Println(s)
	}
}

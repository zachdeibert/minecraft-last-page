package gpu

import (
	ct "context"
	"os"

	"github.com/thecodeteam/goodbye"
)

var (
	cleanupFuncs []func()
)

func init() {
	goodbye.Register(func(ctx ct.Context, sig os.Signal) {
		for i := range cleanupFuncs {
			cleanupFuncs[len(cleanupFuncs)-1-i]()
		}
	})
}

package pool

import (
	"runtime"

	"github.com/panjf2000/ants/v2"
	"github.com/someview/trade-base/log"
)

var routinePool, _ = ants.NewPool(2*1e4, ants.WithNonblocking(true))

func Go(f func()) {
	err := routinePool.Submit(f)
	if err == nil {
		return
	}
	log.Debug("routinePool Go error", log.String("error", err.Error()), log.Int("routineNum", runtime.NumGoroutine()))
	go func() {
		f()
	}()
}

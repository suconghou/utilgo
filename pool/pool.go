package pool

import (
	"sync/atomic"
	"time"
)

// GoPool for go thread
type GoPool struct {
	runing int32
	max    int32
	reduce uint8
	fun    chan func()
}

// New create new pool
func New(max int32, reduce uint8) *GoPool {
	return &GoPool{
		runing: 0,
		max:    max,
		fun:    make(chan func(), max),
		reduce: reduce,
	}
}

// Put create new thread until reach max
func (g *GoPool) Put(f func()) bool {
	g.fun <- f
	if atomic.LoadInt32(&g.runing) >= g.max {
		return false
	}
	go func() {
		defer func() {
			atomic.AddInt32(&g.runing, -1)
		}()
		if g.reduce > 0 {
			for {
				select {
				case f := <-g.fun:
					f()
				case <-time.After(time.Second * time.Duration(g.reduce)):
					return
				}
			}
		} else {
			for {
				f := <-g.fun
				f()
			}
		}
	}()
	atomic.AddInt32(&g.runing, 1)
	return true
}

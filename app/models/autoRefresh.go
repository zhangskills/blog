package models

import (
	"time"
)

type AutoRefresh struct {
	needRefresh bool
}

func NewAutoRefresh(d time.Duration, refreshFn func()) *AutoRefresh {
	autoRefresh := &AutoRefresh{}
	go func() {
		for {
			if autoRefresh.needRefresh {
				refreshFn()
				autoRefresh.needRefresh = false
			}
			time.Sleep(d)
		}
	}()
	return autoRefresh
}

func (a *AutoRefresh) SetStatus(needRefresh bool) {
	a.needRefresh = needRefresh
}

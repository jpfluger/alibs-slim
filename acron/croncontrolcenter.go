package acron

import "sync"

type ICronControlCenter interface {
	GetJRun() IJRun
}

type CronControlCenter struct {
	jrun IJRun

	mu sync.RWMutex
}

func (c *CronControlCenter) GetJRun() IJRun {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.jrun
}

func (c *CronControlCenter) SetJRun(jrun IJRun) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.jrun = jrun
}

package tools

import "sync"

type PriorityMutex struct {
	normalChan       chan struct{}
	highPriorityChan chan struct{}
	mutex            sync.Mutex
}

func NewPriorityMutex() *PriorityMutex {
	return &PriorityMutex{
		normalChan:       make(chan struct{}, 1),
		highPriorityChan: make(chan struct{}, 1),
	}
}

func (pm *PriorityMutex) Lock(highPriority bool) {
	if highPriority {
		pm.highPriorityChan <- struct{}{}
	} else {
		pm.normalChan <- struct{}{}
	}
	pm.mutex.Lock()
}

func (pm *PriorityMutex) Unlock(highPriority bool) {
	pm.mutex.Unlock()
	if highPriority {
		<-pm.highPriorityChan
	} else {
		<-pm.normalChan
	}
}

package tests

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

func (pm *PriorityMutex) TryLock() bool {
	select {
	case <-pm.highPriorityChan:
		if pm.mutex.TryLock() {
			return true
		}
		pm.highPriorityChan <- struct{}{}
	case <-pm.normalChan:
		if pm.mutex.TryLock() {
			return true
		}
		pm.normalChan <- struct{}{}
	default:
	}
	return false
}

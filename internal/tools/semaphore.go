package tools

type Semaphore chan int

func NewSemaphore() *Semaphore {
	sem := make(Semaphore, 1)
	return &sem
}

func (s Semaphore) Acquire() {
	s <- 1
}

func (s Semaphore) Release() {
	<-s
}

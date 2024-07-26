package tests

import (
	"fmt"
	"sync"
	"time"
)

const numFilosofos = 5
const refeicoesDesejadas = 1

type Filosofo struct {
	id                          int
	garfoEsquerda, garfoDireita *sync.Mutex
	refeicoes                   int
}

func (f *Filosofo) pensar() {
	fmt.Printf("Filosofo %d está pensando...\n", f.id)
	time.Sleep(time.Millisecond * time.Duration(50))
}

func (f *Filosofo) comer() {
	fmt.Printf("Filosofo %d está comendo...\n", f.id)
	time.Sleep(time.Second * 10)
	f.refeicoes++
	fmt.Printf("Filosofo %d terminou de comer.\n", f.id)
}

func (f *Filosofo) tentarComer() {
	for f.refeicoes < refeicoesDesejadas {
		f.pensar()

		fmt.Printf("Filosofo %d está tentando pegar o garfo esquerdo...\n", f.id)
		f.garfoEsquerda.Lock()
		fmt.Printf("Filosofo %d pegou o garfo esquerdo (posição %d).\n", f.id, f.id)

		fmt.Printf("Filosofo %d está tentando pegar o garfo direito...\n", f.id)
		f.garfoDireita.Lock()
		fmt.Printf("Filosofo %d pegou o garfo direito (posição %d).\n", f.id, (f.id)%numFilosofos+1)

		f.comer()

		fmt.Printf("Filosofo %d liberou garfo esquerdo...\n", f.id)
		f.garfoEsquerda.Unlock()

		fmt.Printf("Filosofo %d liberou garfo direito...\n", f.id)
		f.garfoDireita.Unlock()
	}
}

func main() {
	var wg sync.WaitGroup
	garfos := make([]*sync.Mutex, numFilosofos)
	filosofos := make([]*Filosofo, numFilosofos)

	for i := range garfos {
		garfos[i] = &sync.Mutex{}
	}

	for i := 0; i < numFilosofos; i++ {
		filosofos[i] = &Filosofo{
			id:            i + 1,
			garfoEsquerda: garfos[i],
			garfoDireita:  garfos[(i+1)%numFilosofos],
		}
	}

	for i := 0; i < numFilosofos; i++ {
		wg.Add(1)
		go func(filosofo *Filosofo) {
			defer wg.Done()
			filosofo.tentarComer()
		}(filosofos[i])
	}

	wg.Wait()
}

package entity

import (
	"fmt"
	"sync"

	"math/rand"

	"github.com/brnocorreia/concurrency/internal/tools"
)

// Implementação da matriz para MUTEX
type MatrixMutex [][]*BlockMutex

func NewMatrixMutex(width, height int) MatrixMutex {
	matrix := make(MatrixMutex, height)

	id := 1
	for i := range matrix {
		matrix[i] = make([]*BlockMutex, width)
		for j := range matrix[i] {
			matrix[i][j] = NewBlockMutex(id, &sync.Mutex{})
			id++
		}
	}
	return matrix
}

func (m MatrixMutex) GetRandomBlock() *BlockMutex {
	rows := len(m)
	cols := len(m[0])
	row := rand.Intn(rows)
	col := rand.Intn(cols)
	return m[row][col]
}

func PrintBlocksMutex(m MatrixMutex) {
	for _, row := range m {
		for _, block := range row {
			fmt.Printf("%3d ", block.Health)
		}
		fmt.Println()
	}
}

// -------------------------------------------------------------------------------------

// Implementação da matriz para SEMAPHORE
type MatrixSemaphore [][]*BlockSemaphore

func NewMatrixSemaphore(width, height int) MatrixSemaphore {
	matrix := make(MatrixSemaphore, height)

	id := 1
	for i := range matrix {
		matrix[i] = make([]*BlockSemaphore, width)
		for j := range matrix[i] {
			matrix[i][j] = NewBlockSemaphore(id, tools.NewSemaphore())
			id++
		}
	}
	return matrix
}

func (m MatrixSemaphore) GetRandomBlock() *BlockSemaphore {
	rows := len(m)
	cols := len(m[0])
	row := rand.Intn(rows)
	col := rand.Intn(cols)
	return m[row][col]
}

func PrintBlocksSemaphore(m MatrixSemaphore) {
	for _, row := range m {
		for _, block := range row {
			fmt.Printf("%3d ", block.Health)
		}
		fmt.Println()
	}
}

// -------------------------------------------------------------------------------------

// Implementação da matriz para TROCA DE MENSAGENS
type MatrixMessage [][]*BlockMessage

func NewMatrixMessage(width, height int) MatrixMessage {
	matrix := make(MatrixMessage, height)

	id := 1
	for i := range matrix {
		matrix[i] = make([]*BlockMessage, width)
		for j := range matrix[i] {
			matrix[i][j] = NewBlockMessage(id)
			id++
		}
	}
	return matrix
}

func (m MatrixMessage) GetRandomBlock() *BlockMessage {
	rows := len(m)
	cols := len(m[0])
	row := rand.Intn(rows)
	col := rand.Intn(cols)
	return m[row][col]
}

func PrintBlocksMessage(m MatrixMessage) {
	for _, row := range m {
		for _, block := range row {
			fmt.Printf("%3d ", block.Health)
		}
		fmt.Println()
	}
}

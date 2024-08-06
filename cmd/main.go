package main

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/brnocorreia/concurrency/config/logger"
	"github.com/brnocorreia/concurrency/entity"
	"go.uber.org/zap"
)

// Semaphore -> https://medium.com/@deckarep/gos-extended-concurrency-semaphores-part-1-5eeabfa351ce

func main() {

	logger.Info("Iniciando o jogo...")

	const numAttacks = 256 // Número de ataques
	const matrixSize = 8   // Tamanho da matriz (30x30)

	logger.Info("Configuração do jogo:", zap.Int("numAttacks", numAttacks), zap.Int("matrixSize", matrixSize))

	// Cria a matriz de blocos
	matrix := NewMatrix(matrixSize, matrixSize)

	logger.Info("Criando os jogadores...")
	// Cria jogadores
	player1 := entity.NewPlayer(1, 30)
	player2 := entity.NewPlayer(2, 30)

	var wg sync.WaitGroup

	results := make(chan string, 2)

	// Função para simular um ataque
	attack := func(player *entity.Player, sequence [][2]int) {
		defer wg.Done()
		for _, coord := range sequence {
			x, y := coord[0], coord[1]
			block := matrix[x][y]
			block.Hit(player)
		}

		result := fmt.Sprintf("O player %d ganhou %d pontos\n", player.Id, player.GetPoints())
		results <- result
	}

	// Inicia as goroutines para os dois jogadores
	wg.Add(2)
	go attack(player1, generateAttackSequence(matrixSize, numAttacks))
	go attack(player2, generateAttackSequence(matrixSize, numAttacks))

	// Aguarda até que ambas as goroutines terminem
	wg.Wait()

	close(results)

	// Imprime o estado final dos blocos
	fmt.Println("------------------------------------------------")
	fmt.Println()
	fmt.Println("Estado final dos blocos:")
	printBlocks(matrix)
	fmt.Println()
	fmt.Println("------------------------------------------------")

	for result := range results {
		fmt.Println(result)
	}
}

type Matrix [][]*entity.Block

func NewMatrix(width, height int) Matrix {
	matrix := make(Matrix, height)

	id := 1
	for i := range matrix {
		matrix[i] = make([]*entity.Block, width)
		for j := range matrix[i] {
			hitTime := rand.Intn(500) + 1
			matrix[i][j] = entity.NewBlock(id, hitTime, &sync.Mutex{})
			id++
		}
	}
	return matrix
}

func (m Matrix) GetRandomBlock() *entity.Block {
	rows := len(m)
	cols := len(m[0])
	row := rand.Intn(rows)
	col := rand.Intn(cols)
	return m[row][col]
}

func printBlocks(matrix Matrix) {
	for _, row := range matrix {
		for _, block := range row {
			fmt.Printf("%3d ", block.Health)
		}
		fmt.Println()
	}
}

func generateAttackSequence(size, numAttacks int) [][2]int {
	sequence := make([][2]int, numAttacks)

	// Preenche a sequência com coordenadas aleatórias
	for i := 0; i < numAttacks; i++ {
		x := rand.Intn(size)
		y := rand.Intn(size)
		sequence[i] = [2]int{x, y}
	}

	return sequence
}

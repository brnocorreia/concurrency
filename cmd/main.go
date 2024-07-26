package main

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/brnocorreia/concurrency/config/logger"
	"github.com/brnocorreia/concurrency/entity"
	"go.uber.org/zap"
)

type Matrix [][]*entity.Block

func NewMatrix(width, height int) Matrix {
	matrix := make(Matrix, height)

	id := 1
	for i := range matrix {
		matrix[i] = make([]*entity.Block, width)
		for j := range matrix[i] {
			hitTime := rand.Intn(5) + 1
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

func main() {
	// V1

	// matrix := NewMatrix(30, 30)
	// block := matrix[0][0]

	// player_1 := entity.NewPlayer(1, 10)
	// player_2 := entity.NewPlayer(2, 10)

	// fmt.Println("Block:", block.String())

	// block.Hit(player_1)
	// block.Hit(player_2)
	// fmt.Println("Block:", block.String())

	// matrix := NewMatrix(30, 30)
	// block := matrix[0][0]

	// player1 := entity.NewPlayer(1, 10)
	// player2 := entity.NewPlayer(2, 10)

	// fmt.Println("Antes dos hits:")
	// fmt.Println("Block:", block.String())

	// var wg sync.WaitGroup

	// // Função para simular o ataque de um jogador
	// attack := func(player *entity.Player) {
	// 	defer wg.Done()
	// 	block.Hit(player)
	// }

	// wg.Add(2) // Adiciona duas goroutines para esperar
	// go attack(player1)
	// go attack(player2)

	// // Aguarda até que ambas as goroutines terminem
	// wg.Wait()

	// fmt.Println("Depois dos hits:")
	// fmt.Println("Block:", block.String())

	// V2

	// matrix := NewMatrix(5, 5)
	// player1 := entity.NewPlayer(1, 10)
	// player2 := entity.NewPlayer(2, 10)

	// fmt.Println("Estado inicial dos blocos:")
	// printBlocks(matrix)

	// var wg sync.WaitGroup
	// numAttacks := 100

	// for i := 0; i < numAttacks; i++ {
	// 	wg.Add(2) // Adiciona duas goroutines para cada ataque

	// 	go func() {
	// 		defer wg.Done()
	// 		block := matrix.GetRandomBlock()
	// 		block.Hit(player1)
	// 	}()

	// 	go func() {
	// 		defer wg.Done()
	// 		block := matrix.GetRandomBlock()
	// 		block.Hit(player2)
	// 	}()
	// }

	// wg.Wait()

	// fmt.Println("Estado final dos blocos:")
	// printBlocks(matrix)
	// fmt.Printf("O player 1 ganhou %d pontos\n", player1.Points)
	// fmt.Printf("O player 2 ganhou %d pontos\n", player2.Points)

	// ----------------------------

	// V3

	logger.Info("Iniciando o jogo...")

	const numAttacks = 10 // Número de ataques
	const matrixSize = 5  // Tamanho da matriz (30x30)

	logger.Info("Configuração do jogo:", zap.Int("numAttacks", numAttacks), zap.Int("matrixSize", matrixSize))

	// Cria a matriz de blocos
	matrix := NewMatrix(matrixSize, matrixSize)

	logger.Info("Criando os jogadores...")
	// Cria jogadores
	player1 := entity.NewPlayer(1, 20)
	player2 := entity.NewPlayer(2, 20)

	var wg sync.WaitGroup

	// Função para simular um ataque
	attack := func(player *entity.Player, sequence [][2]int) {
		defer wg.Done()
		for _, coord := range sequence {
			x, y := coord[0], coord[1]
			block := matrix[x][y]
			block.Hit(player)
		}
	}

	// Inicia as goroutines para os dois jogadores
	wg.Add(2)
	go attack(player1, generateAttackSequence(matrixSize, numAttacks))
	go attack(player2, generateAttackSequence(matrixSize, numAttacks))

	// Aguarda até que ambas as goroutines terminem
	wg.Wait()

	// Imprime o estado final dos blocos
	fmt.Println("Estado final dos blocos:")
	printBlocks(matrix)

	fmt.Printf("O player 1 ganhou %d pontos\n", player1.GetPoints())
	fmt.Printf("O player 2 ganhou %d pontos\n", player2.GetPoints())
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

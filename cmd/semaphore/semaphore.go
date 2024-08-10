package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/brnocorreia/concurrency/internal/config/logger"
	"go.uber.org/zap"
)

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

type Block struct {
	Id        int
	Health    int
	Hit_time  time.Duration
	semaphore *Semaphore
}

func NewBlock(id int, semaphore *Semaphore) *Block {
	var hitTime time.Duration
	if id%2 == 0 {
		hitTime = 500 * time.Millisecond // Duração de 375ms para ids pares
	} else {
		hitTime = 125 * time.Millisecond // Duração de 125ms para ids ímpares
	}

	return &Block{
		Id:        id,
		Health:    100,
		Hit_time:  hitTime,
		semaphore: semaphore,
	}
}

func (b *Block) Hit(player *Player) bool {
	b.semaphore.Acquire()
	defer b.semaphore.Release()

	logger.Info("O player tenta acertar o bloco", zap.Int("playerId", player.Id), zap.Int("blockId", b.Id))

	if b.Health > 0 {
		time.Sleep(b.Hit_time)

		b.Health -= player.GetDamage()
		logger.Info("O player acertou o bloco", zap.Int("playerId", player.Id), zap.Int("blockId", b.Id))
		logger.Info("O bloco agora tem", zap.Int("blockId", b.Id), zap.Int("health", b.Health))
		if b.Health < 0 {
			logger.Info("O player destruiu o bloco", zap.Int("playerId", player.Id), zap.Int("blockId", b.Id))
			player.AddPoint()
			b.Health = 0
		}
		return true
	}
	return false
}

func (b *Block) IsAlive() bool {
	b.semaphore.Acquire()
	defer b.semaphore.Release()
	return b.Health > 0
}

func (b *Block) String() string {
	return fmt.Sprintf("ID=%d, Health=%d, HitTime=%v", b.Id, b.Health, b.Hit_time)
}

type Player struct {
	Id     int
	Power  int
	Points int
}

func NewPlayer(id int, power int) *Player {
	return &Player{
		Id:     id,
		Power:  power,
		Points: 0,
	}
}

func (p *Player) GetDamage() int {
	return p.Power
}

func (p *Player) GetPoints() int {
	return p.Points
}

func (p *Player) AddPoint() {
	p.Points++
}

func main() {

	logger.Info("Iniciando o jogo...")

	const numAttacks = 256 // Número de ataques
	const matrixSize = 8   // Tamanho da matriz (30x30)

	logger.Info("Configuração do jogo:", zap.Int("numAttacks", numAttacks), zap.Int("matrixSize", matrixSize))

	// Cria a matriz de blocos
	matrix := NewMatrix(matrixSize, matrixSize)

	logger.Info("Criando os jogadores...")
	// Cria jogadores
	player1 := NewPlayer(1, 30)
	player2 := NewPlayer(2, 30)

	var wg sync.WaitGroup

	results := make(chan string, 2)

	// Função para simular um ataque
	attack := func(player *Player, sequence [][2]int) {
		defer wg.Done()
		for _, coord := range sequence {
			x, y := coord[0], coord[1]
			block := matrix[x][y]
			block.Hit(player)
		}

		result := fmt.Sprintf("O player %d ganhou %d pontos\n", player.Id, player.GetPoints())
		results <- result
	}

	logger.Info("Carregando a sequência de ataques...")

	sequence_1, err := loadSequenceFromFile("sequence_1.json")
	if err != nil {
		logger.Info("Erro ao carregar a sequência 1")
		return
	}

	sequence_2, err := loadSequenceFromFile("sequence_2.json")
	if err != nil {
		logger.Info("Erro ao carregar a sequência 2")
		return
	}

	// Inicia as goroutines para os dois jogadores
	logger.Info("Iniciando as goroutines...")
	wg.Add(2)

	init := time.Now()
	go attack(player1, sequence_1)
	go attack(player2, sequence_2)

	// Aguarda até que ambas as goroutines terminem
	wg.Wait()
	duration := (time.Since(init))

	logger.Info("Tempo de execução:", zap.Duration("duration", duration))

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

type Matrix [][]*Block

func NewMatrix(width, height int) Matrix {
	matrix := make(Matrix, height)

	id := 1
	for i := range matrix {
		matrix[i] = make([]*Block, width)
		for j := range matrix[i] {
			matrix[i][j] = NewBlock(id, NewSemaphore())
			id++
		}
	}
	return matrix
}

func (m Matrix) GetRandomBlock() *Block {
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

func loadSequenceFromFile(filename string) ([][2]int, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var sequence [][2]int
	err = json.Unmarshal(data, &sequence)
	if err != nil {
		return nil, err
	}

	return sequence, nil
}

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/brnocorreia/concurrency/config/logger"
	"go.uber.org/zap"
)

type Block struct {
	Id       int
	Health   int
	Hit_time time.Duration
	mutex    *sync.Mutex
}

func NewBlock(id int, hit_time int, mutex *sync.Mutex) *Block {
	return &Block{
		Id:       id,
		Health:   100,
		Hit_time: time.Millisecond * time.Duration(hit_time),
		mutex:    mutex,
	}
}

func (b *Block) Hit(player *Player) bool {
	b.mutex.Lock()
	defer b.mutex.Unlock() // Garante que a operação de desbloqueio ocorra após o fim da interação com a memoria

	logger.Info("O player tenta acertar o bloco", zap.Int("playerId", player.Id), zap.Int("blockId", b.Id))

	// TODO: Adicionar contagem de pontos para o ultimo que acertou antes de morrer
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
	b.mutex.Lock()
	defer b.mutex.Unlock()
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

type Matrix [][]*Block

func NewMatrix(width, height int) Matrix {
	matrix := make(Matrix, height)

	id := 1
	for i := range matrix {
		matrix[i] = make([]*Block, width)
		for j := range matrix[i] {
			hitTime := rand.Intn(500) + 1
			matrix[i][j] = NewBlock(id, hitTime, &sync.Mutex{})
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

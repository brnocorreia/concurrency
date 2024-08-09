package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/brnocorreia/concurrency/config/logger"
	"go.uber.org/zap"
)

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

type Block struct {
	Id       int
	Health   int
	Hit_time time.Duration
	mutex    *PriorityMutex
}

func NewBlock(id int) *Block {
	var hitTime time.Duration
	if id%2 == 0 {
		hitTime = 500 * time.Millisecond // Duração de 375ms para ids pares
	} else {
		hitTime = 125 * time.Millisecond // Duração de 250ms para ids ímpares
	}
	return &Block{
		Id:       id,
		Health:   100,
		Hit_time: hitTime,
		mutex:    NewPriorityMutex(),
	}
}

func (b *Block) Hit(player *Player, lockSync chan [4]int, updates chan [4]int, x int, y int) bool {
	b.mutex.Lock(false)
	// Notifica a outra goroutine que o bloco[x][y] está sendo acertado e precisa ser lockado
	lockSync <- [4]int{player.Id, 0, x, y}

	// Ao retornar, a função dá unlock na sua matriz e notifica a outra para dar unlock também
	defer func() {
		b.mutex.Unlock(false)
		lockSync <- [4]int{player.Id, 1, x, y}
	}()

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
		// Preciso notificar a outra goroutine que o bloco[x][y] foi acertado e precisa atualizar o seu estado
		updates <- [4]int{player.Id, b.Health, x, y}
		return true
	}
	return false
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

func syncLocks(lockSync chan [4]int, matrix_1 Matrix, matrix_2 Matrix) {
	// defer wg.Done()
	// Vale lembrar:
	//   	- lock[0] = playerId
	//		- lock[1] = 0 se for operação de Lock e 1 se for operação de Unlock
	//		- lock[2] = Coordenada x do bloco
	//		- lock[3] = Coordenada y do bloco
	for lock := range lockSync {
		id, op, x, y := lock[0], lock[1], lock[2], lock[3]
		// Se for o player 1, eu preciso lockar na matrix_2
		if id == 1 {
			if op == 0 {
				// Operação de lock
				matrix_2[x][y].mutex.Lock(true)
			} else {
				// Operação de unlock
				matrix_2[x][y].mutex.Unlock(true)
			}
		} else {
			// Se for o player 2, eu preciso lockar na matrix_1
			if op == 0 {
				// Operação de lock
				matrix_1[x][y].mutex.Lock(true)
			} else {
				// Operação de unlock
				matrix_1[x][y].mutex.Unlock(true)
			}
		}
	}
}

func updateMatrix(updates chan [4]int, matrix_1 Matrix, matrix_2 Matrix) {
	// defer wg.Done()
	// Vale lembrar:
	//   	- update[0] = playerId
	//		- update[1] = O novo valor da saúde do bloco
	//		- update[2] = Coordenada x do bloco
	//		- update[3] = Coordenada y do bloco
	for update := range updates {
		id, health, x, y := update[0], update[1], update[2], update[3]
		if id == 1 {
			matrix_2[x][y].mutex.Lock(true)
			matrix_2[x][y].Health = health
			matrix_2[x][y].mutex.Unlock(true)
		} else {
			matrix_1[x][y].mutex.Lock(true)
			matrix_1[x][y].Health = health
			matrix_1[x][y].mutex.Unlock(true)
		}
	}
}

func main() {

	logger.Info("Iniciando o jogo...")

	const numAttacks = 256 // Número de ataques
	const matrixSize = 8   // Tamanho da matriz (30x30)

	logger.Info("Configuração do jogo:", zap.Int("numAttacks", numAttacks), zap.Int("matrixSize", matrixSize))

	// Cria a matriz de blocos
	matrix_1 := NewMatrix(matrixSize, matrixSize)
	matrix_2 := NewMatrix(matrixSize, matrixSize)

	logger.Info("Criando os jogadores...")
	// Cria jogadores
	player1 := NewPlayer(1, 30)
	player2 := NewPlayer(2, 30)

	var wg sync.WaitGroup

	results := make(chan string, 2)
	lockSync := make(chan [4]int, 2)
	updates := make(chan [4]int, 2)

	// Função para simular um ataque
	attack := func(player *Player, sequence [][2]int, matrix Matrix) {
		defer wg.Done()
		for _, coord := range sequence {
			x, y := coord[0], coord[1]
			block := matrix[x][y]
			block.Hit(player, lockSync, updates, x, y)
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

	// Inicia as goroutines de sincronização
	// wg.Add(1)
	go syncLocks(lockSync, matrix_1, matrix_2)

	// Inicia as goroutines de atualização
	// wg.Add(1)
	go updateMatrix(updates, matrix_1, matrix_2)

	// Inicia as goroutines para os dois jogadores
	wg.Add(2)
	init := time.Now()
	go attack(player1, sequence_1, matrix_1)
	go attack(player2, sequence_2, matrix_2)

	// Aguarda até que ambas as goroutines terminem
	wg.Wait()
	close(lockSync)
	close(updates)
	close(results)

	duration := (time.Since(init))

	logger.Info("Tempo de execução:", zap.Duration("duration", duration))

	// Imprime o estado final dos blocos
	fmt.Println("------------------------------------------------")
	fmt.Println()
	fmt.Println("Estado final da matriz 1:")
	printBlocks(matrix_1)
	fmt.Println()
	fmt.Println("------------------------------------------------")
	// Imprime o estado final dos blocos
	fmt.Println("------------------------------------------------")
	fmt.Println()
	fmt.Println("Estado final da matriz 2:")
	printBlocks(matrix_2)
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
			matrix[i][j] = NewBlock(id)
			id++
		}
	}
	return matrix
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

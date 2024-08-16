package runner

import (
	"fmt"
	"sync"
	"time"

	"github.com/brnocorreia/concurrency/internal/config/logger"
	"github.com/brnocorreia/concurrency/internal/entity"
	"github.com/brnocorreia/concurrency/internal/tools"
	"go.uber.org/zap"
)

type Runner struct {
	numAttacks  int
	matrixSize  int
	playerPower int
	sequence_1  [][2]int
	sequence_2  [][2]int
}

func NewRunner(numAttacks, matrixSize, playerPower int) *Runner {
	return &Runner{
		numAttacks:  numAttacks,
		matrixSize:  matrixSize,
		playerPower: playerPower,
	}
}

func (r *Runner) LoadSequence() (bool, error) {
	logger.Info("Carregando a sequência de ataques...")
	seq_1, err := tools.LoadSequenceFromFile("sequence_1.json")
	if err != nil {
		logger.Info("Erro ao carregar a sequência 1")
		return false, err
	}
	r.sequence_1 = seq_1

	seq_2, err := tools.LoadSequenceFromFile("sequence_2.json")
	if err != nil {
		logger.Info("Erro ao carregar a sequência 2")
		return false, err
	}
	r.sequence_2 = seq_2
	return true, nil
}

func (r *Runner) RunMutex() {

	logger.Info("Iniciando o jogo para versão MUTEX...")

	// Cria a matriz de blocos
	matrix := entity.NewMatrixMutex(r.matrixSize, r.matrixSize)

	logger.Info("Criando os jogadores...")
	// Cria jogadores
	player1 := entity.NewPlayer(1, r.playerPower)
	player2 := entity.NewPlayer(2, r.playerPower)

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
	logger.Info("Iniciando as goroutines...")
	wg.Add(2)
	init := time.Now()
	go attack(player1, r.sequence_1)
	go attack(player2, r.sequence_2)

	// Aguarda até que ambas as goroutines terminem
	wg.Wait()
	duration := (time.Since(init))

	logger.Info("Tempo de execução [MUTEX]:", zap.Duration("duration", duration))

	close(results)

	// TODO: Armazenar o estado final dos blocos num arquivo/log
	// Imprime o estado final dos blocos
	fmt.Println("------------------------------------------------")
	fmt.Println()
	fmt.Println("Estado final dos blocos:")
	entity.PrintBlocksMutex(matrix)
	fmt.Println()
	fmt.Println("------------------------------------------------")

	for result := range results {
		logger.Info(result)
		fmt.Println(result)
	}
	logger.Info("Finalizando o jogo para versão MUTEX...")
}

func (r *Runner) RunSemaphore() {
	logger.Info("Iniciando o jogo para versão SEMAPHORE...")

	// Cria a matriz de blocos
	matrix := entity.NewMatrixSemaphore(r.matrixSize, r.matrixSize)

	logger.Info("Criando os jogadores...")
	// Cria jogadores
	player1 := entity.NewPlayer(1, r.playerPower)
	player2 := entity.NewPlayer(2, r.playerPower)

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
	logger.Info("Iniciando as goroutines...")
	wg.Add(2)
	init := time.Now()
	go attack(player1, r.sequence_1)
	go attack(player2, r.sequence_2)

	// Aguarda até que ambas as goroutines terminem
	wg.Wait()
	duration := (time.Since(init))

	logger.Info("Tempo de execução [SEMAPHORE]:", zap.Duration("duration", duration))

	close(results)

	// Imprime o estado final dos blocos
	fmt.Println("------------------------------------------------")
	fmt.Println()
	fmt.Println("Estado final dos blocos:")
	entity.PrintBlocksSemaphore(matrix)
	fmt.Println()
	fmt.Println("------------------------------------------------")

	for result := range results {
		logger.Info(result)
		fmt.Println(result)
	}
	logger.Info("Finalizando o jogo para versão SEMAPHORE...")
}

func (r *Runner) RunMessages() {
	logger.Info("Iniciando o jogo para versão TROCA DE MENSAGENS...")

	// Cria a matriz de blocos
	matrix_1 := entity.NewMatrixMessage(r.matrixSize, r.matrixSize)
	matrix_2 := entity.NewMatrixMessage(r.matrixSize, r.matrixSize)

	logger.Info("Criando os jogadores...")
	// Cria jogadores
	player1 := entity.NewPlayer(1, r.playerPower)
	player2 := entity.NewPlayer(2, r.playerPower)

	var wg sync.WaitGroup

	results := make(chan string, 2)
	lockSync := make(chan [4]int, 2)
	updates := make(chan [4]int, 2)

	// Função para simular um ataque
	attack := func(player *entity.Player, sequence [][2]int, matrix entity.MatrixMessage) {
		defer wg.Done()
		for _, coord := range sequence {
			x, y := coord[0], coord[1]
			block := matrix[x][y]
			block.Hit(player, lockSync, updates, x, y)
		}

		result := fmt.Sprintf("O player %d ganhou %d pontos\n", player.Id, player.GetPoints())
		results <- result
	}

	logger.Info("Iniciando as goroutines...")
	// Inicia as goroutines de sincronização
	// wg.Add(1)
	go entity.SyncLocks(lockSync, matrix_1, matrix_2)

	// Inicia as goroutines de atualização
	// wg.Add(1)
	go entity.UpdateMatrix(updates, matrix_1, matrix_2)

	// Inicia as goroutines para os dois jogadores
	wg.Add(2)
	init := time.Now()
	go attack(player1, r.sequence_1, matrix_1)
	go attack(player2, r.sequence_2, matrix_2)

	// Aguarda até que ambas as goroutines terminem
	wg.Wait()
	close(lockSync)
	close(updates)
	close(results)

	duration := (time.Since(init))

	logger.Info("Tempo de execução [TROCA DE MENSAGENS]:", zap.Duration("duration", duration))

	// Imprime o estado final dos blocos
	fmt.Println("------------------------------------------------")
	fmt.Println()
	fmt.Println("Estado final da matriz 1:")
	entity.PrintBlocksMessage(matrix_1)
	fmt.Println()
	fmt.Println("------------------------------------------------")
	// Imprime o estado final dos blocos
	fmt.Println("------------------------------------------------")
	fmt.Println()
	fmt.Println("Estado final da matriz 2:")
	entity.PrintBlocksMessage(matrix_2)
	fmt.Println()
	fmt.Println("------------------------------------------------")

	for result := range results {
		logger.Info(result)
		fmt.Println(result)
	}
	logger.Info("Finalizando o jogo para versão TROCA DE MENSAGENS...")
}

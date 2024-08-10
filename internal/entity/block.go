package entity

import (
	"fmt"
	"sync"
	"time"

	"github.com/brnocorreia/concurrency/internal/config/logger"
	"github.com/brnocorreia/concurrency/internal/tools"
	"go.uber.org/zap"
)

// Implementação dos blocos para MUTEX
type BlockMutex struct {
	Id       int
	Health   int
	Hit_time time.Duration
	mutex    *sync.Mutex
}

func NewBlockMutex(id int, mutex *sync.Mutex) *BlockMutex {
	var hitTime time.Duration
	if id%2 == 0 {
		hitTime = 500 * time.Millisecond // Duração de 375ms para ids pares
	} else {
		hitTime = 125 * time.Millisecond // Duração de 250ms para ids ímpares
	}
	return &BlockMutex{
		Id:       id,
		Health:   100,
		Hit_time: hitTime,
		mutex:    mutex,
	}
}

func (b *BlockMutex) Hit(player *Player) bool {
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

func (b *BlockMutex) IsAlive() bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.Health > 0
}

func (b *BlockMutex) String() string {
	return fmt.Sprintf("ID=%d, Health=%d, HitTime=%v", b.Id, b.Health, b.Hit_time)
}

// -------------------------------------------------------------------------------------

// Implementação dos blocos para SEMAPHORE
type BlockSemaphore struct {
	Id        int
	Health    int
	Hit_time  time.Duration
	semaphore *tools.Semaphore
}

func NewBlockSemaphore(id int, semaphore *tools.Semaphore) *BlockSemaphore {
	var hitTime time.Duration
	if id%2 == 0 {
		hitTime = 500 * time.Millisecond // Duração de 375ms para ids pares
	} else {
		hitTime = 125 * time.Millisecond // Duração de 125ms para ids ímpares
	}

	return &BlockSemaphore{
		Id:        id,
		Health:    100,
		Hit_time:  hitTime,
		semaphore: semaphore,
	}
}

func (b *BlockSemaphore) Hit(player *Player) bool {
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

// -------------------------------------------------------------------------------------

// Implementação dos blocos para troca de mensagens
type BlockMessage struct {
	Id       int
	Health   int
	Hit_time time.Duration
	mutex    *tools.PriorityMutex
}

func NewBlockMessage(id int) *BlockMessage {
	var hitTime time.Duration
	if id%2 == 0 {
		hitTime = 500 * time.Millisecond // Duração de 375ms para ids pares
	} else {
		hitTime = 125 * time.Millisecond // Duração de 250ms para ids ímpares
	}
	return &BlockMessage{
		Id:       id,
		Health:   100,
		Hit_time: hitTime,
		mutex:    tools.NewPriorityMutex(),
	}
}

func (b *BlockMessage) Hit(player *Player, lockSync chan [4]int, updates chan [4]int, x int, y int) bool {
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

func SyncLocks(lockSync chan [4]int, matrix_1 MatrixMessage, matrix_2 MatrixMessage) {
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

func UpdateMatrix(updates chan [4]int, matrix_1 MatrixMessage, matrix_2 MatrixMessage) {
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

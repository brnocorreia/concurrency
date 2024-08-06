package entity

import (
	"fmt"
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

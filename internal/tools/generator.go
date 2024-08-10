package tools

import (
	"encoding/json"
	"fmt"
	"os"

	"math/rand"

	"github.com/brnocorreia/concurrency/internal/config/logger"
	"go.uber.org/zap"
)

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

func saveSequenceToFile(sequence [][2]int, filename string) error {
	data, err := json.Marshal(sequence)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// Função para ler a sequência de ataques de um arquivo
func LoadSequenceFromFile(filename string) ([][2]int, error) {
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

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func Generate(size, numAttacks int) (bool, error) {
	for i := 1; i <= 2; i++ {
		sequence := generateAttackSequence(size, numAttacks)
		filename := fmt.Sprintf("sequence_%d.json", i)
		err := saveSequenceToFile(sequence, filename)
		if err != nil {
			logger.Error("Erro ao salvar a sequência", err)
			return false, err
		}
		logger.Info("Sequência carregada em:", zap.String("filename", filename))
	}
	return true, nil
}

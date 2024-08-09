package main

import (
	"encoding/json"
	"fmt"
	"os"

	"math/rand"
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

func main() {
	var size, numAttacks int
	fmt.Print("Digite o tamanho do tabuleiro: ")
	fmt.Scanln(&size)
	fmt.Print("Digite o número de ataques: ")
	fmt.Scanln(&numAttacks)

	for i := 1; i <= 2; i++ {
		sequence := generateAttackSequence(size, numAttacks)
		filename := fmt.Sprintf("sequence_%d.json", i)
		err := saveSequenceToFile(sequence, filename)
		if err != nil {
			fmt.Println("Erro ao salvar a sequência:", err)
			return
		}
		fmt.Println("Sequência carregada em:", filename)
	}
}

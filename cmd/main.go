package main

import (
	"fmt"
	"os"

	"github.com/brnocorreia/concurrency/internal/config/logger"
	"github.com/brnocorreia/concurrency/internal/runner"
	"github.com/brnocorreia/concurrency/internal/tools"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// Semaphore -> https://medium.com/@deckarep/gos-extended-concurrency-semaphores-part-1-5eeabfa351ce

var (
	numAttacks  int
	matrixSize  int
	playerPower int
	mode        string
	regenerate  bool
)

var rootCmd = &cobra.Command{
	Use:   "concurrency",
	Short: "Concurrency is a game developed for study purposes under UFBA's MATA58 course",
	Long:  `Simple and fast, Concurrency explores all that Go has to offer when it comes to concurrency and parallelism.`,
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the game",
	Run: func(cmd *cobra.Command, args []string) {
		runner := runner.NewRunner(numAttacks, matrixSize, playerPower)

		if regenerate {
			logger.Info("Regenerating attack sequences...")
			_, err := tools.Generate(matrixSize, numAttacks)
			if err != nil {
				logger.Info("Error generating sequences:", zap.Error(err))
				os.Exit(1)
			}
		}

		if !tools.FileExists("sequence_1.json") || !tools.FileExists("sequence_2.json") {
			logger.Info("Attack sequences not found, generating...")
			_, err := tools.Generate(matrixSize, numAttacks)
			if err != nil {
				logger.Info("Error generating sequences:", zap.Error(err))
				os.Exit(1)
			}
		}

		_, err := runner.LoadSequence()
		if err != nil {
			logger.Info("Error loading sequences:", zap.Error(err))
			os.Exit(1)
		}

		switch mode {
		case "mutex":
			runner.RunMutex()
		case "semaphore":
			runner.RunSemaphore()
		case "messages":
			runner.RunMessages()
		case "all":
			runner.RunMutex()
			runner.RunSemaphore()
			runner.RunMessages()
		default:
			logger.Info("Invalid mode:", zap.String("mode", mode))
			os.Exit(1)
		}
	},
}

func init() {
	runCmd.Flags().IntVarP(&numAttacks, "attacks", "a", 256, "Number of attacks")
	runCmd.Flags().IntVarP(&matrixSize, "size", "s", 8, "Matrix size")
	runCmd.Flags().IntVarP(&playerPower, "power", "p", 30, "Player power")
	runCmd.Flags().StringVarP(&mode, "mode", "m", "all", "Execution mode (mutex, semaphore, or messages)")
	runCmd.Flags().BoolVarP(&regenerate, "regenerate", "r", false, "Regenerate attack sequences")

	rootCmd.AddCommand(runCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

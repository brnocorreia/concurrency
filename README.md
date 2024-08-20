# Concurrency Game

Concurrency is a game developed for study purposes under UFBA's MATA58 course. Simple and fast, it explores all that Go has to offer when it comes to concurrency and parallelism.

## Getting Started

### Prerequisites

- Nothing! :laughing: :laughing: I'm kidding, you just need a good terminal.

### Installation

1. Clone and navigate to the repository:

```console
git clone https://github.com/brnocorreia/concurrency.git && cd concurrency
```

### Running

- You can run the game using the go command (make sure you have Go installed in your machine):

```console
go run cmd/main.go
```

- Or you can use the **recommended** way, which is to use the binary. Make sure to select the correct binary for your OS:
- We will use the `concurrency-linux-amd64` binary in this guide. You can clone the entire repo or just download the binary on the table below.

```console
./bin/concurrency-linux-amd64
```

#### Available binaries:

| OS      | ARCH  | Filename                      | Status | Download                                                                                                      |
| ------- | ----- | ----------------------------- | ------ | ------------------------------------------------------------------------------------------------------------- |
| Linux   | amd64 | concurrency-linux-amd64       | ✅     | [Download](https://github.com/brnocorreia/concurrency/releases/latest/download/concurrency-linux-amd64)       |
| Linux   | arm64 | concurrency-linux-arm64       | ✅     | [Download](https://github.com/brnocorreia/concurrency/releases/latest/download/concurrency-linux-arm64)       |
| Windows | amd64 | concurrency-windows-amd64.exe | ✅     | [Download](https://github.com/brnocorreia/concurrency/releases/latest/download/concurrency-windows-amd64.exe) |
| macOS   | amd64 | concurrency-darwin-amd64      | ✅     | [Download](https://github.com/brnocorreia/concurrency/releases/latest/download/concurrency-darwin-amd64)      |

## Usage

### Run

- The game should be run using the `run` command, as shown below.

```console
./bin/concurrency-linux-amd64 run
```

- You can use flags to change the game's parameters like the number of attacks of each player, the matrix size, and the player's power (how much damage they can deal with one hit).

- Feel free to use the `-h` flag to see the available options, but basically, you can use the following:

```console
./bin/concurrency-linux-amd64 run -h

Usage:
  concurrency run [flags]

Flags:
  -a, --attacks int   Number of attacks (default 256)
  -h, --help          help for run
  -m, --mode string   Execution mode (mutex, semaphore, or messages) (default "all")
  -p, --power int     Player power (default 30)
  -r, --regenerate    Regenerate attack sequences
  -s, --size int      Matrix size (default 8)
```

#### Examples

- Run the game with the default parameters:

```console
./bin/concurrency-linux-amd64 run
```

- Run the game with a different number of attacks and matrix size:

```console
./bin/concurrency-linux-amd64 run -a 512 -s 16
```

- Run the game in a specific mode, this is the most important flag:

```console
./bin/concurrency-linux-amd64 run -m semaphore
```

## Game Modes

The game supports multiple execution modes:

- `mutex`: Uses mutual exclusion for concurrency control
- `semaphore`: Uses semaphores for concurrency control
- `messages`: Uses message passing for concurrency control
- `all`: Runs all modes sequentially

## Outputs

- You can check the attack sequences in the sequence_1.json and sequence_2.json files.
- You can check the default logs in the log.log file and the results (player points and final stage of the game matrix) in the results.json file.

## Additional Information

- The game automatically generates attack sequences if they don't exist.
- Use the `--regenerate` flag to force regeneration of attack sequences.
- The game will exit with an error message if an invalid mode is specified.

## Acknowledgments

- UFBA's MATA58 course for inspiring this project, this was the main job of the course.
- The Go community for providing excellent tools and libraries

# Concurrency Game

Concurrency é um jogo via CLI desenvolvido para propósitos educacionais da matéria de Sistemas Operacionais (MATA58). Simples e rápido, concurrency explora tudo que o Golang tem a oferecer quando o tema é código concorrente e paralelo.

## Como jogar

### Pré-requisitos

- Apenas o seu terminal :laughing:
- Go 1.22 ou superior, se não quiser ir pela maneira recomendada.

### Instalação

1. Clone e navegue até o repositório:

```console
git clone https://github.com/brnocorreia/concurrency.git && cd concurrency
```

### Rodando o jogo

- Você pode rodar o jogo usando o comando `go run` (certifique-se de que o Go está instalado na sua máquina):

```console
go run cmd/main.go
```

- Ou você pode usar a forma **recomendada** de rodar o jogo, que é através do binário. Certifique-se de escolher o binário apropriado para o seu sistema operacional:
- Iremos usar o binário `concurrency-linux-amd64` nesses guia.

```console
./bin/concurrency-linux-amd64
```

#### Binários disponíveis:

| OS      | ARCH  | Nome do arquivo               | Status | Download                                                                                                      |
| ------- | ----- | ----------------------------- | ------ | ------------------------------------------------------------------------------------------------------------- |
| Linux   | amd64 | concurrency-linux-amd64       | ✅     | [Download](https://github.com/brnocorreia/concurrency/releases/latest/download/concurrency-linux-amd64)       |
| Linux   | arm64 | concurrency-linux-arm64       | ✅     | [Download](https://github.com/brnocorreia/concurrency/releases/latest/download/concurrency-linux-arm64)       |
| Windows | amd64 | concurrency-windows-amd64.exe | ✅     | [Download](https://github.com/brnocorreia/concurrency/releases/latest/download/concurrency-windows-amd64.exe) |
| macOS   | amd64 | concurrency-darwin-amd64      | ✅     | [Download](https://github.com/brnocorreia/concurrency/releases/latest/download/concurrency-darwin-amd64)      |

## Uso

### Rodar

- O jogo deve ser executado usando o comando `run`, como mostrado abaixo.

```console
./bin/concurrency-linux-amd64 run
```

- Você pode usar flags para alterar os parâmetros do jogo, como o número de ataques de cada jogador, o tamanho da matriz, e a força do jogador (quanto dano ele pode causar com um golpe).

- Sinta-se livre para usar o sinalizador `-h` para ver as opções disponíveis, mas básicamente, você pode usar o seguinte:

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

#### Exemplos

- Rodar o jogo com os parâmetros padrão:

```console
./bin/concurrency-linux-amd64 run
```

- Rodar o jogo com um número diferente de ataques e tamanho da matriz:

```console
./bin/concurrency-linux-amd64 run -a 512 -s 16
```

- Rodar o jogo em um modo específico, este é o mais importante sinalizador:

```console
./bin/concurrency-linux-amd64 run -m semaphore
```

## Modos de Execução

O jogo suporta vários modos de execução:

- `mutex`: Usa exclusão mutua para controle de concorrência
- `semaphore`: Usa semáforos para controle de concorrência
- `messages`: Usa passagem de mensagens para controle de concorrência
- `all`: Executa todos os modos sequencialmente

## Saídas

- Você pode verificar as sequências de ataques em arquivos sequence_1.json e sequence_2.json.
- Você pode verificar os logs padrão em arquivo log.log e os resultados (pontos do jogador e etapa final da matriz do jogo) em arquivo results.json.

## Informações Adicionais

- O jogo gera automaticamente sequências de ataques se elas não existirem.
- Use o sinalizador `--regenerate` para forçar a regeneração de sequências de ataques.
- O jogo sairá com uma mensagem de erro se um modo inválido for especificado.

## Agradecimentos

- A matéria de Sistemas Operacionais (MATA58) da UFBA para inspirar este projeto, esse trabalho foi o principal trabalho do curso.
- A comunidade Go para fornecer excelentes ferramentas e bibliotecas

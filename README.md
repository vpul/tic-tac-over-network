# Tic-Tac-Over-Network

Networked two-player tic-tac-toe written in Go. The server uses TCP with newline-delimited JSON, pairs clients in FIFO order, validates moves server-side, and keeps each game independent.

## Requirements

- Go 1.25 or newer

## Run

From the repository root, start the server:

```sh
go run ./cmd/server
```

The server listens on `:8080` by default. Use `-addr` to choose another address:

```sh
go run ./cmd/server -addr :9090
```

In two separate terminals, start one client per player:

```sh
go run ./cmd/client -addr localhost:8080
```

Clients wait until another client connects, then are paired automatically. Enter a cell number from `1` to `9` to make a move:

```text
1 | 2 | 3
--+---+--
4 | 5 | 6
--+---+--
7 | 8 | 9
```

Enter `q` or `quit` to leave. If a player leaves, the opponent is notified and both clients exit.

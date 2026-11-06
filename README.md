# cli-tools-go

A collection of small, independent command-line tools written in Go.

Each tool lives in its own directory with its own Go module.

## Tools

| Tool | Description |
|------|-------------|
| [word-counter](word-counter) | Count words, lines or bytes from stdin |

## Build

```sh
cd <tool>
go build -o ../bin/<tool> .
```

## Usage

```sh
echo "hello world" | bin/word-counter
echo "hello world" | bin/word-counter -b
```

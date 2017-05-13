tcp-port-forward
================

## Usage

```
$ go run main.go -h
a tcp port-forwarder written in golang

Usage:
  tcp-port-forward [command]

Available Commands:
  help        Help about any command
  local       the local port of the port-forwarder
  remote      init remote part of the port-forwarder

Flags:
  -f, --from string   the local address of the proxy (default "localhost:5000")
  -h, --help          help for tcp-port-forward
  -t, --to string     the upstream address (default "localhost:10000")

Use "tcp-port-forward [command] --help" for more information about a command.
```

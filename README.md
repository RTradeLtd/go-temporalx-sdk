# go-temporalx-sdk

This repository is an open-source fork of TemporalX's command line client, with all server functionality remove. It provides a feature complete CLI client for using **all** TemporalX functionality, and can also help you to build other applicatinos ontop of the gRPC API.

# installation

to build a copy of this locally you'll need go 1.13+, and have downloaded all the dependencies using `go mod download`. After you can then built the cli tool with a simple `make`.

# usage

To get an overview you can simply invoke the binary with `tex-cli` which should display output similar to:

```
NAME:
   tex-cli - TemporalX command-line management tool

USAGE:
   tex-cli [global options] command [command options] [arguments...]

VERSION:
   v1.0.0

AUTHORS:
   Alex Trottier <postables@rtradetechnologies.com>
   George Xie <georgex@rtradetechnologies.com>

COMMANDS:
   client   gRPC client subcommands
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --bootstrap, --bp  bootstrap against public ipfs (default: false)
   --help, -h         show help (default: false)
   --version, -v      print the version (default: false)
```

You can then access all client commands with `tex-cli client` which should display the following output:

```
NAME:
   tex-cli client - gRPC client subcommands

USAGE:
   tex-cli client command [command options] [arguments...]

DESCRIPTION:
   Enables access to a rich gRPC client library

COMMANDS:
   node      node management commands
   file      file upload/download commands
   extras    node extras management
   pubsub    pubsub commands
   namesys   namesys commands
   keystore  keystore commands
   api       low level api maintenance commands
   admin     admin commands
   help, h   Shows a list of commands or help for one command

OPTIONS:
   --endpoint.address value, --ea value  temporalx endpoint address (default: "127.0.0.1")
   --insecure                            enable insecure connections to temporalx (default: true)
   --help, -h                            show help (default: false)
   --version, -v                         print the version (default: false)
```

If you want to upload a file to the public TemporalX endpoint you can do:

```shell
$> tex-cli client --ea xapi.temporal.cloud:9090 file upload --fn /path/to/file
```

So for example lets upload the binary itself (this will require having put the `tex-cli` in to your `$PATH` env variable):

```shell
$> tex-cli client --ea xapi.temporal.cloud:9090 file upload --fn $(which tex-cli)
hash of file: bafybeihqruaz4k3iux43vysporzcpgpkkqihnavm3zxwcj7o5zgbqgu77a
```

Alternatively if you want to view the progress on your upload:

```shell
$> ./tex-cli client --ea xapi.temporal.cloud:9090 file upload --fn tex-cli --pp
 100% |████████████████████████████████████████|  [2s:0s]
hash of file: bafybeihqruaz4k3iux43vysporzcpgpkkqihnavm3zxwcj7o5zgbqgu77a
```


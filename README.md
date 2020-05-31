# Typing Test

Typing Test is a simple web application to implement WebSocket functionality. The server code is written in Golang.

## Installation

Use the git clone  to download the repository

```bash
git clone https://github.com/nishanksrj/TypingTest.git
```


Before running, make sure you have Go installed in your system.

[Click to visit the Installation page](https://golang.org/doc/install)

&nbsp;

After installing Go, Install the following dependencies

```bash
go get github.com/gorilla/websocket
go get github.com/dgraph-io/badger/v2
```

## Usage

Run the following bash command to start the server. The server will listen for the requests on localhost at port 8000.
```bash
go run main.go
```

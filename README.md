# Ruthenium

## Description

The blockchain motivating to do local business and to not capitalize tokens.

## Installation

* Clone the project
* Build the project (**requires [go](https://go.dev/dl/), version 1.17 at least)**

## Usage

* In src/blockchain_server, run:
    * go run main.go
    * go run main.go -port=5001
* In src/wallet_server, run:
    * go run main.go -host-ip=<your ip>
    * go run main.go -port=8081 -host-ip=<your ip> -host-port=5001
* In two separated tabs using Mozilla Firefox web browser, go to:
    * http://localhost:8080
    * http://localhost:8081
* Start sending money!

## Contributing

Contribution is not open yet

## Authors and acknowledgment

For a tutorial to create a first blockchain in go, thanks to:

* [Yuko Sakai](https://www.udemy.com/user/myeigoworld/)
* [Jun Sakai](https://udemy.com/user/jun-sakai/)

For an incredibly simple [P2P go library](https://github.com/leprosus/golang-p2p), thanks to:

* [Denis Korolev](https://github.com/leprosus)

## License

![img.png](img.png)

http://unlicense.org/

## Project status

Still under heavy development

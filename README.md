# Ruthenium

## Description

The blockchain motivating to do local business and to not capitalize tokens.

## Installation

* Clone the project
* Build the project (**requires [go](https://go.dev/dl/), version 1.17 at least)**

## Usage

* Open your DNS port 8106
* At root level (ruthenium folder), run:
    * go run src/ui/main.go -host-ip=<your external IP address>
* Using a web browser, go to:
    * http://localhost:8080
    * store your public key, private key and wallet address
* In src/blockchain_server, run:
    * go run src/node/main.go -public-key=<your public key> -private-key=<your private key>
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

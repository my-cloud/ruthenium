# Ruthenium

## Description

The blockchain motivating to not capitalize tokens.

## Installation

* Clone the project
* Build the project (**requires [go](https://go.dev/dl/), version 1.17 at least)**

## Usage

* Open your DNS port 8106
* At root level (ruthenium folder), run:
    * go run src/ui/main.go -host-ip={your external IP address}
* Using a web browser, go to:
    * http://localhost:8080
    * store your private key
* In src/blockchain_server, run:
    * go run src/node/main.go -private-key={your private key}
* Start sending money!

## Authors and acknowledgment

For a tutorial to create a first blockchain in go, thanks to:

* [Yuko Sakai](https://www.udemy.com/user/myeigoworld/)
* [Jun Sakai](https://udemy.com/user/jun-sakai/)

## License

![img.png](doc/img.png)

http://unlicense.org/

## Project status

The main principles have been implemented.  
Now it needs a lot of refactoring and tests to improve maintainability and
reliability.

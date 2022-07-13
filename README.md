# Ruthenium

## Description

The blockchain motivating to do local business and to not capitalize money.

## Installation

go 1.16 is recommended

## Usage

* Clone this project
* In src/blockchain_server, run:
    * go run main.go
    * go run main.go -port=5001
* In src/wallet_server, run:
    * go run main.go
    * go run main.go -port=8081 -gateway=http://127.0.0.1:5001
* In two separated tabs using Mozilla Firefox web browser, go to:
    * http://localhost:8080
    * http://localhost:8081
* Start sending money!

## Contributing

Contribution is not open yet

## Authors and acknowledgment

For a first blockchain in go tutorial, thanks to:

* [Yuko Sakai](https://www.udemy.com/user/myeigoworld/)
* [Jun Sakai](https://udemy.com/user/jun-sakai/)

## License

![img.png](img.png)

http://unlicense.org/

## Project status

**Done**: Blockchain created using REST APIs

**WIP**: Replace REST by P2P

**TODO**: Customize the blockchain

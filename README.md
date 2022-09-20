# Ruthenium

## Description

The blockchain motivating to not capitalize tokens.

## Prerequisites
* Building from sources
  * Install [go](https://go.dev/dl/) version 1.17 or higher

* Using docker image
  * Install [docker](https://www.docker.com/)

## Installation
* Building from sources
  * Clone the project: `https://github.com/my-cloud/ruthenium.git`

* Using docker image
  * Pull the image: `sudo docker pull ghcr.io/my-cloud/ruthenium:main`

## Usage

* Open your DNS port 8106
* To be able to validate blocks, you need be registered in the [Proof of Humanity](https://app.proofofhumanity.id/) registry with an Ethereum wallet address for which you are the owner of `<your private key>`

* Building from sources
  * At root level (ruthenium folder), run the node:
    ```
    go run src/node/main.go -private-key=<your private key>
    ```
  * At root level (ruthenium folder), run the ui:
    ```
    go run src/ui/main.go -host-ip=<your external IP address> -private-key=<your private key>
    ```

* Using docker image
  * Run the node:
    ```
    sudo docker run -ti ghcr.io/my-cloud/ruthenium:main \app\node -host-ip=<your external IP address> -private-key=<your private key>
    ```
  * Run the ui:
    ```
    sudo docker run -ti ghcr.io/my-cloud/ruthenium:main \app\ui -private-key=<your private key>
    ```

* Using a web browser, go to:
  * http://localhost:8080
  * if not provided, the private key will be generated, then you need to securely store it.

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

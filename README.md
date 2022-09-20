# Ruthenium

## Description
The blockchain motivating to not capitalize tokens.

For further information, see [the dedicated wiki page](https://github.com/my-cloud/ruthenium/wiki/Home)

## Prerequisites
* Using sources (option A)
  * Install [go](https://go.dev/dl/) version 1.17 or higher
* Using docker image (option B)
  * Install [docker](https://www.docker.com/)
* Open your DNS port 8106
* To be able to validate blocks, you need be registered in the [Proof of Humanity](https://app.proofofhumanity.id/) registry with an Ethereum wallet address for which you are the owner of `<your private key>`

For further information, see [the dedicated wiki page](https://github.com/my-cloud/ruthenium/wiki/Usage#prerequisites)

## Installation
* Using sources (option A)
  * Clone the project: `git clone https://github.com/my-cloud/ruthenium.git`
* Using docker image (option B)
  * Pull the image: `sudo docker pull ghcr.io/my-cloud/ruthenium:main`

For further information, see [the dedicated wiki page](https://github.com/my-cloud/ruthenium/wiki/Usage#installation)

## Launch
* Using sources (option A)
  * At root level (ruthenium folder), run the node:
    ```
    go run src/node/main.go -private-key=<your private key>
    ```
  * At root level (ruthenium folder), run the ui:
    ```
    go run src/ui/main.go -host-ip=<your external IP address> -private-key=<your private key>
    ```
* Using docker image (option B)
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

For further information, see [the dedicated wiki page](https://github.com/my-cloud/ruthenium/wiki/Usage#launch)

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

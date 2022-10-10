# Ruthenium

## Description
The Ruthenium blockchain.

This README contains essential information for a quick start. You will find a detailed description of the project in [the dedicated wiki page](https://github.com/my-cloud/ruthenium/wiki/Home). If you want to know what reasons led to create this blockchain, you can directly dive into the [Ruthenium whitepaper](https://github.com/my-cloud/ruthenium/wiki/Whitepaper). 

## Usage
There are two ways to easily use the blockchain. You can either use your own build from [sources](https://github.com/my-cloud/ruthenium/releases) (Option A) or use a docker image provided in the [repository packages](https://github.com/my-cloud/ruthenium/pkgs/container/ruthenium) (Option B).

### Prerequisites
* Option A: using sources
  * [Go](https://go.dev/dl/) >= 1.17
* Option B: using docker image
  * [Docker](https://www.docker.com/)
* Your DNS port 8106 must be open.
* You must be registered in the [Proof of Humanity](https://app.proofofhumanity.id/) registry with an Ethereum wallet address for which you are the owner of the `<private key>`.

### Installation
* Option A: using sources
  * Download the sources archive:
    ```
    https://github.com/my-cloud/ruthenium/releases/latest
    ```
* Option B: using docker image
  * Pull the image:
    ```
    sudo docker pull ghcr.io/my-cloud/ruthenium:main
    ```

### Launch
* Option A: using sources
  * Extract files from the sources archive
  * At root level (ruthenium folder), run the node:
    ```
    go run src/node/main.go -private-key=<private key>
    ```
  * At root level (ruthenium folder), run the ui:
    ```
    go run src/ui/main.go -host-ip=<your external IP address> -private-key=<private key>
    ```
* Option B: using docker image
  * Run the node:
    ```
    sudo docker run -p 8106:8106 -ti ghcr.io/my-cloud/ruthenium:main \app\node -host-ip=<your external IP address> -private-key=<private key>
    ```
  * Run the ui:
    ```
    sudo docker run -p 8080:8080 -ti ghcr.io/my-cloud/ruthenium:main \app\ui -private-key=<private key>
    ```
* Using a web browser, go to:
  * http://localhost:8080
  * if not provided, the private key will be generated, then you need to securely store it.

For further details concerning the usage, see [the dedicated wiki page](https://github.com/my-cloud/ruthenium/wiki/Usage)

## Contributing
Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are greatly appreciated.

If you have a suggestion that would make this better, please do not hesitate to [report a new bug](https://github.com/my-cloud/ruthenium/issues/new?assignees=&labels=bug&template=bug_report.md&title=) or [request a new feature](https://github.com/my-cloud/ruthenium/issues/new?assignees=&labels=enhancement&template=feature_request.md&title=). Of course, you are welcome to fork the repository and create a pull request. In any case, please let's take a look at the [contributing](https://github.com/my-cloud/ruthenium/blob/dev/CONTRIBUTING.md) file.

[⭐](https://github.com/login?return_to=%2Fmy-cloud%2Fruthenium) Don't forget to give the project a [star](https://github.com/login?return_to=%2Fmy-cloud%2Fruthenium)! Thanks again!

## Contact
[Jérémy Pansier](https://github.com/JeremyPansier) - jpansier@my-cloud.me

Project Link: https://github.com/my-cloud/ruthenium

## Acknowledgment
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

# Ruthenium
[![Release](https://img.shields.io/github/release/my-cloud/ruthenium?logo=github)](https://github.com/my-cloud/ruthenium/releases)

## Description
[![Wiki](https://img.shields.io/badge/wiki-home-blue?logo=github)](https://github.com/my-cloud/ruthenium/wiki)
[![Wiki](https://img.shields.io/badge/wiki-whitepaper-blue?logo=github)](https://github.com/my-cloud/ruthenium/wiki/Whitepaper)

The Ruthenium blockchain protocol.

This README contains essential information for a quick start. You will find a detailed description of the project in the [wiki](https://github.com/my-cloud/ruthenium/wiki/Home). If you want to know what reasons led to create this blockchain, you can directly dive into the Ruthenium [whitepaper](https://github.com/my-cloud/ruthenium/wiki/Whitepaper). 

## Usage
[![Build](https://img.shields.io/github/workflow/status/my-cloud/ruthenium/Build?logo=github)](https://github.com/my-cloud/ruthenium/actions?query=workflow%3ABuild+event%3Apush+branch%3Amain)

There are two ways to use the Ruthenium blockchain. You can either use your own build from [sources](https://github.com/my-cloud/ruthenium/releases) (Option A) or use a docker image provided in the [repository packages](https://github.com/my-cloud/ruthenium/pkgs/container/ruthenium) (Option B).

### Prerequisites
* Option A: using sources
  * [![Go](https://img.shields.io/github/go-mod/go-version/my-cloud/ruthenium?logo=go)](https://go.dev/dl/)
* Option B: using docker image
  * [![Docker](https://img.shields.io/badge/docker-grey?logo=docker)](https://www.docker.com/)
* Your DNS port 8106 must be open.
* You must be registered in the [Proof of Humanity](https://github.com/my-cloud/ruthenium/Whitepaper#proof-of-humanity) registry with an Ethereum wallet address for which you are the owner of the `<private key>`.

### Installation
* Option A: using sources
  * Download the sources archive:
    ```
    https://github.com/my-cloud/ruthenium/releases/latest
    ```
* Option B: using docker image
  * Pull the image:
    ```
    sudo docker pull ghcr.io/my-cloud/ruthenium:latest
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
    sudo docker run -p 8106:8106 -ti ghcr.io/my-cloud/ruthenium:latest \app\node -host-ip=<your external IP address> -private-key=<private key>
    ```
  * Run the ui:
    ```
    sudo docker run -p 8080:8080 -ti ghcr.io/my-cloud/ruthenium:latest \app\ui -private-key=<private key>
    ```
* Using a web browser, go to:
  * http://localhost:8080
  * if not provided, the private key will be generated, then you need to securely store it.

For further details concerning the usage, see [the dedicated wiki page](https://github.com/my-cloud/ruthenium/wiki/Usage)

## Contributing
[![GitHub forks](https://img.shields.io/github/forks/my-cloud/ruthenium?logo=github)](https://github.com/my-cloud/ruthenium/fork)
[![GitHub Repo stars](https://img.shields.io/github/stars/my-cloud/ruthenium?logo=github)](https://github.com/my-cloud/ruthenium)

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are greatly appreciated.

If you have a suggestion that would make this better, please do not hesitate to [report a new bug](https://github.com/my-cloud/ruthenium/issues/new?assignees=&labels=bug&template=bug_report.md&title=) or [request a new feature](https://github.com/my-cloud/ruthenium/issues/new?assignees=&labels=enhancement&template=feature_request.md&title=). Of course, you are welcome to fork the repository and create a pull request. In any case, please let's take a look at the [contributing](CONTRIBUTING.md) file.

⭐ Don't forget to give the project a [star](https://docs.github.com/en/get-started/exploring-projects-on-github/saving-repositories-with-stars)! Thanks again!

## Contact
[![Project](https://img.shields.io/badge/sources-repository-blue?logo=github)](https://github.com/my-cloud/ruthenium)

Founder: [Jérémy Pansier](https://github.com/JeremyPansier) - jpansier@my-cloud.me

## Authors and Acknowledgments
[![Contributors](https://img.shields.io/github/contributors/my-cloud/ruthenium?logo=github)](https://github.com/my-cloud/ruthenium/graphs/contributors)

Special thanks to [Gwenall Pansier](https://github.com/Gwenall) who contributed since the early developments.

For a [tutorial to create a first blockchain in Go][1], thanks to [Yuko Sakai][2] & [Jun Sakai][3].

## License
[![License](https://img.shields.io/github/license/my-cloud/ruthenium?label=⚖&nbsp;license)](LICENSE)

![license.png](doc/license.png)

## Project status
[![GitHub commit activity](https://img.shields.io/github/commit-activity/m/my-cloud/ruthenium?logo=github)](https://github.com/my-cloud/ruthenium/commits/main)
[![GitHub milestones](https://img.shields.io/github/milestones/open/my-cloud/ruthenium?logo=github)](https://github.com/my-cloud/ruthenium/milestones)
[![GitHub milestones](https://img.shields.io/github/milestones/closed/my-cloud/ruthenium?logo=github)](https://github.com/my-cloud/ruthenium/milestones)

[![Maintainability](https://sonarcloud.io/api/project_badges/measure?project=my-cloud_ruthenium&metric=sqale_rating)](https://sonarcloud.io/component_measures?id=my-cloud_ruthenium&metric=sqale_rating)
[![Security](https://sonarcloud.io/api/project_badges/measure?project=my-cloud_ruthenium&metric=security_rating)](https://sonarcloud.io/component_measures?id=my-cloud_ruthenium&metric=security_rating)
[![Reliability](https://sonarcloud.io/api/project_badges/measure?project=my-cloud_ruthenium&metric=reliability_rating)](https://sonarcloud.io/component_measures?id=my-cloud_ruthenium&metric=reliability_rating)
[![Coverage](https://img.shields.io/sonar/coverage/my-cloud_ruthenium/main?logo=sonarcloud&server=https%3A%2F%2Fsonarcloud.io)](https://sonarcloud.io/component_measures?id=my-cloud_ruthenium&metric=coverage)

The main principles have been implemented.

Now it needs a lot of refactoring and tests to improve maintainability and
reliability.

[1]: https://www.udemy.com/course/golang-how-to-build-a-blockchain-in-go/ "Udemy tutorial to build a blockchain in Go"
[2]: https://www.udemy.com/user/myeigoworld/ "Yuko Sakai LinkedIn profile"
[3]: https://udemy.com/user/jun-sakai/ "Jun Sakai LinkedIn profile"

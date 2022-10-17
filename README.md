# Ruthenium

[![Build Status](https://github.com/my-cloud/ruthenium/workflows/Go/badge.svg?branch=main)](https://github.com/my-cloud/ruthenium/actions?query=workflow%3AGo+event%3Apush+branch%3Amain)
[![Quality Gate](https://sonarcloud.io/api/project_badges/measure?project=my-cloud_ruthenium&metric=alert_status)](https://sonarcloud.io/project/overview?id=my-cloud_ruthenium)
[![Maintainability](https://sonarcloud.io/api/project_badges/measure?project=my-cloud_ruthenium&metric=sqale_rating)](https://sonarcloud.io/component_measures?id=my-cloud_ruthenium&metric=sqale_rating)
[![Security](https://sonarcloud.io/api/project_badges/measure?project=my-cloud_ruthenium&metric=security_rating)](https://sonarcloud.io/component_measures?id=my-cloud_ruthenium&metric=security_rating)
[![Reliability](https://sonarcloud.io/api/project_badges/measure?project=my-cloud_ruthenium&metric=reliability_rating)](https://sonarcloud.io/component_measures?id=my-cloud_ruthenium&metric=reliability_rating)
[![Duplicated Lines](https://sonarcloud.io/api/project_badges/measure?project=my-cloud_ruthenium&metric=duplicated_lines_density)](https://sonarcloud.io/component_measures?id=my-cloud_ruthenium&metric=duplicated_lines_density)

## Description
The Ruthenium blockchain protocol.

This README contains essential information for a quick start. You will find a detailed description of the project in [the dedicated wiki page](https://github.com/my-cloud/ruthenium/wiki/Home). If you want to know what reasons led to create this blockchain, you can directly dive into the [Ruthenium whitepaper](https://github.com/my-cloud/ruthenium/wiki/Whitepaper). 

## Usage
There are two ways to use the Ruthenium blockchain. You can either use your own build from [sources](https://github.com/my-cloud/ruthenium/releases) (Option A) or use a docker image provided in the [repository packages](https://github.com/my-cloud/ruthenium/pkgs/container/ruthenium) (Option B).

### Prerequisites
* Option A: using sources
  * [Go][1] >= 1.17
* Option B: using docker image
  * [Docker][2]
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
Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are greatly appreciated.

If you have a suggestion that would make this better, please do not hesitate to [report a new bug](https://github.com/my-cloud/ruthenium/issues/new?assignees=&labels=bug&template=bug_report.md&title=) or [request a new feature](https://github.com/my-cloud/ruthenium/issues/new?assignees=&labels=enhancement&template=feature_request.md&title=). Of course, you are welcome to fork the repository and create a pull request. In any case, please let's take a look at the [contributing](https://github.com/my-cloud/ruthenium/blob/dev/CONTRIBUTING.md) file.

[⭐](https://github.com/login?return_to=%2Fmy-cloud%2Fruthenium) Don't forget to give the project a [star](https://github.com/login?return_to=%2Fmy-cloud%2Fruthenium)! Thanks again!

## Contact
Founder: [Jérémy Pansier](https://github.com/JeremyPansier) - jpansier@my-cloud.me

Project Link: https://github.com/my-cloud/ruthenium

## Authors and Acknowledgments
Special thanks to [Gwenall Pansier](https://github.com/Gwenall) who contributed since the early developments.

For a tutorial to create a first blockchain in go, thanks to [Yuko Sakai][3] & [Jun Sakai][4].

See also the list of [contributors](https://github.com/my-cloud/ruthenium/graphs/contributors) who participated in this project.

## License
![img.png](doc/img.png)

http://unlicense.org/

## Project status
The main principles have been implemented.  
Now it needs a lot of refactoring and tests to improve maintainability and
reliability.

[1]: https://go.dev/dl/ "Go website"
[2]: https://www.docker.com/ "Docker website"
[3]: https://www.udemy.com/user/myeigoworld/ "Yuko Sakai LinkedIn profile"
[4]: https://udemy.com/user/jun-sakai/ "Jun Sakai LinkedIn profile"

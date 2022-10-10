# Welcome to Ruthenium contributing guide
Thank you for investing your time in contributing to our project!

Please read our [Code of Conduct](https://github.com/my-cloud/ruthenium/blob/dev/CODE_OF_CONDUCT.md) to keep our community approachable and respectable.

In this guide you will get an overview of the contribution workflow from opening an issue, creating a PR, reviewing, merging the PR and releasing.

## New contributor guide
To get an overview of the project, read the [README](https://github.com/my-cloud/ruthenium#readme).



## Create a new issue
If you spot a problem with the docs, [search if an issue already exists](https://docs.github.com/en/github/searching-for-information-on-github/searching-on-github/searching-issues-and-pull-requests#search-by-the-title-body-or-comments). If a related issue doesn't exist, you can open a new issue using a relevant [issue form](https://github.com/my-cloud/ruthenium/issues/new/choose).

## Solve an issue
Scan through our [existing issues](https://github.com/my-cloud/ruthenium/issues) to find one that interests you. You can narrow down the search using [labels](https://github.com/my-cloud/ruthenium/labels) as filters. You can also take a look at the [open milestones](https://github.com/my-cloud/ruthenium/milestones) to have an idea of the issues priorities. To follow the progress of issues, let's take a look the [projects](https://github.com/my-cloud/ruthenium/projects?query=is%3Aopen)

As a general rule, we don’t assign issues to anyone. If you find an issue to work on, you are welcome to open a PR with changes to solve it.

### Prerequisites
* [Go](https://go.dev/dl/) >= 1.17
* [Git](https://git-scm.com/)

### Make changes
1. [Fork the project](https://github.com/my-cloud/ruthenium/fork).
1. Checkout a new branch (`git checkout -b feature/amazing-feature`).
1. Implement your solution following [our code conventions](#Go)
1. Commit and push your changes (`git add .; commit -m 'feat(blockchain): add some amazing feature'; git push origin feature/amazing-feature`).
1. [Create a pull request](https://github.com/my-cloud/ruthenium/compare) ([PR](https://docs.github.com/en/pull-requests)) targeting the `dev` branch.

## Review
🛡 Restricted to write access members.
1. Check if the changes follows [our conventions](#Go).
1. Manually [test](https://github.com/my-cloud/ruthenium/wiki/Usage) the changes. 
1. Submit your [review](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/reviewing-changes-in-pull-requests/reviewing-proposed-changes-in-a-pull-request).

## Merge
🛡 Restricted to write access members.
1. Check if the [targeted branch](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/changing-the-base-branch-of-a-pull-request) is the good one (most of the time, it should be `dev`).
1. Check if the squash commit message respects [our conventions](#Git) and finishes with the PR number (*ie* `fix(ui): message (#1)`).

## Release
🛡 Restricted to code owners.
* Set the `main` branch as the default branch.
* Merge the `dev` branch **LOCALLY** using fast-forward merge into `main`.
* Set the `dev` branch as the default branch.

The tag and release are automatically created using [semantic-release](https://github.com/go-semantic-release/action).  
The package is automatically pushed.  
The versioning follows the [semantic versioning convention](https://semver.org/).

## Conventions
### Git
The commit messages on the `dev` and `main` branches must follow the [Angular commit message format](https://github.com/angular/angular/blob/main/CONTRIBUTING.md#-commit-message-format)

### Go
We try follow the [Golang clean code conventions](https://github.com/Pungyeon/clean-go-article).
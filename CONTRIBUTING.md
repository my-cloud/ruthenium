# Resolving issues
1. Create an issue if it doesn't exist.
1. Fork the project.
1. Checkout a new branch.
1. Commit your changes.
1. Push to the branch.
1. Open a pull request.

# Following conventions
## Git
The commit messages on the `dev` and `main` branches must follow the [Angular commit message format](https://github.com/angular/angular/blob/main/CONTRIBUTING.md#-commit-message-format)

## Code
[Golang clean code conventins](https://github.com/Pungyeon/clean-go-article).

## Release
* Create a branch from `dev` named with the next version number check [how the release number will be generated](https://github.com/semantic-release/semantic-release/blob/master/docs/recipes/release-workflow/maintenance-releases.md#publishing-maintenance-releases).
* Open a pull request in the [ruthenium Github repository](https://github.com/my-cloud/ruthenium) from this branch to master
* Merge this branch without squashing commits into `main`.
The tag and release are automatically created using [semantic-release](https://github.com/go-semantic-release/action).  
The package is automatically pushed.  
The versioning follows the [semantic versioning convention](https://semver.org/).

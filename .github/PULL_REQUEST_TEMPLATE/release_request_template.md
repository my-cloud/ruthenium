---
name: Release pull request checklist
about: Open a pull request to the `main` to trigger a release
title: ''
labels: ''
assignees: ''
---

# [Release](https://github.com/my-cloud/ruthenium/blob/83-improve-contributing-documentation/CONTRIBUTING.md#release) checklist:
- [ ] The code follows the [our conventions](https://github.com/my-cloud/ruthenium/blob/83-improve-contributing-documentation/CONTRIBUTING.md#go).
- [ ] Important principle changes have been documented in the [wiki](https://github.com/my-cloud/ruthenium/wiki).
- [ ] The new code is unitary tested.
- [ ] The release branch is named with the [next version number](https://github.com/semantic-release/semantic-release/blob/master/docs/recipes/release-workflow/maintenance-releases.md#publishing-maintenance-releases).
- [ ] The target branch is `main`.
- [ ] ⚠ The commit will **NOT** be squashed.
- [ ] The branch will be deleted after being merged.
# Contributing Guidelines

Thank you for your interest in contributing to Klaytn Contracts. As an open source project, Klaytn Contracts is always open to the developer community and we welcome your contribution to help more developer onboarding resources for the Klaytn developer community. Please read the guideline below and follow it in all interactions with the project.

## How to Contribute on Klaytn Docs

1. Read this [contributing document](./CONTRIBUTING.md).
2. Sign [Contributor Licensing Agreement (CLA)](#contributor-license-agreement-cla).
3. Submit an issue with a proper [labelling](#usage-of-labels).
4. Please wait until the label changes to `contribution welcome` - otherwise, it is not ready to be worked on.
5. Only after the label changed to `contribution welcome`, you can start submitting the changes. To avoid any duplicate efforts, it is recommended to update the issue so that other contributors could see someone working on the issue.
6. Before making a Pull Request (PR), please make sure the suggested content changes are accurate and linked with the corresponding issue reported. After submitting the PR, wait for code review and approval. The reviewer may ask you for additional commits or changes.
7. All PRs should be made against the `master` branch, once the change has been approved, the PR is merged by the project moderator.
8. After merging the PR, the pull request will be closed. You can then delete the now obsolete branch.

## Types of Contribution
There are various ways to contribute and participate. Please read the guidelines below regarding the process of each type of contribution.

-   [Issues and Bugs](#issues-and-bugs)
-   [Feature Requests](#feature-requests)
-   [Code Contribution](#code-contribution)

### Issues and Bugs

If you find a bug or other issues in Klaytn, please [submit an issue](https://github.com/klaytn/klaytn-contracts/issues). Before submitting an issue, please invest some extra time to figure out that:

- The issue is not a duplicate issue.
- The issue has not been fixed in the latest release of Klaytn Contracts.
Please do not use the issue tracker for personal support requests. Use developer@klaytn.foundation for the personal support requests.

When you report a bug, please make sure that your report has the following information.
- Steps to reproduce the issue.
- A clear and complete description of the issue.
- Code and/or screen captures are highly recommended.

After confirming your report meets the above criteria, [submit the issue](https://github.com/klaytn/klaytn-contracts/issues). Please use [labels](#usage-of-labels) to categorize your issue.

### Feature Requests

You can also use the [issue tracker](https://github.com/klaytn/klaytn-contracts/issues) to request a new feature or enhancement. Note that any code contribution without an issue link will not be accepted. 

Please submit an issue explaining your proposal first so that the Klaytn developer community can fully understand and discuss the idea. Please use [labels](#usage-of-labels) for your feature request as well.

### Code Contribution 

Smart contracts manage value and are highly vulnerable to errors and attacks. Hence, please make sure to review the following guidelines. 

#### Style Guidelines

The design guidelines have quite a high abstraction level. These style guidelines are more concrete and easier to apply, and also more opinionated. We value clean code and consistency, and those are prerequisites for us to include new code in the repository. Before proposing a change, please read these guidelines and take some time to familiarize yourself with the style of the existing codebase.

#### Solidity code

In order to be consistent with all the other Solidity projects, we follow the
[official recommendations documented in the Solidity style guide](http://solidity.readthedocs.io/en/latest/style-guide.html) and [openzepplin](https://github.com/OpenZeppelin/openzeppelin-contracts/blob/master/GUIDELINES.md) style guide.

#### Tests

* Tests Must be Written Elegantly

    Tests are a good way to show how to use the library, and maintaining them is extremely necessary. Don't write long tests, write helper functions to make them be as short and concise as possible (they should take just a few lines each), and use good variable names.

* Tests Must not be Random

    Inputs for tests should not be generated randomly. Accounts used to create test contracts are an exception, those can be random. Also, the type and structure of outputs should be checked.

#### Creating a Pull Request(PRs)
As a contributor, you are expected to fork this repository, work on your own fork and then submit pull requests. The pull requests will be reviewed and eventually merged into the main repo. See ["Fork-a-Repo"](https://help.github.com/articles/fork-a-repo/) for how this works.

1) Make sure your fork is up to date with the main repository:

```
cd klaytn-contracts
git remote add upstream https://github.com/klaytn/klaytn-contracts.git
git fetch upstream
git pull --rebase upstream master
```
NOTE: The directory `klaytn-contracts` represents your fork's local copy.

2) Branch out from `master` into `fix/some-bug-#123`/`feature/add-timelock`. 

```
git checkout -b fix/some-bug-#123
```

3) Make your changes, add your files, commit, and push to your fork.  
```
git add SomeFile.js
git commit "Fixed bug #123"
git push origin fix/some-bug-#123
```

4) Run tests, linter, etc. This can be done by running local continuous integration and make sure it passes.

```bash
npm test
npm run lint:fix
npm run lint
```

5) Create a new pull request. If its a issue, make sure you link the issue in the [Related issues](https://github.com/klaytn/klaytn-contracts/blob/4b9efe2976d892e11f839bd3e05739226d4fd81f/.github/PULL_REQUEST_TEMPLATE.md)

*IMPORTANT* Read the PR template very carefully and make sure to follow all the instructions. These instructions
refer to some very important conditions that your PR must meet in order to be accepted, such as making sure that all tests pass, JS linting tests pass, Solidity linting tests pass, etc.

6) Maintainers will review your code and possibly ask for changes before your code is pulled in to the main repository. We'll check that all tests pass, review the coding style, and check for general code correctness. If everything is OK, we'll merge your pull request and your code will be part of Klaytn Contracts.


#### Usage of Labels

You can use the following labels:

Labels for initial issue categories:

- issue/bug: Issues with the code-level bugs.
- issue/enhancement: Issues for enhancement requests.
- issue/questions: Questions to klaytn contracts and other issues not related to bug or enhancement.

Status of open issues (will be tagged by the project moderators):

- (no label): The default status.
- open/need more information : The issue's creator needs to provide additional information to review.
- open/reviewing: The issue is under review.
- open/re-label needed: The label needs to be changed to confirmed as being a `bug` or future `enhancement`.
- open/approved: The issue is confirmed as being a `bug` to be fixed or `enhancement` to be developed.
- open/contribution welcome: The fix or enhancement is approved and you are invited to contribute to it.

Status of closed issues:

- closed/fixed: A fix for the issue was provided.
- closed/duplicate: The issue is also reported in a different issue and is being managed there.
- closed/invalid: The issue cannot be reproduced.
- closed/reject: The issue is rejected after review.
- closed/wontfix: This issue will not be worked on.


## Contributor License Agreement (CLA)

Keep in mind when you submit your pull request, you'll need to sign the CLA via [CLA-Assistant](https://github.com/klaytn/klaytn-contracts) for legal purposes. You will have to sign the CLA just one time, either as an individual or corporation.

You will be prompted to sign the agreement by CLA Assistant (bot) when you open a Pull Request for the first time.
# Contributing Guidelines

Thank you for your interest in contributing to Klaytn. As an open source project, Klaytn is always open to the developer community and we welcome your contribution. Please read the guideline below and follow it in all interactions with the project.

## How to Contribute

1. Read this [contributing document](./CONTRIBUTING.md).
2. Sign [Contributor Licensing Agreement (CLA)](#contributor-license-agreement-cla).
3. Submit an issue with proper [labeling](#usage-of-labels).
4. Please wait until the label changes to `contribution welcome` - otherwise, it is not ready to be worked on.
5. After the label changed to `contribution welcome`, you can start working on the implementation. To avoid any duplicate efforts, it is recommended to update the issue so that other developers see someone working on the issue.
6. Before making a PR, please make sure you fully tested the code. It is highly recommended to provide the test code as well. After submitting the PR, wait for code review and approval. The reviewer may ask you for additional commits or changes.
7. Once the change has been approved, the PR is merged by the project moderator.
8. After merging the PR, we close the pull request. You can then delete the now obsolete branch.

## Types of Contribution
There are various ways to contribute and participate. Please read the guidelines below regarding the process of each type of contribution.

-   [Issues and Bugs](#issues-and-bugs)
-   [Feature Requests](#feature-requests)
-   [Code Contribution](#code-contribution)

### Issues and Bugs

If you find a bug or other issues in Klaytn, please [submit an issue](https://github.com/klaytn/klaytn/issues). Before submitting an issue, please invest some extra time to figure out that:

- The issue is not a duplicate issue.
- The issue has not been fixed in the latest release of Klaytn.
Please do not use the issue tracker for personal support requests. Use developer@klaytn.com for the personal support requests.

When you report a bug, please make sure that your report has the following information.
- Steps to reproduce the issue.
- A clear and complete description of the issue.
- Code and/or screen captures are highly recommended.

After confirming your report meets the above criteria, [submit the issue](https://github.com/klaytn/klaytn/issues). Please use [labels](#usage-of-labels) to categorize your issue.

### Feature Requests

You can also use the [issue tracker](https://github.com/klaytn/klaytn/issues) to request a new feature or enhancement. Note that any code contribution without an issue link will not be accepted. Please submit an issue explaining your proposal first so that Klaytn community can fully understand and discuss the idea. Please use [labels](#usage-of-labels) for your feature request as well.

#### Usage of Labels

You can use the following labels:

Labels for initial issue categories:

- issue/bug: Issues with the code-level bugs.
- issue/documentation: Issues with the documentation.
- issue/enhancement: Issues for enhancement requests.

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

### Code Contribution

Please follow the coding style and quality requirements to satisfy the product standards. You must follow the coding style as best as you can when submitting code. Take note of naming conventions, separation of concerns, and formatting rules.

The go implementation of Klaytn uses [godoc](https://godoc.org/golang.org/x/tools/cmd/godoc)
to document its source code. For the guideline of official Go language, please
refer to the following websites:
- https://golang.org/doc/effective_go.html#commentary
- https://blog.golang.org/godoc-documenting-go-code

## Contributor License Agreement (CLA)

Keep in mind when you submit your pull request, you'll need to sign the CLA via [CLA-Assistant](https://cla-assistant.io/klaytn/klaytn) for legal purposes. You will have to sign the CLA just one time, either as an individual or corporation.

You will be prompted to sign the agreement by CLA Assistant (bot) when you open a Pull Request for the first time.

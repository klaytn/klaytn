# tags for TODO & NOTE

This document is to collect the existing tags for TODO & NOTE and their descriptions. A Tag here means a word or words used to categorize related comments. (e.g., `TODO-Klaytn-TAG`, `TODO-Klaytn-ServiceChain` ) Before leaving a comment using `TODO` or `NOTE`, please check this document to use proper tag. If a proper tag exists, use it. However, if there is no proper tag, use a new tag and also add it to this document.



## Why use tags?

As described above, its primary purpose is to **categorize related comments and issues by tag**. If you want to search for something to do in ServiceChain category, you can just `grep TODO-Klaytn-ServiceChain` to get the list of TODOs in ServiceChain category. Also, by putting `Klaytn` between prefix and tag, **we can distinguish TODOs from geth and Klaytn**. Below are examples of comments using tags.

```
// TODO-Klaytn-TAG
// e.g., TODO-Klaytn-ServiceChain Need to add an option, not to write receipts.
// NOTE-Klaytn-TAG
// e.g., NOTE-Klaytn-RemoveLater Below Prove is only used in tests, not in core codes.
```



## Prefixes - TODO, NOTE

- Above prefixes should be used in `PREFIX-Klaytn-TAG` format.
- **TODO**
  - When you leave a comment to describe **something needed to be done**.
- **NOTE**
  - When you leave a comment to explain **something important to be noticed** by Klaytn developers.
  - `NOTE` is not highlighted by GoLand IDE by default, however, [you can change the setting to highlight NOTE keyword.](https://www.jetbrains.com/help/idea/using-todo.html)



## List of Tags

Please note that the tag name always

- **starts with a capital letter**
- **uses camel-case if it has two words or more**.

| TagName                                        | Description                                               |
| ---------------------------------------------- | --------------------------------------------------------- |
| Issue[IssueNumber], e.g., TODO-Klaytn-Issue833 | If a corresponding issue exists.                          |
| Docs                                           | Documentation or copyright related comments.              |
| BN                                             | Related to BootNode.                                      |
| FailedTest                                     | Tests currently fail on Klaytn. (therefore commented out) |
| RemoveLater                                    | Codes which can be removed later due to some reasons.     |
| FeePayer                                       | Related to FeePayer feature.                              |
| Gas                                            | Related to gas calculation and its policy.                |
| MultiSig                                       | Related to MultiSig feature.                              |
| TestNet                                        | Related to TestNet.                                       |
| NodeDiscovery                                  | Related to NodeDiscovery.                                 |
| Accounts                                       | Related to Klaytn Accounts.                               |
| Downloader                                     | Related to Downloader.                                    |
| NodeCmd                                        | Related to cmd/utils/nodecmd.                             |
| Storage                                        | Related to storage.                                       |
| Refactoring                                    | Related to refactoring.                                   |
| gRPC                                           | Related to gRPC.                                          |
| HF                                             | Related to HardFork.                                      |
| StateDB                                        | Related to StateDB and stateObject.                       |
| DataArchiving                                  | Related to Data Archiving feature.                        |

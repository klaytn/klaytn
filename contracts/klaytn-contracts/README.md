# Klaytn Contracts

[![Docs](https://img.shields.io/badge/docs-%F0%9F%93%84-blue)](https://docs.klaytn.foundation/)
[![NPM Package](https://badge.fury.io/js/@klaytn%2Fcontracts.svg)](https://www.npmjs.com/package/@klaytn/contracts)
[![Gitter](https://badges.gitter.im/klaytn/klaytn-contracts.svg)](https://gitter.im/klaytn/klaytn-contracts?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

**A library for secure smart contract development.** Build on a solid foundation of community-vetted code.
It is a fork of openzepplin contracts. In addition to that, this repository contains Klaytn's token standards such as KIP-7, KIP-17, and KIP-37 compatible with ERC-20, ERC-721, and ERC-1155 respectively.

Please refer to [this](https://kips.klaytn.com/token) link for the mapping of Ethereum to Klaytn token standards.

## Overview

### Installation

```console
$ npm install @klaytn/contracts
```

An alternative to npm is to use the GitHub repository `klaytn/klaytn-contracts` to retrieve the contracts. When doing this, make sure to specify the tag for a release such as `v1.0.0`, instead of using the `master` branch.

### Usage

Once installed, you can use the contracts in the library by importing them:

```solidity
pragma solidity ^0.8.0;

import "@klaytn/contracts/contracts/KIP/token/KIP17/KIP17.sol";

contract MyCollectible is KIP17 {
    constructor() KIP17("MyCollectible", "MCO") {
    }
}
```

To keep your system secure, you should **always** use the installed code as-is, and neither copy-paste it from online sources, nor modify it yourself. The library is designed so that only the contracts and functions you use are deployed, so you don't need to worry about it needlessly increasing gas costs.

## Contribute

In line with our commitment to decentralization, all Klaytn codebase and its documentations are completely open source. Klaytn always welcomes your contribution. Anyone can view, edit, fix its contents and make suggestions. You can either create a pull request on GitHub or use GitBook. Make sure to sign our [Contributor License Agreement (CLA)](https://cla-assistant.io/klaytn/klaytn-contracts) first and there are also a few guidelines our contributors would check out before contributing:

- [Contribution Guide](./CONTRIBUTING.md)
- [License](./LICENSE)
- [Code of Conducts](./code-of-conduct.md)

## Need Help? <a href="#need-help" id="need-help"></a>

If you have any questions, please visit our [Gitter channel](https://gitter.im/klaytn/klaytn-contracts?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge), [Klaytn Developers Forum](https://forum.klaytn.foundation/) and [Discord channel](https://discord.gg/mWsHFqN5Zf).

## License

 Contracts is released under the [MIT License](LICENSE).

## Acknowledgments 

Thanks for Openzepplin Team for providing the  contracts.

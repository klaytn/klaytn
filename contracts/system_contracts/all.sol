// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

// This file import all contracts so that single `abigen` can generate all bindings at once,
// Because otherwise, separate `abigen` per contract may cause symbol conflict
// when one .sol file is imported from multiple .sol files.

import "./registry/Registry.sol";
import "./registry/RegistryMock.sol";
import "./kip113/KIP113.sol";
import "./kip113/KIP113Mock.sol";
import "./lib/ERC1967Proxy.sol";

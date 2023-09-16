// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

// This file import all contracts so that single `abigen` can generate all bindings at once.
// In addition, separate `abigen` for each contract may cause symbol conflict
// when one .sol file is imported from multiple .sol files.

import "kip103/TreasuryRebalance.sol";
import "kip103/TreasuryRebalanceMock.sol";

## Klaytn system contracts
Package system deals with system contracts in Klaytn.

The system contracts are smart contracts that controls the blockchain consensus.
- AddressBook: the list of consensus nodes. It stores their nodeId, staking contracts and reward
  addresses.
- GovParam: the governance parameter storage, meant to be modified by on-chain governance vote.
  It overrides the existing header-based governance (i.e. block header votes) and it can be
  optionally configured after the Kore hardfork.
- TreasuryRebalance: records the rebalance of treasury funds. Configured by the [KIP103](https://kips.klaytn.foundation/KIPs/kip-103) hardfork.

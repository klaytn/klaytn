// Copyright 2019 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

/*
Package sc implements an auxiliary blockchain called Service Chain.

Service Chains in Klaytn are auxiliary blockchains independent from the Klaytn main chain.
They mostly act like normal Klaytn blockchains but has additional features to connect them to another Klaytn network.
They can be used as separate public/private blockchains or child chains of a Klaytn chain (or another Service Chain).
The followings describe main features of Service chain.
  - Anchoring block data of Service Chain
  - Value Transfer (KLAY, KCT)
  - Various bridge contract configurations
  - Support high availability


Service Chain provides the inter-connectivity to another Klaytn chain through two bridges, MainBridge and SubBridge.
Each bridge has a bridge contract on different blockchain networks and pairs up with another bridge to interact with.
They are directly connected on the network layer and communicate with each other through p2p messages enabling inter-chain operations.

MainBridge is configured on the node of a parent chain and SubBridge is configured on the node of a child chain.
Both of a Klaytn chain and a Service Chain can be a parent chain, but only a Service Chain can be a child chain.
The block data of a child chain can be anchored to the bridge contract of MainBridge with the chain data anchoring transaction.

Unlike the block data anchoring, user data transfer is bi-directional.
For example, users can transfer KLAY of Klaytn main chain to an address of a Service Chain or vice versa.
This kind of inter-chain operation requires read/write ability on both chains but does not use MainBridge functions in the process.
Instead of the MainBridge, the SubBridge in the child chain directly calls read/write operations to the parent chain node through RPC (In the basic configuration, the parent chain node is the same with the MainBridge enabled node).
Of course, the accounts of both chains should be registered on the SubBridge to generate transactions.
Following is the process of the KLAY transfer from Klaytn main chain to a Service Chain.
1. A user executes the inter-chain operation by sending a transaction with KLAY to the bridge contract of Klaytn main chain.
2. The bridge contract keeps KLAY on its account and creates an event for the inter-chain request.
3. The SubBridge subscribes the event log on the main chain node through RPC.
4. The SubBridge generates a transaction on the child chain node to the bridge contract of the SubBridge.
5. Finally, The bridge contract mints (or uses its KLAY) and sends KLAY to the target address.


Source Files

Functions and variables related to Service Chain are defined in the files listed below.
  - api_bridge.go : provides APIs for MainBridge or SubBridge.
  - bridge_accounts.go : generates inter-chain transactions between a parent chain and a child chain.
  - bridge_addr_journal.go : provides a journal mechanism for bridge addresses to provide the persistence service.
  - bridge_manager.go : handles the bridge information and manages the bridge operations.
  - bridgepeer.go : implements data structures of p2p peers used for Service Chain bridges.
  - config.go : provides configurations of Service Chain nodes.
  - gen_config.go : provides marshalling and unmarshalling functions of the Service Chain configuration.
  - local_backend.go : provides read/write operations for the child chain block data.
  - main_bridge_handler.go : implements a p2p message handler of MainBridge.
  - main_event_handler.go : implements a event handler of MainBridge.
  - mainbridge.go : implements MainBridge of the parent chain node.
  - metrics.go : contains metrics used for sc package.
  - protocol.go : defines protocols of Service Chain.
  - remote_backend.go : provides read/write RPC calls for the parent chain block data.
  - sub_bridge_handler.go : implements a p2p message handler of SubBridge.
  - sub_event_handler.go : implements a event handler of SubBridge.
  - subbridge.go : implements SubBridge of the child chain node.
  - vt_recovery.go : provides recovery from the service failure for inter-chain value transfer.
*/
package sc

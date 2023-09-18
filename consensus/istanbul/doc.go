// Copyright 2018 The klaytn Authors
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
Package istanbul is a BFT based consensus engine which implements consensus/Engine interface.
Istanbul engine was inspired by Clique POA, Hyperledger's SBFT, Tendermint, HydraChain and NCCU BFT.

Istanbul engine is using 3-phase consensus and it can tolerate F faulty nodes where N = 3F + 1

In Klaytn, it is being used as the main consensus engine after modification for supports of Committee, Reward and Governance.
Package istanbul has three sub-packages, core, backend, and validator. Please refer to each package's doc.go for more information.

Source Files

Various interfaces, constants and utility functions for Istanbul consensus engine
 - `backend.go`: Defines Backend interface which provides application specific functions for Istanbul core
 - `config.go`: Provides default configuration for Istanbul engine
 - `errors.go`: Defines three errors used in Istanbul engine
 - `events.go`: Defines events which are used for Istanbul engine communication
 - `types.go`: Defines message structs such as Proposal, Request, View, Preprepare, Subject and ConsensusMsg
 - `utils.go`: Provides three utility functions: RLPHash, GetSignatureAddress and CheckValidatorSignature
 - `validator.go`: Defines Validator, ValidatorSet interfaces and Validators, ProposalSelector types
*/
package istanbul

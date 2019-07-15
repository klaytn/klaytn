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

package contracts

//go:generate abigen --sol ./bridge/Bridge.sol --pkg bridge --out ./bridge/Bridge.go
//go:generate abigen --sol ./sc_erc721/sc_nft.sol --pkg scnft --out ./sc_erc721/sc_nft.go
//go:generate abigen --sol ./sc_erc20/sc_token.sol --pkg sctoken --out ./sc_erc20/sc_token.go
//go:generate abigen --sol ./cypress/credit.sol --pkg cypress --out ./cypress/credit.go

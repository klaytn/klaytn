// Copyright 2020 The klaytn Authors
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
Package chaindatafetcher implements blockchain data load to KAS-specific database, or kafka.
Source Files
  - api.go                   : includes chaindatafetcher-related APIs
  - chaindata_fetcher.go     : implements chaindatafetcher main operations
  - config.go                : includes chaindatafetcher configurations
  - metrics.go               : includes chaindatafetcher metrics
  - repository.go            : implements repository interface
*/

package chaindatafetcher

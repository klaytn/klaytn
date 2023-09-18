// Copyright 2018 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
// This file is part of go-ethereum.
//
// Copyright 2017 AMIS Technologies
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//

/*
service package provides various templates to build a docker-compose.yml

Source Files

Each file contains following contents
 - constellation.go : Deprecated. This is not being used in Klaytn
 - eth_stats.go : Defines `KlayStats` and provides a yaml template for a KlayStats service
 - grafana.go : Defines `GrafanaService` and provides a yaml template for a Grafana service
 - prometheus.go : Defines `PrometheusService` and provides a yaml template for a Prometheus service
 - txgen.go : Defines `TxGenService` and provides a yaml template for a txgen service
 - validator.go : Defines `Validator` and provides a yaml template for a validator configuration

*/
package service

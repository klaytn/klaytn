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

package discover

func NewDiscovery(cfg *Config) (Discovery, error) {
	return newTable(cfg)
}

func ParseNodeType(nt string) NodeType {
	switch nt {
	case "cn":
		return NodeTypeCN
	case "pn":
		return NodeTypePN
	case "en":
		return NodeTypeEN
	case "bn":
		return NodeTypeBN
	default:
		return NodeTypeUnknown
	}
}

// StringNodeType converts NodeType to string
func StringNodeType(nType NodeType) string { // TODO-Klaytn-Node Consolidate p2p.NodeType and common.ConnType
	switch nType {
	case NodeTypeCN:
		return "cn"
	case NodeTypePN:
		return "pn"
	case NodeTypeEN:
		return "en"
	case NodeTypeBN:
		return "bn"
	default:
		return "unknown"
	}
}

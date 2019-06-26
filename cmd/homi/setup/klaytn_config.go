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

package setup

import (
	"bytes"
	"fmt"
	"text/template"
)

type KlaytnConfig struct {
	NetworkId int
	RPCPort   int
	WSPort    int
	P2PPort   int
	DataDir   string
	LogDir    string
	RunDir    string
	NodeType  string
}

func New(rpcPort int, wsPort int, p2pPort int, dataDir string, logDir string, runDir string, nodeType string) *KlaytnConfig {
	kConf := &KlaytnConfig{
		RPCPort:  rpcPort,
		WSPort:   wsPort,
		P2PPort:  p2pPort,
		DataDir:  dataDir,
		LogDir:   logDir,
		RunDir:   runDir,
		NodeType: nodeType,
	}
	return kConf
}

func (k KlaytnConfig) String() string {
	tmpl, err := template.New("KlaytnConfig").Parse(kTemplate)
	if err != nil {
		fmt.Printf("Failed to parse template, %v", err)
		return ""
	}

	res := new(bytes.Buffer)
	err = tmpl.Execute(res, k)
	if err != nil {
		fmt.Printf("Failed to render template, %v", err)
		return ""
	}

	return res.String()
}

var kTemplate = `# Configuration file for the klay service.

NETWORK_ID={{ .NetworkId }}

RPC_PORT={{ .RPCPort }}
WS_PORT={{ .WSPort }}
PORT={{ .P2PPort }}

DATA_DIR={{ .DataDir }}
LOG_DIR={{ .LogDir }}
RUN_DIR={{ .RunDir }}

# NODE_TYPE [CN, PN, RN]
NODE_TYPE={{ .NodeType }}
`

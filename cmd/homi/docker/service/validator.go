// Copyright 2018 The klaytn Authors
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

package service

import (
	"bytes"
	"fmt"
	"text/template"
)

type Validator struct {
	Identity       int
	Genesis        string
	SCGenesis      string
	Address        string
	NodeKey        string
	StaticNodes    string
	BridgeNodes    string
	Port           int
	RPCPort        int
	PrometheusPort int
	IP             string
	EthStats       string
	Name           string
	DockerImageId  string
	UseFastHttp    bool
	NetworkId      int
	ParentChainId  int
	NodeType       string
	AddPrivKey     bool
}

func NewValidator(identity int, genesis, scGenesis string, nodeAddress string, nodeKey string, staticNodes, bridgeNodes string, port int, rpcPort int,
	prometheusPort int, ethStats string, ip string, dockerImageId string, useFastHttp bool, networkId, parentChainId int,
	namePrefix string, nodeType string, addPrivKey bool) *Validator {
	return &Validator{
		Identity:       identity,
		Genesis:        genesis,
		SCGenesis:      scGenesis,
		Address:        nodeAddress,
		NodeKey:        nodeKey,
		Port:           port,
		RPCPort:        rpcPort,
		PrometheusPort: prometheusPort,
		EthStats:       ethStats,
		IP:             ip,
		Name:           fmt.Sprintf("%s-%v", namePrefix, identity),
		DockerImageId:  dockerImageId,
		UseFastHttp:    useFastHttp,
		NetworkId:      networkId,
		ParentChainId:  parentChainId,
		NodeType:       nodeType,
		AddPrivKey:     addPrivKey,
		StaticNodes:    staticNodes,
		BridgeNodes:    bridgeNodes,
	}
}

func (v Validator) String() string {
	tmpl, err := template.New("validator").Parse(validatorTemplate)
	if err != nil {
		fmt.Printf("Failed to parse template, %v", err)
		return ""
	}

	result := new(bytes.Buffer)
	err = tmpl.Execute(result, v)
	if err != nil {
		fmt.Printf("Failed to render template, %v", err)
		return ""
	}
	return result.String()
}

var validatorTemplate = `{{ .Name }}:
    hostname: {{ .Name }}
    image: {{ .DockerImageId }}
    ports:
      - '{{ .Port }}:32323'
      - '{{ .RPCPort }}:8551'
      - '{{ .PrometheusPort }}:61001'
{{- if eq .Name "EN-0" }}
      - '50505:50505'
{{- else if eq .Name "SCN-0" }}
      - '50506:50506'
{{- end }}
    entrypoint:
      - /bin/sh
      - -c
      - |
        mkdir -p /klaytn
{{- if eq .NodeType "scn" }}
        echo '{{ .SCGenesis }}' > /klaytn/genesis.json
{{- else if eq .NodeType "spn" }}
        echo '{{ .SCGenesis }}' > /klaytn/genesis.json
{{- else if eq .NodeType "sen" }}
        echo '{{ .SCGenesis }}' > /klaytn/genesis.json
{{- else }}
        echo '{{ .Genesis }}' > /klaytn/genesis.json
{{- end }}
        echo '{{ .StaticNodes }}' > /klaytn/static-nodes.json
{{- if .BridgeNodes }}
        echo '{{ .BridgeNodes }}' > /klaytn/main-bridges.json
{{- end }}
        k{{ .NodeType }} --datadir "/klaytn" init "/klaytn/genesis.json"

{{- if .AddPrivKey}}
        echo '{"address":"75a59b94889a05c03c66c3c84e9d2f8308ca4abd","crypto":{"cipher":"aes-128-ctr","ciphertext":"347fef8ab9aaf9d41b6114dfc0d9fd6ecab9d660fa86f687dc7aa1e094b76184","cipherparams":{"iv":"5070268dfc64ced716cf407bee943def"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"2cf44cb912515c5de2aacf4d133fd97a56efece5a8fc691381296300b42fb6c8"},"mac":"ea20c0019321f2a91f1ca51bc99d58c8f7cf1f37cff5cc47ae17ad747c060046"},"id":"a122a8da-a787-4d3d-a627-02034553e674","version":1}' > /klaytn/keystore/mykey
        echo "SuperSecret1231" > /klaytn/password.txt
{{- end}}
        echo "# docker-compose" >> /klaytn-docker-pkg/conf/k{{ .NodeType }}d.conf
        echo 'NETWORK=""' >> /klaytn-docker-pkg/conf/k{{ .NodeType }}d.conf
        echo 'DATA_DIR="/klaytn"' >> /klaytn-docker-pkg/conf/k{{ .NodeType }}d.conf
        echo 'LOG_DIR="$$DATA_DIR/log"' >> /klaytn-docker-pkg/conf/k{{ .NodeType }}d.conf
        echo 'RPC_ENABLE=1' >> /klaytn-docker-pkg/conf/k{{ .NodeType }}d.conf
        echo 'RPC_API="db,eth,klay,net,web3,miner,personal,txpool,debug,admin,istanbul,mainbridge,subbridge"' >> /klaytn-docker-pkg/conf/k{{ .NodeType }}d.conf
        echo 'NETWORK_ID="{{ .NetworkId }}"' >> /klaytn-docker-pkg/conf/k{{ .NodeType }}d.conf
        echo 'NO_DISCOVER=1' >> /klaytn-docker-pkg/conf/k{{ .NodeType }}d.conf
        echo 'ADDITIONAL="$$ADDITIONAL --identity \"{{ .Name }}\""' >> /klaytn-docker-pkg/conf/k{{ .NodeType }}d.conf
        echo 'ADDITIONAL="$$ADDITIONAL --nodekeyhex {{ .NodeKey }}"' >> /klaytn-docker-pkg/conf/k{{ .NodeType }}d.conf
{{- if .AddPrivKey}}
        echo 'ADDITIONAL="$$ADDITIONAL --unlock \"75a59b94889a05c03c66c3c84e9d2f8308ca4abd\""' >> /klaytn-docker-pkg/conf/k{{ .NodeType }}d.conf
        echo 'ADDITIONAL="$$ADDITIONAL --password "/klaytn/password.txt"' >> /klaytn-docker-pkg/conf/k{{ .NodeType }}d.conf
{{- end}}
{{- if eq .NodeType "scn" }}
        echo 'PORT=32323' >> /klaytn-docker-pkg/conf/k{{ .NodeType }}d.conf
{{- end}}
{{- if .ParentChainId }}
        echo 'ADDITIONAL="$$ADDITIONAL --parentchainid {{ .ParentChainId }}"' >> /klaytn-docker-pkg/conf/k{{ .NodeType }}d.conf
        echo 'SC_SUB_BRIDGE=1' >> /klaytn-docker-pkg/conf/k{{ .NodeType }}d.conf
        echo 'SC_SUB_BRIDGE_PORT=50506' >> /klaytn-docker-pkg/conf/k{{ .NodeType }}d.conf
{{- end}}
{{- if .UseFastHttp}}
        echo 'SERVER_TYPE=fasthttp' >> /klaytn-docker-pkg/conf/k{{ .NodeType }}d.conf
{{- end}}
{{- if eq .NodeType "cn" }}
        echo 'REWARDBASE={{ .Address }}' >> /klaytn-docker-pkg/conf/k{{ .NodeType }}d.conf
{{- else if eq .NodeType "pn" }}
        echo 'ADDITIONAL="$$ADDITIONAL --txpool.nolocals"' >> /klaytn-docker-pkg/conf/k{{ .NodeType }}d.conf
{{- else if eq .Name "EN-0" }}
        echo 'SC_MAIN_BRIDGE=1' >> /klaytn-docker-pkg/conf/k{{ .NodeType }}d.conf
        echo 'SC_MAIN_BRIDGE_PORT=50505' >> /klaytn-docker-pkg/conf/k{{ .NodeType }}d.conf
        echo 'SC_MAIN_BRIDGE_INDEXING=1' >> /klaytn-docker-pkg/conf/k{{ .NodeType }}d.conf
{{- end }}
        echo 'ADDITIONAL="$$ADDITIONAL --debug --metrics --prometheus"' >> /klaytn-docker-pkg/conf/k{{ .NodeType }}d.conf
        /klaytn-docker-pkg/bin/k{{ .NodeType }}d start
{{- if eq .NodeType "cn"}}
        sleep 1
        ken attach --exec "personal.importRawKey('{{ .NodeKey }}', '')" http://localhost:{{ .RPCPort }}
        ken attach --exec "personal.unlockAccount('{{ .Address }}', '', 999999999)" http://localhost:{{ .RPCPort }}
{{- end }}
        tail -F klaytn/log/k{{ .NodeType }}d.out
    networks:
      app_net:
        ipv4_address: {{ .IP }}
    restart: "no"`

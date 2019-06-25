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
	Address        string
	NodeKey        string
	StaticNodes    string
	Port           int
	RPCPort        int
	PrometheusPort int
	IP             string
	EthStats       string
	Name           string
	DockerImageId  string
	UseFastHttp    bool
	NetworkId      int
	NodeType       string
	AddPrivKey     bool
}

func NewValidator(identity int, genesis string, nodeAddress string, nodeKey string, staticNodes string, port int, rpcPort int,
	prometheusPort int, ethStats string, ip string, dockerImageId string, useFastHttp bool, networkId int,
	namePrefix string, nodeType string, addPrivKey bool) *Validator {
	return &Validator{
		Identity:       identity,
		Genesis:        genesis,
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
		NodeType:       nodeType,
		AddPrivKey:     addPrivKey,
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
    entrypoint:
      - /bin/sh
      - -c
      - |
        mkdir -p /klaytn
        echo '{{ .Genesis }}' > /klaytn/genesis.json
        echo '{{ .StaticNodes }}' > /klaytn/static-nodes.json
        k{{ .NodeType }} --datadir "/klaytn" init "/klaytn/genesis.json"

{{- if .AddPrivKey}}
        echo '{"address":"75a59b94889a05c03c66c3c84e9d2f8308ca4abd","crypto":{"cipher":"aes-128-ctr","ciphertext":"347fef8ab9aaf9d41b6114dfc0d9fd6ecab9d660fa86f687dc7aa1e094b76184","cipherparams":{"iv":"5070268dfc64ced716cf407bee943def"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"2cf44cb912515c5de2aacf4d133fd97a56efece5a8fc691381296300b42fb6c8"},"mac":"ea20c0019321f2a91f1ca51bc99d58c8f7cf1f37cff5cc47ae17ad747c060046"},"id":"a122a8da-a787-4d3d-a627-02034553e674","version":1}' > /klaytn/keystore/mykey
        echo "SuperSecret1231" > /klaytn/password.txt
{{- end}}
        k{{ .NodeType }} \
        --identity "{{ .Name }}" \
        --rpc \
        --rpcaddr "0.0.0.0" \
        --rpcport "8551" \
        --rpccorsdomain "*" \
        --datadir "/klaytn" \
        --port "32323" \
        --rpcapi "db,klay,net,web3,miner,personal,txpool,debug,admin,istanbul" \
        --networkid "{{ .NetworkId }}" \
        --nat "any" \
        --nodekeyhex "{{ .NodeKey }}" \
        --nodiscover \
        --debug \
        --metrics \
        --prometheus \
        --syncmode "full" \
{{- if .AddPrivKey}}
        --unlock "75a59b94889a05c03c66c3c84e9d2f8308ca4abd" \
        --password "/klaytn/password.txt" \
{{- end}}
{{- if .UseFastHttp}}
        --srvtype fasthttp \
{{- end}}
{{- if eq .NodeType "cn" }}
        --rewardbase {{ .Address }}
{{- else if eq .NodeType "pn" }}
        --txpool.nolocals
{{- else }}
{{- end}}

    networks:
      app_net:
        ipv4_address: {{ .IP }}
    restart: "no"`

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

package compose

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/klaytn/klaytn/cmd/homi/docker/service"
)

type Homi struct {
	IPPrefix          string
	EthStats          *service.KlayStats
	Services          []*service.Validator
	PrometheusService *service.PrometheusService
	GrafanaService    *service.GrafanaService
	UseGrafana        bool
	Proxies           []*service.Validator
	UseTxGen          bool
	TxGenService      *service.TxGenService
	TxGenOpt          service.TxGenOption
}

func New(ipPrefix string, number int, secret string, addresses, nodeKeys []string,
	genesis, scGenesis, staticCNNodes, staticPNNodes, staticENNodes, staticSCNNodes, staticSPNNodes, staticSENNodes,
	bridgeNodes, dockerImageId string, useFastHttp bool, networkId, parentChainId int,
	useGrafana bool, proxyNodeKeys, enNodeKeys, scnNodeKeys, spnNodeKeys, senNodeKeys []string, useTxGen bool,
	txGenOpt service.TxGenOption) *Homi {
	ist := &Homi{
		IPPrefix:   ipPrefix,
		EthStats:   service.NewEthStats(fmt.Sprintf("%v.9", ipPrefix), secret),
		UseGrafana: useGrafana,
		UseTxGen:   useTxGen,
	}
	ist.init(number, addresses, nodeKeys, genesis, scGenesis, staticCNNodes, staticPNNodes, staticENNodes,
		staticSCNNodes, staticSPNNodes, staticSENNodes, bridgeNodes, dockerImageId, useFastHttp, networkId,
		parentChainId, proxyNodeKeys, enNodeKeys, scnNodeKeys, spnNodeKeys, senNodeKeys, txGenOpt)
	return ist
}

func (ist *Homi) init(number int, addresses, nodeKeys []string, genesis, scGenesis, staticCNNodes, staticPNNodes,
	staticENNodes, staticSCNNodes, staticSPNNodes, staticSENNodes, bridgeNodes, dockerImageId string, useFastHttp bool,
	networkId, parentChainId int, proxyNodeKeys, enNodeKeys, scnNodeKeys, spnNodeKeys, senNodeKeys []string, txGenOpt service.TxGenOption) {
	var validatorNames []string
	for i := 0; i < number; i++ {
		s := service.NewValidator(i,
			genesis,
			"",
			addresses[i],
			nodeKeys[i],
			"",
			"",
			32323+i,
			8551+i,
			61001+i,
			ist.EthStats.Host(),
			// from subnet ip 10
			fmt.Sprintf("%v.%v", ist.IPPrefix, i+10),
			dockerImageId,
			useFastHttp,
			networkId,
			0,
			"CN",
			"cn",
			false,
		)

		staticCNNodes = strings.Replace(staticCNNodes, "0.0.0.0", s.IP, 1)
		ist.Services = append(ist.Services, s)
		validatorNames = append(validatorNames, s.Name)
	}

	numPNs := len(proxyNodeKeys)
	for i := 0; i < numPNs; i++ {
		s := service.NewValidator(i,
			genesis,
			"",
			"",
			proxyNodeKeys[i],
			"",
			"",
			32323+number+i,
			8551+number+i,
			61001+number+i,
			ist.EthStats.Host(),
			// from subnet ip 10
			fmt.Sprintf("%v.%v", ist.IPPrefix, number+i+10),
			dockerImageId,
			useFastHttp,
			networkId,
			0,
			"PN",
			"pn",
			false,
		)

		staticPNNodes = strings.Replace(staticPNNodes, "0.0.0.0", s.IP, 1)
		ist.Services = append(ist.Services, s)
		validatorNames = append(validatorNames, s.Name)
	}

	numENs := len(enNodeKeys)
	for i := 0; i < len(enNodeKeys); i++ {
		s := service.NewValidator(i,
			genesis,
			"",
			"",
			enNodeKeys[i],
			"",
			"",
			32323+number+numPNs+i,
			8551+number+numPNs+i,
			61001+number+numPNs+i,
			ist.EthStats.Host(),
			// from subnet ip 10
			fmt.Sprintf("%v.%v", ist.IPPrefix, number+numPNs+i+10),
			dockerImageId,
			useFastHttp,
			networkId,
			0,
			"EN",
			"en",
			false,
		)
		if i == 0 {
			bridgeNodes = strings.Replace(bridgeNodes, "0.0.0.0", s.IP, 1)
			bridgeNodes = strings.Replace(bridgeNodes, "32323", "50505", 1)
		}
		staticENNodes = strings.Replace(staticENNodes, "0.0.0.0", s.IP, 1)
		ist.Services = append(ist.Services, s)
		validatorNames = append(validatorNames, s.Name)
	}

	numSCNs := len(scnNodeKeys)
	scnParentChainId := parentChainId
	scnBridgeNodes := bridgeNodes
	if len(spnNodeKeys) > 0 || len(senNodeKeys) > 0 {
		scnParentChainId = 0
		scnBridgeNodes = ""
	}
	for i := 0; i < len(scnNodeKeys); i++ {
		s := service.NewValidator(i,
			"",
			scGenesis,
			"",
			scnNodeKeys[i],
			"",
			scnBridgeNodes,
			32323+number+numPNs+numENs+i,
			8551+number+numPNs+numENs+i,
			61001+number+numPNs+numENs+i,
			ist.EthStats.Host(),
			// from subnet ip 10
			fmt.Sprintf("%v.%v", ist.IPPrefix, number+numPNs+numENs+i+10),
			dockerImageId,
			useFastHttp,
			networkId,
			scnParentChainId,
			"SCN",
			"scn",
			false,
		)

		staticSCNNodes = strings.Replace(staticSCNNodes, "0.0.0.0", s.IP, 1)
		ist.Services = append(ist.Services, s)
		validatorNames = append(validatorNames, s.Name)
	}

	numSPNs := len(spnNodeKeys)
	spnParentChainId := parentChainId
	spnBridgeNodes := bridgeNodes
	if len(senNodeKeys) > 0 {
		spnParentChainId = 0
		spnBridgeNodes = ""
	}
	for i := 0; i < len(spnNodeKeys); i++ {
		s := service.NewValidator(i,
			"",
			scGenesis,
			"",
			spnNodeKeys[i],
			"",
			spnBridgeNodes,
			32323+number+numPNs+numENs+numSCNs+i,
			8551+number+numPNs+numENs+numSCNs+i,
			61001+number+numPNs+numENs+numSCNs+i,
			ist.EthStats.Host(),
			// from subnet ip 10
			fmt.Sprintf("%v.%v", ist.IPPrefix, number+numPNs+numENs+numSCNs+i+10),
			dockerImageId,
			useFastHttp,
			networkId,
			spnParentChainId,
			"SPN",
			"spn",
			false,
		)

		staticSPNNodes = strings.Replace(staticSPNNodes, "0.0.0.0", s.IP, 1)
		ist.Services = append(ist.Services, s)
		validatorNames = append(validatorNames, s.Name)
	}

	for i := 0; i < len(senNodeKeys); i++ {
		s := service.NewValidator(i,
			"",
			scGenesis,
			"",
			senNodeKeys[i],
			"",
			bridgeNodes,
			32323+number+numPNs+numENs+numSCNs+numSPNs+i,
			8551+number+numPNs+numENs+numSCNs+numSPNs+i,
			61001+number+numPNs+numENs+numSCNs+numSPNs+i,
			ist.EthStats.Host(),
			// from subnet ip 10
			fmt.Sprintf("%v.%v", ist.IPPrefix, number+numPNs+numENs+numSCNs+numSPNs+i+10),
			dockerImageId,
			useFastHttp,
			networkId,
			parentChainId,
			"SEN",
			"sen",
			false,
		)

		ist.Services = append(ist.Services, s)
		validatorNames = append(validatorNames, s.Name)
	}

	// update static nodes
	for i := range ist.Services {
		if ist.Services[i].NodeType == "scn" {
			ist.Services[i].StaticNodes = staticSCNNodes
		} else if ist.Services[i].NodeType == "spn" {
			ist.Services[i].StaticNodes = staticSCNNodes
		} else if ist.Services[i].NodeType == "sen" {
			ist.Services[i].StaticNodes = staticSPNNodes
		} else if ist.Services[i].NodeType == "en" {
			ist.Services[i].StaticNodes = staticPNNodes
		} else {
			ist.Services[i].StaticNodes = staticCNNodes
		}
	}

	ist.PrometheusService = service.NewPrometheusService(
		fmt.Sprintf("%v.%v", ist.IPPrefix, 9),
		validatorNames)

	if ist.UseGrafana {
		ist.GrafanaService = service.NewGrafanaService(fmt.Sprintf("%v.%v", ist.IPPrefix, 8))
	}

	ist.TxGenService = service.NewTxGenService(
		fmt.Sprintf("%v.%v", ist.IPPrefix, 7),
		fmt.Sprintf("http://%v.%v:8551", ist.IPPrefix, number+10),
		txGenOpt)
}

func (ist Homi) String() string {
	tmpl, err := template.New("istanbul").Parse(istanbulTemplate)
	if err != nil {
		fmt.Printf("Failed to parse template, %v", err)
		return ""
	}

	result := new(bytes.Buffer)
	err = tmpl.Execute(result, ist)
	if err != nil {
		fmt.Printf("Failed to render template, %v", err)
		return ""
	}

	return result.String()
}

var istanbulTemplate = `version: '3'
services:
  {{- range .Services }}
  {{ . }}
  {{- end }}
  {{- range .Proxies }}
  {{ . }}
  {{- end }}
  {{ .PrometheusService }}
  {{- if .UseGrafana }}
  {{ .GrafanaService }}
  {{- end }}
  {{- if .UseTxGen }}
  {{ .TxGenService }}
  {{- end }}
networks:
  app_net:
    driver: bridge
    ipam:
      driver: default
      config:
      - subnet: {{ .IPPrefix }}.0/24`

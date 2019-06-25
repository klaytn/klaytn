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

type KlayStats struct {
	Secret string
	IP     string
}

func NewEthStats(ip string, secret string) *KlayStats {
	return &KlayStats{
		IP:     ip,
		Secret: secret,
	}
}

func (c KlayStats) Host() string {
	return fmt.Sprintf("%v@%v:3000", c.Secret, c.IP)
}

func (c KlayStats) String() string {
	tmpl, err := template.New("klay_stats").Parse(klayStatsTemplate)
	if err != nil {
		fmt.Printf("Failed to parse template, %v", err)
		return ""
	}

	result := new(bytes.Buffer)
	err = tmpl.Execute(result, c)
	if err != nil {
		fmt.Printf("Failed to render template, %v", err)
		return ""
	}

	return result.String()
}

var klayStatsTemplate = `klay-stats:
    image: quay.io/amis/klaystats:latest
    ports:
      - '3000:3000'
    environment:
      - WS_SECRET={{ .Secret }}
    restart: always
    networks:
      app_net:
        ipv4_address: {{ .IP }}`

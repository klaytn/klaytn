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
	"strings"
	"text/template"
)

type Constellation struct {
	Identity   int
	Name       string
	IP         string
	Port       int
	OtherNodes string
	PublicKey  string
	PrivateKey string
	SocketPath string
	ConfigPath string
	Folder     string
	KeyPath    string
}

func (c *Constellation) SetOtherNodes(nodes []string) {
	c.OtherNodes = strings.Join(nodes, ",")
}

func (c Constellation) Host() string {
	return fmt.Sprintf("http://%v:%v/", c.IP, c.Port)
}

func (c Constellation) String() string {
	tmpl, err := template.New("constellation").Parse(constellationTemplate)
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

var constellationTemplate = `{{ .Name }}:
    hostname: {{ .Name }}
    image: quay.io/amis/constellation:latest
    ports:
      - '{{ .Port }}:{{ .Port }}'
    volumes:
      - {{ .Identity }}:{{ .Folder }}:z
      - .:/tmp/
    entrypoint:
      - /bin/sh
      - -c
      - |
        mkdir -p {{ .Folder }}
        echo "socket=\"{{ .SocketPath }}\"\npublickeys=[\"{{ .PublicKey }}\"]\n" > {{ .ConfigPath }}
        constellation-node --generatekeys={{ .KeyPath }}
        cp {{ .KeyPath }}.pub /tmp/tm{{ .Identity }}.pub
        constellation-node \
          --url={{ .Host }} \
          --port={{ .Port }} \
          --socket={{ .SocketPath }} \
          --othernodes={{ .OtherNodes }} \
          --publickeys={{ .PublicKey }} \
          --privatekeys={{ .PrivateKey }} \
          --storage={{ .Folder }} \
          --verbosity=4
    networks:
      app_net:
        ipv4_address: {{ .IP }}
    restart: always`

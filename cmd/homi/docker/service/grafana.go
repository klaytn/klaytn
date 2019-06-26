package service

import (
	"bytes"
	"fmt"
	"text/template"
)

type GrafanaService struct {
	IP string
}

func (gf GrafanaService) String() string {
	tmpl, err := template.New("GrafanaService").Parse(grafanaServiceTemplate)
	if err != nil {
		fmt.Printf("Failed to parse template, %v", err)
		return ""
	}

	result := new(bytes.Buffer)
	err = tmpl.Execute(result, gf)
	if err != nil {
		fmt.Printf("Failed to render template, %v", err)
		return ""
	}
	return result.String()
}

var grafanaServiceTemplate = `grafana:
    hostname: grafana
    image: grafana/grafana:5.2.2
    ports:
      - 3000:3000
    networks:
      app_net:
        ipv4_address: {{ .IP }}
    restart: "no"
`

func NewGrafanaService(ip string) *GrafanaService {
	return &GrafanaService{
		IP: ip,
	}
}

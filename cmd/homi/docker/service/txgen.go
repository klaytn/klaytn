package service

import (
	"bytes"
	"fmt"
	"text/template"
)

type TxGenService struct {
	Name       string
	IP         string
	TargetUrl  string
	Rate       int
	ThreadSize int
	ConnSize   int
	Duration   string
}

func NewTxGenService(ip string, targetUrl string, opt TxGenOption) *TxGenService {
	return &TxGenService{
		Name:       "txgen",
		IP:         ip,
		TargetUrl:  targetUrl,
		Rate:       opt.TxGenRate,
		ThreadSize: opt.TxGenThreadSize,
		ConnSize:   opt.TxGenConnSize,
		Duration:   opt.TxGenDuration,
	}
}

func (s TxGenService) String() string {
	tmpl, err := template.New("TxGen").Parse(dockerComposeTemplate)
	if err != nil {
		fmt.Printf("Failed to parse template, %v", err)
		return ""
	}

	result := new(bytes.Buffer)
	err = tmpl.Execute(result, s)
	if err != nil {
		fmt.Printf("Failed to render template, %v", err)
		return ""
	}
	return result.String()
}

var dockerComposeTemplate = `{{ .Name }}:
    hostname: {{ .Name }}
    image: 428948643293.dkr.ecr.ap-northeast-2.amazonaws.com/groundx/txgen:latest
    entrypoint: ["txgen", "-r", "{{ .Rate }}", "-t", "{{ .ThreadSize }}", "-c", "{{ .ConnSize }}", "-d", "{{ .Duration }}", "{{ .TargetUrl }}"]
    networks:
      app_net:
        ipv4_address: {{ .IP }}
    restart: "no"
`

type TxGenOption struct {
	TxGenThreadSize int
	TxGenDuration   string
	TxGenConnSize   int
	TxGenRate       int
}

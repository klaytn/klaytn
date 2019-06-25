package service

import (
	"bytes"
	"fmt"
	"text/template"
)

// docker-compose service
type PrometheusService struct {
	Config PrometheusConfig
	Name   string
	IP     string
}

func (prom PrometheusService) String() string {
	tmpl, err := template.New("PrometheusService").Parse(prometheusServieTemplate)
	if err != nil {
		fmt.Printf("Failed to parse template, %v", err)
		return ""
	}

	result := new(bytes.Buffer)
	err = tmpl.Execute(result, prom)
	if err != nil {
		fmt.Printf("Failed to render template, %v", err)
		return ""
	}
	return result.String()
}

var prometheusServieTemplate = `{{ .Name }}:
    hostname: {{ .Name }}
    image: prom/prometheus:v2.3.2
    ports:
      - 9090:9090
    volumes:
      - "./prometheus.yml:/etc/prometheus/prometheus.yml"
    networks:
      app_net:
        ipv4_address: {{ .IP }}
    restart: "no"
`

func NewPrometheusService(ip string, targetHostNames []string) *PrometheusService {
	return &PrometheusService{
		Name:   "prometheus",
		IP:     ip,
		Config: NewPrometheusConfig(targetHostNames),
	}
}

// prometheus server config
type PrometheusConfig struct {
	Targets []string
}

func NewPrometheusConfig(targetHostNames []string) PrometheusConfig {
	return PrometheusConfig{
		Targets: targetHostNames,
	}
}

func (promConfig PrometheusConfig) String() string {
	tmpl, err := template.New("PrometheusConfig").Parse(prometheusConfigTemplate)
	if err != nil {
		fmt.Printf("Failed to parse template, %v", err)
		return ""
	}

	result := new(bytes.Buffer)
	err = tmpl.Execute(result, promConfig)
	if err != nil {
		fmt.Printf("Failed to render template, %v", err)
		return ""
	}

	return result.String()
}

var prometheusConfigTemplate = `global:
  scrape_interval:     15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'klaytn'
    static_configs:
    - targets:
      {{- range .Targets }}
      - "{{ . }}:61001"
      {{- end }}
`

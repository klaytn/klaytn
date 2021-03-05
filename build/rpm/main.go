package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/klaytn/klaytn/params"
	"github.com/urfave/cli"
)

const (
	CN   = "kcn"
	PN   = "kpn"
	EN   = "ken"
	SCN  = "kscn"
	SPN  = "kspn"
	SEN  = "ksen"
	BN   = "kbn"
	HOMI = "homi"
	GEN  = "kgen"
)

type NodeInfo struct {
	daemon      string
	summary     string
	description string
}

var BINARY_TYPE = map[string]NodeInfo{
	CN: {"kcnd",
		"Klaytn consensus node daemon",
		"kcnd is a daemon for Klaytn consensus node (kcn). For more information, please refer to https://docs.klaytn.com.",
	},
	PN: {"kpnd",
		"Klaytn proxy node daemon",
		"kpnd is a daemon for Klaytn proxy node (kpn). For more information, please refer to https://docs.klaytn.com.",
	},
	EN: {"kend",
		"Klaytn endpoint node daemon",
		"kend is a daemon for Klaytn endpoint node (ken). For more information, please refer to https://docs.klaytn.com.",
	},
	SCN: {"kscnd",
		"Klaytn servicechain consensus node daemon",
		"kscnd is a daemon for Klaytn servicechain consensus node (kscn). For more information, please refer to https://docs.klaytn.com.",
	},
	SPN: {"kspnd",
		"Klaytn servicechain proxy node daemon",
		"kspnd is a daemon for Klaytn servicechain proxy node (kspn). For more information, please refer to https://docs.klaytn.com.",
	},
	SEN: {"ksend",
		"Klaytn servicechain endpoint node daemon",
		"ksend is a daemon for Klaytn servicechain endpoint node (ksen). For more information, please refer to https://docs.klaytn.com.",
	},
	BN: {"kbnd",
		"Klaytn boot node daemon",
		"kbnd is a daemon for Klaytn boot node (kbn). For more information, please refer to https://docs.klaytn.com.",
	},
	HOMI: {"homi",
		"genesis.json generator",
		"homi is a generator of genesis.json."},
	GEN: {"kgen",
		"private key generator",
		"kgen is a generator of private keys.",
	},
}

type RpmSpec struct {
	BuildNumber int
	Version     string
	Name        string
	Summary     string
	MakeTarget  string
	ProgramName string // kcn, kpn, ken, kscn, kspn, ksen, kbn
	DaemonName  string // kcnd, kpnd, kend, kscnd, kspnd, ksend, kbnd
	PostFix     string // baobab
	Description string
}

func (r RpmSpec) String() string {
	tmpl, err := template.New("rpmspec").Parse(rpmSpecTemplate)
	if err != nil {
		fmt.Printf("Failed to parse template, %v", err)
		return ""
	}

	result := new(bytes.Buffer)
	err = tmpl.Execute(result, r)
	if err != nil {
		fmt.Printf("Failed to render template, %v", err)
		return ""
	}
	return result.String()
}

func main() {
	app := cli.NewApp()
	app.Name = "klaytn_rpmtool"
	app.Version = "0.2"
	app.Commands = []cli.Command{
		{
			Name:    "gen_spec",
			Aliases: []string{"a"},
			Usage:   "generate rpm spec file",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "binary_type",
					Usage: "Klaytn binary type (kcn, kpn, ken, kscn, kspn, ksen, kbn, kgen, homi)",
				},
				cli.BoolFlag{
					Name:  "devel",
					Usage: "generate spec for devel version",
				},
				cli.BoolFlag{
					Name:  "baobab",
					Usage: "generate spec for baobab version",
				},
				cli.IntFlag{
					Name:  "build_num",
					Usage: "build number",
				},
			},
			Action: genspec,
		},
		{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "return klaytn version",
			Action: func(c *cli.Context) error {
				fmt.Print(params.Version)
				return nil
			},
		},
		{
			Name:    "release_num",
			Aliases: []string{"r"},
			Usage:   "return klaytn release number",
			Action: func(c *cli.Context) error {
				fmt.Print(params.ReleaseNum)
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func genspec(c *cli.Context) error {
	rpmSpec := new(RpmSpec)

	binaryType := c.String("binary_type")
	if _, ok := BINARY_TYPE[binaryType]; ok != true {
		return fmt.Errorf("binary_type[\"%s\"] is not supported. Use --binary_type [kcn, kpn, ken, kscn, kspn, ksen, kbn, kgen, homi]", binaryType)
	}

	rpmSpec.ProgramName = strings.ToLower(binaryType)
	rpmSpec.DaemonName = BINARY_TYPE[binaryType].daemon
	rpmSpec.PostFix = ""

	if c.Bool("devel") {
		buildNum := c.Int("build_num")
		if buildNum == 0 {
			fmt.Println("BuildNumber should be set")
			os.Exit(1)
		}
		rpmSpec.BuildNumber = buildNum
		rpmSpec.Name = BINARY_TYPE[binaryType].daemon + "-devel"
	} else if c.Bool("baobab") {
		rpmSpec.BuildNumber = params.ReleaseNum
		rpmSpec.Name = BINARY_TYPE[binaryType].daemon + "-baobab"
		rpmSpec.PostFix = "_baobab"
	} else {
		rpmSpec.BuildNumber = params.ReleaseNum
		rpmSpec.Name = BINARY_TYPE[binaryType].daemon
	}
	rpmSpec.Summary = BINARY_TYPE[binaryType].summary
	rpmSpec.Description = BINARY_TYPE[binaryType].description
	rpmSpec.Version = params.Version
	fmt.Println(rpmSpec)
	return nil
}

var rpmSpecTemplate = `Name:               {{ .Name }}
Version:            {{ .Version }}
Release:            {{ .BuildNumber }}%{?dist}
Summary:            {{ .Summary }}

Group:              Application/blockchain
License:            GNU
URL:                https://www.klaytn.com
Source0:            %{name}-%{version}.tar.gz
BuildRoot:          %(mktemp -ud %{_tmppath}/%{name}-%{version}-%{release}-XXXXXX)

%description
  {{ .Description }}

%prep
%setup -q

%build
make {{ .ProgramName }}

%define is_daemon %( if [ {{ .ProgramName }} != {{ .DaemonName }} ]; then echo "1"; else echo "0"; fi )

%install
mkdir -p $RPM_BUILD_ROOT/usr/bin
mkdir -p $RPM_BUILD_ROOT/etc/{{ .DaemonName }}/conf
mkdir -p $RPM_BUILD_ROOT/etc/init.d
mkdir -p $RPM_BUILD_ROOT/var/log/{{ .DaemonName }}

cp build/bin/{{ .ProgramName }} $RPM_BUILD_ROOT/usr/bin/{{ .ProgramName }}
%if %is_daemon
cp build/rpm/etc/init.d/{{ .DaemonName }} $RPM_BUILD_ROOT/etc/init.d/{{ .DaemonName }}
cp build/rpm/etc/{{ .DaemonName }}/conf/{{ .DaemonName }}{{ .PostFix }}.conf $RPM_BUILD_ROOT/etc/{{ .DaemonName }}/conf/{{ .DaemonName }}.conf
%endif

%files
%attr(755, -, -) /usr/bin/{{ .ProgramName }}
%if %is_daemon
%config(noreplace) %attr(644, -, -) /etc/{{ .DaemonName }}/conf/{{ .DaemonName }}.conf
%attr(754, -, -) /etc/init.d/{{ .DaemonName }}
%endif
%exclude /usr/local/var/lib/rpm/*
%exclude /usr/local/var/lib/rpm/.*
%exclude /usr/local/var/tmp/*

%pre
%if %is_daemon
if [ $1 -eq 2 ]; then
	# Package upgrade
	systemctl stop {{ .DaemonName }}.service > /dev/null 2>&1
fi
%endif

%post
%if %is_daemon
if [ $1 -eq 1 ]; then
	# Package installation
	systemctl daemon-reload >/dev/null 2>&1
fi
if [ $1 -eq 2 ]; then
	# Package upgrade
	systemctl daemon-reload >/dev/null 2>&1
fi
%endif

%preun
%if %is_daemon
if [ $1 -eq 0 ]; then
	# Package removal, not upgrade
	systemctl --no-reload disable {{ .DaemonName }}.service > /dev/null 2>&1
	systemctl stop {{ .DaemonName }}.service > /dev/null 2>&1
fi
%endif

%postun
%if %is_daemon
if [ $1 -eq 0 ]; then
	# Package uninstallation
	systemctl daemon-reload >/dev/null 2>&1
fi
%endif
`

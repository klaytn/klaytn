package node

import (
	"ground-x/go-gxplatform/event"
	"ground-x/go-gxplatform/log"
	"path/filepath"
	"strings"
	"errors"
	"ground-x/go-gxplatform/rpc"
)

type VNode struct {
	Node
}

// New creates a new P2P node, ready for protocol registration.
func NewVNode(conf *Config) (*VNode, error) {
	// Copy config and resolve the datadir so future changes to the current
	// working directory don't affect the node.
	confCopy := *conf
	conf = &confCopy
	if conf.DataDir != "" {
		absdatadir, err := filepath.Abs(conf.DataDir)
		if err != nil {
			return nil, err
		}
		conf.DataDir = absdatadir
	}
	// Ensure that the instance name doesn't cause weird conflicts with
	// other files in the data directory.
	if strings.ContainsAny(conf.Name, `/\`) {
		return nil, errors.New(`Config.Name must not contain '/' or '\'`)
	}
	if conf.Name == datadirDefaultKeyStore {
		return nil, errors.New(`Config.Name cannot be "` + datadirDefaultKeyStore + `"`)
	}
	if strings.HasSuffix(conf.Name, ".ipc") {
		return nil, errors.New(`Config.Name cannot end in ".ipc"`)
	}

	// Ensure that the AccountManager method works before the node has started.
	// We rely on this in cmd/geth.
	am, ephemeralKeystore, err := makeAccountManager(conf)
	if err != nil {
		return nil, err
	}
	if conf.Logger == nil {
		conf.Logger = log.New()
	}

	// Note: any interaction with Config that would create/touch files
	// in the data directory or instance directory is delayed until Start.
	return &VNode{
		Node{accman: am,
			ephemeralKeystore: ephemeralKeystore,
			config: conf,
			serviceFuncs: []ServiceConstructor{},
			ipcEndpoint: conf.IPCEndpoint(),
			httpEndpoint: conf.HTTPEndpoint(),
			wsEndpoint: conf.WSEndpoint(),
			eventmux: new(event.TypeMux),
			log: conf.Logger,
		}}, nil
}

func (n *VNode) apis() []rpc.API {
	return []rpc.API{
		{
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPublicDebugAPI(&n.Node),
			Public:    true,
		}, {
			Namespace: "web3",
			Version:   "1.0",
			Service:   NewPublicWeb3API(&n.Node),
			Public:    true,
		},
	}
}
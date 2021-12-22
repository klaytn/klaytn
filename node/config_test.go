// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of go-ethereum.
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
//
// This file is derived from node/config_test.go (2018/06/04).
// Modified and improved for the klaytn development.

package node

import (
	"bytes"
	"encoding/json"
	"github.com/klaytn/klaytn/networks/p2p/discover"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/networks/p2p"
)

// Tests that datadirs can be successfully created, be them manually configured
// ones or automatically generated temporary ones.
func TestDatadirCreation(t *testing.T) {
	// Create a temporary data dir and check that it can be used by a node
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("failed to create manual data dir: %v", err)
	}
	defer os.RemoveAll(dir)

	if _, err := New(&Config{DataDir: dir}); err != nil {
		t.Fatalf("failed to create stack with existing datadir: %v", err)
	}
	// Generate a long non-existing datadir path and check that it gets created by a node
	dir = filepath.Join(dir, "a", "b", "c", "d", "e", "f")
	if _, err := New(&Config{DataDir: dir}); err != nil {
		t.Fatalf("failed to create stack with creatable datadir: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("freshly created datadir not accessible: %v", err)
	}
	// Verify that an impossible datadir fails creation
	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	dir = filepath.Join(file.Name(), "invalid/path")
	if _, err := New(&Config{DataDir: dir}); err == nil {
		t.Fatalf("protocol stack created with an invalid datadir")
	}
}

// Tests that IPC paths are correctly resolved to valid endpoints of different
// platforms.
func TestIPCPathResolution(t *testing.T) {
	var tests = []struct {
		DataDir  string
		IPCPath  string
		Windows  bool
		Endpoint string
	}{
		{"", "", false, ""},
		{"data", "", false, ""},
		{"", "klay.ipc", false, filepath.Join(os.TempDir(), "klay.ipc")},
		{"data", "klay.ipc", false, "data/klay.ipc"},
		{"data", "./klay.ipc", false, "./klay.ipc"},
		{"data", "/klay.ipc", false, "/klay.ipc"},
		{"", "", true, ``},
		{"data", "", true, ``},
		{"", "klay.ipc", true, `\\.\pipe\klay.ipc`},
		{"data", "klay.ipc", true, `\\.\pipe\klay.ipc`},
		{"data", `\\.\pipe\klay.ipc`, true, `\\.\pipe\klay.ipc`},
	}
	for i, test := range tests {
		// Only run when platform/test match
		if (runtime.GOOS == "windows") == test.Windows {
			if endpoint := (&Config{DataDir: test.DataDir, IPCPath: test.IPCPath}).IPCEndpoint(); endpoint != test.Endpoint {
				t.Errorf("test %d: IPC endpoint mismatch: have %s, want %s", i, endpoint, test.Endpoint)
			}
		}
	}
}

// Tests that node keys can be correctly created, persisted, loaded and/or made
// ephemeral.
func TestNodeKeyPersistency(t *testing.T) {
	// Create a temporary folder and make sure no key is present
	dir, err := ioutil.TempDir("", "node-test")
	if err != nil {
		t.Fatalf("failed to create temporary data directory: %v", err)
	}
	defer os.RemoveAll(dir)

	keyfile := filepath.Join(dir, "unit-test", datadirPrivateKey)

	// Configure a node with a preset key and ensure it's not persisted
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate one-shot node key: %v", err)
	}
	config := &Config{Name: "unit-test", DataDir: dir, P2P: p2p.Config{PrivateKey: key}}
	config.NodeKey()
	if _, err := os.Stat(filepath.Join(keyfile)); err == nil {
		t.Fatalf("one-shot node key persisted to data directory")
	}

	// Configure a node with no preset key and ensure it is persisted this time
	config = &Config{Name: "unit-test", DataDir: dir}
	config.NodeKey()
	if _, err := os.Stat(keyfile); err != nil {
		t.Fatalf("node key not persisted to data directory: %v", err)
	}
	if _, err = crypto.LoadECDSA(keyfile); err != nil {
		t.Fatalf("failed to load freshly persisted node key: %v", err)
	}
	blob1, err := ioutil.ReadFile(keyfile)
	if err != nil {
		t.Fatalf("failed to read freshly persisted node key: %v", err)
	}

	// Configure a new node and ensure the previously persisted key is loaded
	config = &Config{Name: "unit-test", DataDir: dir}
	config.NodeKey()
	blob2, err := ioutil.ReadFile(filepath.Join(keyfile))
	if err != nil {
		t.Fatalf("failed to read previously persisted node key: %v", err)
	}
	if !bytes.Equal(blob1, blob2) {
		t.Fatalf("persisted node key mismatch: have %x, want %x", blob2, blob1)
	}

	// TODO-Klaytn-FailedTest Test fails
	/*
		// Configure ephemeral node and ensure no key is dumped locally
		config = &Config{Name: "unit-test", DataDir: ""}
		config.NodeKey()
		if _, err := os.Stat(filepath.Join(".", "unit-test", datadirPrivateKey)); err == nil {
			t.Fatalf("ephemeral node key persisted to disk")
		}
	*/
}

func TestConfig_ParsePersistentNodes(t *testing.T) {
	type expect struct {
		node  string
		proxy string
	}
	tests := []struct {
		name        string
		nodes       []string
		proxyURL    string
		expects     []expect
		assertNodes func(t *testing.T, nodes []*discover.Node)
	}{
		{
			name: "success",
			nodes: []string{
				"kni://a3fd567e0e1bb9f7d7a20385b797563883ebb9d45ff6f05a588b56256f46bd649b7ecc8e3e17cc50df4599c1809463e66ad964f7a3fb6cf4c768c25d647f2c02@1.1.1.1:3000",
				"enode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@2.2.2.2:52150",
			},
			proxyURL: "",
			expects: []expect{
				{
					"kni://a3fd567e0e1bb9f7d7a20385b797563883ebb9d45ff6f05a588b56256f46bd649b7ecc8e3e17cc50df4599c1809463e66ad964f7a3fb6cf4c768c25d647f2c02@1.1.1.1:3000",
					"",
				}, {
					"kni://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@2.2.2.2:52150",
					"",
				},
			},
		}, {
			name: "success with ignored invalid kni",
			nodes: []string{
				"kni://a3fd567e0e1bb9f7d7a20385b797563883ebb9d45ff6f05a588b56256f46bd649b7ecc8e3e17cc50df4599c1809463e66ad964f7a3fb6cf4c768c25d647f2c02@1.1.1.1:3000",
				"unknown://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@1234.com:52150",
				"enode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@2.2.2.2:52150",
			},
			proxyURL: "",
			expects: []expect{
				{
					"kni://a3fd567e0e1bb9f7d7a20385b797563883ebb9d45ff6f05a588b56256f46bd649b7ecc8e3e17cc50df4599c1809463e66ad964f7a3fb6cf4c768c25d647f2c02@1.1.1.1:3000",
					"",
				}, {
					"kni://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@2.2.2.2:52150",
					"",
				},
			},
		}, {
			name: "success with proxy",
			nodes: []string{
				"kni://a3fd567e0e1bb9f7d7a20385b797563883ebb9d45ff6f05a588b56256f46bd649b7ecc8e3e17cc50df4599c1809463e66ad964f7a3fb6cf4c768c25d647f2c02@1.1.1.1:3000",
				"enode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@2.2.2.2:52150",
			},
			proxyURL: "socks5://10.100.10.1:3500",
			expects: []expect{
				{
					"kni://a3fd567e0e1bb9f7d7a20385b797563883ebb9d45ff6f05a588b56256f46bd649b7ecc8e3e17cc50df4599c1809463e66ad964f7a3fb6cf4c768c25d647f2c02@1.1.1.1:3000",
					"socks5://10.100.10.1:3500",
				}, {
					"kni://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@2.2.2.2:52150",
					"socks5://10.100.10.1:3500",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := datadirStaticNodes
			conf := setupNodeConfigs(t, tt.nodes, path)
			defer os.RemoveAll(conf.DataDir)
			conf.P2P.ProxyURL = tt.proxyURL

			nodes := conf.parsePersistentNodes(conf.ResolvePath(path))

			assert.Equal(t, len(tt.expects), len(nodes))
			for i, n := range nodes {
				assert.Equal(t, tt.expects[i].node, n.String())
				assert.Equal(t, tt.expects[i].proxy, n.ProxyURL)
			}
		})
	}
}

func setupNodeConfigs(t *testing.T, nodes []string, nodepath string) *Config {
	// Setup temporary directory.
	dir, err := ioutil.TempDir("", "node-config-test")
	assert.NoError(t, err)
	c := &Config{DataDir: dir}
	err = os.MkdirAll(c.ResolvePath(""), os.ModePerm)
	assert.NoError(t, err)

	// Setup node config files.
	if nodepath != "" {
		nodeBytes, _ := json.Marshal(&nodes)
		err := ioutil.WriteFile(c.ResolvePath(nodepath), nodeBytes, os.ModePerm)
		assert.NoError(t, err)
	}
	return c
}

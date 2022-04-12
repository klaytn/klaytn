// Modifications Copyright 2018 The klaytn Authors
// Copyright 2019 The go-ethereum Authors
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
//
// This file is derived from common/compiler/solidity.go (2018/06/04).
// Modified and improved for the klaytn development.

// Package compiler wraps the Solidity compiler executables (solc).
package compiler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"

	"github.com/klaytn/klaytn/log"
)

var logger = log.NewModuleLogger(log.Common)

// Solidity contains information about the solidity compiler.
type Solidity struct {
	Path, Version, FullVersion string
	Major, Minor, Patch        int
}

// --combined-output format
type solcOutput struct {
	Contracts map[string]struct {
		BinRuntime                                  string `json:"bin-runtime"`
		SrcMapRuntime                               string `json:"srcmap-runtime"`
		Bin, SrcMap, Abi, Devdoc, Userdoc, Metadata string
		Hashes                                      map[string]string
	}
	Version string `json:"version"`
}

// solidity v.0.8 changes the way ABI, Devdoc and Userdoc are serialized
type solcOutputV8 struct {
	Contracts map[string]struct {
		BinRuntime            string `json:"bin-runtime"`
		SrcMapRuntime         string `json:"srcmap-runtime"`
		Bin, SrcMap, Metadata string
		Abi                   interface{}
		Devdoc                interface{}
		Userdoc               interface{}
		Hashes                map[string]string
	}
	Version string `json:"version"`
}

func (s *Solidity) makeArgs() []string {
	p := []string{
		"--combined-json", "bin,bin-runtime,srcmap,srcmap-runtime,abi,userdoc,devdoc",
		"--optimize",                  // code optimizer switched on
		"--allow-paths", "., ./, ../", //default to support relative path： ./  ../  .
	}
	if s.Major > 0 || s.Minor > 4 || s.Patch > 6 {
		p[1] += ",metadata,hashes"
	}
	return p
}

// SolidityVersion runs solc and parses its version output.
func SolidityVersion(solc string) (*Solidity, error) {
	if solc == "" {
		solc = "solc"
	}
	var out bytes.Buffer
	cmd := exec.Command(solc, "--version")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	matches := versionRegexp.FindStringSubmatch(out.String())
	if len(matches) != 4 {
		return nil, fmt.Errorf("can't parse solc version %q", out.String())
	}
	s := &Solidity{Path: cmd.Path, FullVersion: out.String(), Version: matches[0]}
	if s.Major, err = strconv.Atoi(matches[1]); err != nil {
		return nil, err
	}
	if s.Minor, err = strconv.Atoi(matches[2]); err != nil {
		return nil, err
	}
	if s.Patch, err = strconv.Atoi(matches[3]); err != nil {
		return nil, err
	}
	return s, nil
}

// CompileSolidityString builds and returns all the contracts contained within a source string.
func CompileSolidityString(solc, source string) (map[string]*Contract, error) {
	if len(source) == 0 {
		return nil, errors.New("solc: empty source string")
	}
	s, err := SolidityVersion(solc)
	if err != nil {
		return nil, err
	}
	args := append(s.makeArgs(), "--")
	cmd := exec.Command(s.Path, append(args, "-")...)
	cmd.Stdin = strings.NewReader(source)
	return s.run(cmd, source)
}

// CompileSolidity compiles all given Solidity source files.
func CompileSolidity(solc string, sourcefiles ...string) (map[string]*Contract, error) {
	if len(sourcefiles) == 0 {
		return nil, errors.New("solc: no source files")
	}
	source, err := slurpFiles(sourcefiles)
	if err != nil {
		return nil, err
	}
	s, err := SolidityVersion(solc)
	if err != nil {
		return nil, err
	}
	args := append(s.makeArgs(), "--")
	cmd := exec.Command(s.Path, append(args, sourcefiles...)...)
	return s.run(cmd, source)
}

// CompileSolidityOrLoad compiles all given Solidity source files.
// If suitable compiler is not available, try to load combinedJSON stored as file.
// The combinedJSON file should be named as *.sol.json.
// Create combinedJSON with following command:
//   solc --combined-json bin,bin-runtime,srcmap,srcmap-runtime,abi,userdoc,devdoc \
//       --optimize --allow-paths '., ./, ../' test.sol > test.sol.json
func CompileSolidityOrLoad(solc string, sourcefiles ...string) (map[string]*Contract, error) {
	// Extract solidity version requirements from source codes
	if len(sourcefiles) == 0 {
		return nil, errors.New("solc: no source files")
	}
	source, err := slurpFiles(sourcefiles)
	if err != nil {
		return nil, err
	}

	// Find available compiler, if any
	s, err := SolidityVersion(solc)
	if err != nil {
		logger.Warn("Solidity compiler not found. Loading from file.", "err", err)
		return loadCombinedJSON(source, sourcefiles...)
	}

	args := append(s.makeArgs(), "--")
	cmd := exec.Command(s.Path, append(args, sourcefiles...)...)
	contracts, err := s.run(cmd, source)
	if err != nil {
		logger.Warn("Solidity compiler cannot compile source versions. Loading from file.", "err", err)
		return loadCombinedJSON(source, sourcefiles...)
	}
	return contracts, nil
}

// Search for CombinedJSON at <solidity_file_name>.json
func loadCombinedJSON(source string, sourcefiles ...string) (map[string]*Contract, error) {
	path := sourcefiles[0] + ".json"
	combinedJSON, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Cannot open combined json at %s", path)
	}

	// Find compiler version from the loaded json
	v := struct {
		Version string `json:"version"`
	}{}
	if err := json.Unmarshal(combinedJSON, &v); err != nil {
		return nil, fmt.Errorf("can't parse combined json version %s", err)
	}
	matches := versionRegexp.FindStringSubmatch(v.Version)
	if len(matches) != 4 {
		return nil, fmt.Errorf("can't parse combined json version %q", v.Version)
	}
	version := matches[0]

	return ParseCombinedJSON(combinedJSON, source, version, version, "")
}

func (s *Solidity) run(cmd *exec.Cmd, source string) (map[string]*Contract, error) {
	var stderr, stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("solc: %v\n%s", err, stderr.Bytes())
	}

	return ParseCombinedJSON(stdout.Bytes(), source, s.Version, s.Version, strings.Join(s.makeArgs(), " "))
}

// ParseCombinedJSON takes the direct output of a solc --combined-output run and
// parses it into a map of string contract name to Contract structs. The
// provided source, language and compiler version, and compiler options are all
// passed through into the Contract structs.
//
// The solc output is expected to contain ABI, source mapping, user docs, and dev docs.
//
// Returns an error if the JSON is malformed or missing data, or if the JSON
// embedded within the JSON is malformed.
func ParseCombinedJSON(combinedJSON []byte, source string, languageVersion string, compilerVersion string, compilerOptions string) (map[string]*Contract, error) {
	var output solcOutput
	if err := json.Unmarshal(combinedJSON, &output); err != nil {
		// Try to parse the output with the new solidity v.0.8.0 rules
		return parseCombinedJSONV8(combinedJSON, source, languageVersion, compilerVersion, compilerOptions)
	}
	// Compilation succeeded, assemble and return the contracts.
	contracts := make(map[string]*Contract)
	for name, info := range output.Contracts {
		// Parse the individual compilation results.
		var abi interface{}
		if err := json.Unmarshal([]byte(info.Abi), &abi); err != nil {
			return nil, fmt.Errorf("solc: error reading abi definition (%v)", err)
		}
		var userdoc, devdoc interface{}
		json.Unmarshal([]byte(info.Userdoc), &userdoc)
		json.Unmarshal([]byte(info.Devdoc), &devdoc)

		contracts[name] = &Contract{
			Code:        "0x" + info.Bin,
			RuntimeCode: "0x" + info.BinRuntime,
			Hashes:      info.Hashes,
			Info: ContractInfo{
				Source:          source,
				Language:        "Solidity",
				LanguageVersion: languageVersion,
				CompilerVersion: compilerVersion,
				CompilerOptions: compilerOptions,
				SrcMap:          info.SrcMap,
				SrcMapRuntime:   info.SrcMapRuntime,
				AbiDefinition:   abi,
				UserDoc:         userdoc,
				DeveloperDoc:    devdoc,
				Metadata:        info.Metadata,
			},
		}
	}
	return contracts, nil
}

// parseCombinedJSONV8 parses the direct output of solc --combined-output
// and parses it using the rules from solidity v.0.8.0 and later.
func parseCombinedJSONV8(combinedJSON []byte, source string, languageVersion string, compilerVersion string, compilerOptions string) (map[string]*Contract, error) {
	var output solcOutputV8
	if err := json.Unmarshal(combinedJSON, &output); err != nil {
		return nil, err
	}
	// Compilation succeeded, assemble and return the contracts.
	contracts := make(map[string]*Contract)
	for name, info := range output.Contracts {
		contracts[name] = &Contract{
			Code:        "0x" + info.Bin,
			RuntimeCode: "0x" + info.BinRuntime,
			Hashes:      info.Hashes,
			Info: ContractInfo{
				Source:          source,
				Language:        "Solidity",
				LanguageVersion: languageVersion,
				CompilerVersion: compilerVersion,
				CompilerOptions: compilerOptions,
				SrcMap:          info.SrcMap,
				SrcMapRuntime:   info.SrcMapRuntime,
				AbiDefinition:   info.Abi,
				UserDoc:         info.Userdoc,
				DeveloperDoc:    info.Devdoc,
				Metadata:        info.Metadata,
			},
		}
	}
	return contracts, nil
}

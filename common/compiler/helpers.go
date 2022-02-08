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

// Package compiler wraps the Solidity and Vyper compiler executables (solc; vyper).
package compiler

import (
	"bytes"
	"io/ioutil"
	"regexp"

	semver "github.com/Masterminds/semver/v3"
)

var versionRegexp = regexp.MustCompile(`([0-9]+)\.([0-9]+)\.([0-9]+)`)
var pragmaVersionRegexp = regexp.MustCompile(`(?m)^pragma solidity (.+)\;$`)

// Contract contains information about a compiled contract, alongside its code and runtime code.
type Contract struct {
	Code        string            `json:"code"`
	RuntimeCode string            `json:"runtime-code"`
	Info        ContractInfo      `json:"info"`
	Hashes      map[string]string `json:"hashes"`
}

// ContractInfo contains information about a compiled contract, including access
// to the ABI definition, source mapping, user and developer docs, and metadata.
//
// Depending on the source, language version, compiler version, and compiler
// options will provide information about how the contract was compiled.
type ContractInfo struct {
	Source          string      `json:"source"`
	Language        string      `json:"language"`
	LanguageVersion string      `json:"languageVersion"`
	CompilerVersion string      `json:"compilerVersion"`
	CompilerOptions string      `json:"compilerOptions"`
	SrcMap          interface{} `json:"srcMap"`
	SrcMapRuntime   string      `json:"srcMapRuntime"`
	AbiDefinition   interface{} `json:"abiDefinition"`
	UserDoc         interface{} `json:"userDoc"`
	DeveloperDoc    interface{} `json:"developerDoc"`
	Metadata        string      `json:"metadata"`
}

func slurpFiles(files []string) (string, error) {
	var concat bytes.Buffer
	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return "", err
		}
		concat.Write(content)
	}
	return concat.String(), nil
}

func extractSourceVersion(source string) []string {
	matches := pragmaVersionRegexp.FindAllSubmatch([]byte(source), -1)
	versions := make([]string, 0)
	for _, match := range matches {
		versions = append(versions, string(match[1]))
	}
	return versions
}

func solcCanCompile(solcVersion string, sourceVersions []string) (bool, error) {
	v, err := semver.NewVersion(solcVersion)
	if err != nil {
		return false, err
	}

	for _, sourceVersion := range sourceVersions {
		c, err := semver.NewConstraint(sourceVersion)
		if err != nil {
			return false, err
		}
		if !c.Check(v) {
			return false, nil
		}
	}
	return true, nil
}

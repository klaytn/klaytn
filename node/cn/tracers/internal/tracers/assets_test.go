// Copyright 2021 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

package tracers

import (
	"crypto/sha1"
	"io/ioutil"
	"os/exec"
	"testing"

	"log"

	"github.com/stretchr/testify/assert"
)

func TestBindata(t *testing.T) {
	originalData, err := ioutil.ReadFile("./assets.go")

	if err != nil {
		log.Fatalln(err)
	}

	originalHash := sha1.New()
	originalHash.Write(originalData)
	cmd := exec.Command("go-bindata", "-nometadata", "-o", "assets.go", "-pkg", "tracers", "-ignore", "tracers.go", "-ignore", "assets.go", "./...")

	err = cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}

	cmd = exec.Command("gofmt", "-w", "-s", "assets.go")

	err = cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}

	updatedData, err := ioutil.ReadFile("./assets.go")

	if err != nil {
		log.Fatalln(err)
	}

	updatedHash := sha1.New()
	updatedHash.Write(updatedData)

	assert.Equal(t, updatedHash, originalHash)
}

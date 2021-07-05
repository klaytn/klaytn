package deps

import (
	"crypto/sha1"
	"io/ioutil"
	"log"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBindata(t *testing.T) {
	originalData, err := ioutil.ReadFile("./bindata.go")

	if err != nil {
		log.Fatalln(err)
	}

	originalHash := sha1.New()
	originalHash.Write(originalData)
	cmd := exec.Command("go-bindata", "-nometadata", "-pkg", "deps", "-o", "bindata.go", "bignumber.js", "web3.js")

	err = cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}

	cmd = exec.Command("gofmt", "-w", "-s", "bindata.go")

	err = cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}

	updatedData, err := ioutil.ReadFile("./bindata.go")

	if err != nil {
		log.Fatalln(err)
	}

	updatedHash := sha1.New()
	updatedHash.Write(updatedData)

	assert.Equal(t, updatedHash, originalHash)
}

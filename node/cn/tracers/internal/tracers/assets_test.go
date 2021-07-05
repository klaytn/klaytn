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

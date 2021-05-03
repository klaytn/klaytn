package params

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCaseWithBinaryCodeFormat = []struct {
	cf                            CodeFormat
	validateExpectVal             bool
	isDeployedIstanbulHFExpectVal bool
	stringExpectVal               string
}{
	{0b00000000, true, false, "CodeFormatEVM"},
	{0b10000000, true, true, "CodeFormatEVM"},
	{0b00000001, false, false, "UndefinedCodeFormat"},
	{0b10000001, false, true, "UndefinedCodeFormat"},
}

func TestString(t *testing.T) {
	for _, tc := range testCaseWithBinaryCodeFormat {
		assert.Equal(t, tc.stringExpectVal, tc.cf.String())
	}
}

func TestValidate(t *testing.T) {
	for _, tc := range testCaseWithBinaryCodeFormat {
		assert.Equal(t, tc.validateExpectVal, tc.cf.Validate())
	}
}

func TestAfterIstanbulHF(t *testing.T) {
	for _, tc := range testCaseWithBinaryCodeFormat {
		assert.Equal(t, tc.isDeployedIstanbulHFExpectVal, tc.cf.IsDeployedAfterIstanbulHF())
	}
}

func TestSetIstanbulHFField(t *testing.T) {
	testCaseSetIstanbulHFField := []struct {
		cf                        CodeFormat
		isDeployedAfterIstanbulHF bool
		expectCf                  CodeFormat
	}{
		{CodeFormatEVM, false, 0b00000000},
		{CodeFormatEVM, true, 0b10000000},
	}
	for _, tc := range testCaseSetIstanbulHFField {
		tc.cf.SetIstanbulHFField(tc.isDeployedAfterIstanbulHF)
		assert.Equal(t, tc.expectCf, tc.cf)
	}
}

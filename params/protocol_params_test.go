package params

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCaseWithBinaryCodeFormat = []struct {
	cf                     CodeInfo
	getCodeFormatExpectVal CodeFormat
	stringExpectVal        string
	validateExpectVal      bool
	getVmVersionExpectVal  VmVersion
}{
	{0b00000000, CodeFormatEVM, "CodeFormatEVM", true, VmVersionConstantinople},
	{0b00010000, CodeFormatEVM, "CodeFormatEVM", true, VmVersionIstanbul},
	{0b00000001, CodeFormatLast, "UndefinedCodeFormat", false, VmVersionConstantinople},
	{0b00010001, CodeFormatLast, "UndefinedCodeFormat", false, VmVersionIstanbul},
}

func TestGetCodeFormat(t *testing.T) {
	for _, tc := range testCaseWithBinaryCodeFormat {
		assert.Equal(t, tc.getCodeFormatExpectVal, tc.cf.GetCodeFormat())
	}
}

func TestString(t *testing.T) {
	for _, tc := range testCaseWithBinaryCodeFormat {
		assert.Equal(t, tc.stringExpectVal, tc.cf.GetCodeFormat().String())
	}
}

func TestValidate(t *testing.T) {
	for _, tc := range testCaseWithBinaryCodeFormat {
		assert.Equal(t, tc.validateExpectVal, tc.cf.GetCodeFormat().Validate())
	}
}

func TestGetVmVersion(t *testing.T) {
	for _, tc := range testCaseWithBinaryCodeFormat {
		assert.Equal(t, tc.getVmVersionExpectVal, tc.cf.GetVmVersion())
	}
}

func TestSetVmVersion(t *testing.T) {
	testCaseSetIstanbulHFField := []struct {
		cf                        CodeFormat
		isDeployedAfterIstanbulHF bool
		expectCi                  CodeInfo
	}{
		{CodeFormatEVM, false, 0b00000000},
		{CodeFormatEVM, true, 0b00010000},
	}
	for _, tc := range testCaseSetIstanbulHFField {
		assert.Equal(t, tc.expectCi, tc.cf.GenerateCodeInfo(Rules{IsIstanbul: tc.isDeployedAfterIstanbulHF}))
	}
}

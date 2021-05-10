package params

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCaseWithBinaryCodeFormat = []struct {
	cf                            CodeInfo
	getCodeFormatExpectVal        CodeFormat
	stringExpectVal               string
	validateExpectVal             bool
	isDeployedIstanbulHFExpectVal bool
	getVmVersionExpectVal         VmVersion
}{
	{0b00000000, CodeFormatEVM, "CodeFormatEVM", true, false, VmVersionConstantinople},
	{0b00010000, CodeFormatEVM, "CodeFormatEVM", true, true, VmVersionIstanbul},
	{0b00000001, CodeFormatLast, "UndefinedCodeFormat", false, false, VmVersionConstantinople},
	{0b00010001, CodeFormatLast, "UndefinedCodeFormat", false, true, VmVersionIstanbul},
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

func TestIsDeployedAfterHF(t *testing.T) {
	for _, tc := range testCaseWithBinaryCodeFormat {
		assert.Equal(t, tc.isDeployedIstanbulHFExpectVal, tc.cf.GetVmVersion().IsDeployedAfterHF(VmVersionIstanbul))
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

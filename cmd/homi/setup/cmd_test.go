package setup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveIfDefaultAddress(t *testing.T) {
	testcases := []struct {
		input    string
		expected string
	}{
		{
			`{
				"govParamContract" : "0x0000000000000000000000000000000000000000"
			}`,
			`{
			}`,
		},
		{
			`{
				"x": 1,
				"govParamContract" : "0x0000000000000000000000000000000000000000"
			}`,
			`{
				"x": 1
			}`,
		},
		{
			`{
				"x": 1,
				"govParamContract" : "0x0000000000000000000000000000000000000000",
				"y": 1
			}`,
			`{
				"x": 1,
				"y": 1
			}`,
		},
		{
			`{
				"govParamContract" : "0x0000000000000000000000000000000000000000",
				"y": 1
			}`,
			`{
				"y": 1
			}`,
		},
	}

	for i, tc := range testcases {
		actual := removeIfDefaultAddress([]byte(tc.input), "govParamContract")
		if actual == nil {
			actual = []byte("")
		}
		assert.Equal(t, []byte(tc.expected), actual, "testcases[%d] failed", i)
	}
}

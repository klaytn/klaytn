// Copyright 2018 The klaytn Authors
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

package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// ToCamelCase converts an under-score string to a camel-case string
func ToCamelCase(inputUnderScoreStr string) (camelCase string) {
	isToUpper := false

	for k, v := range inputUnderScoreStr {
		if k == 0 {
			camelCase = strings.ToUpper(string(inputUnderScoreStr[0]))
		} else {
			if isToUpper {
				camelCase += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					camelCase += string(v)
				}
			}
		}
	}
	return
}

// ToUnderScore converts a camel-case string to a under-score string
func ToUnderScore(s string) string {
	return SplitAndJoin(s, "_")
}

// ToHyphen converts a camel-case string to a hyphen-style string
func ToHyphen(s string) string {
	return SplitAndJoin(s, "-")
}

// SplitAndJoin converts a camel-case string to a string joined by the provided symbol
func SplitAndJoin(s string, symbol string) string {
	var camel = regexp.MustCompile("(^[^A-Z0-9]*)?([A-Z0-9]{2,}|[A-Z0-9][^A-Z]+|$)")
	var a []string
	for _, sub := range camel.FindAllStringSubmatch(s, -1) {
		if sub[1] != "" {
			a = append(a, sub[1])
		}
		if sub[2] != "" {
			a = append(a, sub[2])
		}
	}

	result := strings.ToLower(strings.Join(a, symbol))
	result = strings.TrimPrefix(result, "_")
	result = strings.TrimSuffix(result, "_")
	return result
}

func FormatPackage(name string) string {
	if name == "" {
		return ""
	}
	return fmt.Sprintf("%v.", name)
}

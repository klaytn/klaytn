// Modifications Copyright 2019 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from cmd/geth/usage.go (2018/06/04).
// Modified and improved for the klaytn development.

// Contains the Klaytn node command usage template and generator.

package utils

import (
	"io"
	"strings"

	"gopkg.in/urfave/cli.v1"
)

func NewHelpPrinter(fg []FlagGroup) func(w io.Writer, tmp string, data interface{}) {
	originalHelpPrinter := cli.HelpPrinter
	return func(w io.Writer, tmpl string, data interface{}) {
		type helpData struct {
			App        interface{}
			FlagGroups []FlagGroup
		}

		if tmpl == GlobalAppHelpTemplate {
			// Iterate over all the flags and add any uncategorized ones
			categorized := make(map[string]struct{})
			for _, group := range fg {
				for _, flag := range group.Flags {
					categorized[flag.String()] = struct{}{}
				}
			}
			var uncategorized []cli.Flag
			for _, flag := range data.(*cli.App).Flags {
				if _, ok := categorized[flag.String()]; !ok {
					if strings.HasPrefix(flag.GetName(), "dashboard") {
						continue
					}
					uncategorized = append(uncategorized, flag)
				}
			}
			if len(uncategorized) > 0 {
				// Append all ungategorized options to the misc group
				miscs := len(fg[len(fg)-1].Flags)
				fg[len(fg)-1].Flags = append(fg[len(fg)-1].Flags, uncategorized...)

				// Make sure they are removed afterwards
				defer func() {
					fg[len(fg)-1].Flags = fg[len(fg)-1].Flags[:miscs]
				}()
			}
			// Render out custom usage screen
			originalHelpPrinter(w, tmpl, helpData{data, fg})
		} else if tmpl == CommandHelpTemplate {
			// Iterate over all command specific flags and categorize them
			categorized := make(map[string][]cli.Flag)
			for _, flag := range data.(cli.Command).Flags {
				if _, ok := categorized[flag.String()]; !ok {
					categorized[flagCategory(flag, fg)] = append(categorized[flagCategory(flag, fg)], flag)
				}
			}

			// sort to get a stable ordering
			sorted := make([]FlagGroup, 0, len(categorized))
			for cat, flgs := range categorized {
				sorted = append(sorted, FlagGroup{cat, flgs})
			}

			// add sorted array to data and render with default printer
			originalHelpPrinter(w, tmpl, map[string]interface{}{
				"cmd":              data,
				"categorizedFlags": sorted,
			})
		} else {
			originalHelpPrinter(w, tmpl, data)
		}
	}
}

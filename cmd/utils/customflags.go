// Modifications Copyright 2018 The klaytn Authors
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
// This file is derived from cmd/utils/customflags.go (2018/06/04).
// Modified and improved for the klaytn development.

package utils

import (
	"encoding"
	"errors"
	"flag"
	"fmt"
	"math/big"
	"os"
	"os/user"
	"path"
	"strings"
	"syscall"

	"github.com/klaytn/klaytn/common/math"
	"github.com/klaytn/klaytn/datasync/downloader"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

// NOTE-klaytn: The custom directoryFlag deprecated.
// urfave.v2 new flag, PathFlag, replaced the directoryFlag.

// Custom type which is registered in the flags library which cli uses for
// argument parsing. This allows us to expand Value to an absolute path when
// the argument is parsed
type DirectoryString struct {
	Value string
}

func (self *DirectoryString) String() string {
	return self.Value
}

func (self *DirectoryString) Set(value string) error {
	self.Value = expandPath(value)
	return nil
}

type WrappedDirectoryFlag struct {
	DirectoryFlag
	set *flag.FlagSet
}

func NewWrappedDirectoryFlag(fl DirectoryFlag) *WrappedDirectoryFlag {
	return &WrappedDirectoryFlag{DirectoryFlag: fl, set: nil}
}

func (f *WrappedDirectoryFlag) Apply(set *flag.FlagSet) {
	f.set = set
	f.DirectoryFlag.Apply(set)
}

func (f *WrappedDirectoryFlag) ApplyInputSourceValue(context *cli.Context, isc altsrc.InputSourceContext) error {
	if f.set != nil {
		if !isEnvVarSet(f.EnvVar) {
			value, err := isc.String(f.DirectoryFlag.Name)
			if err != nil {
				return err
			}
			if value != "" {
				eachName(f.Name, func(name string) {
					f.set.Set(f.Name, value)
				})
			}
		}
	}
	return nil
}

// Custom cli.Flag type which expand the received string to an absolute path.
// e.g. ~/.ethereum -> /home/username/.ethereum
type DirectoryFlag struct {
	Name   string
	Value  DirectoryString
	Usage  string
	EnvVar string
}

func (self DirectoryFlag) String() string {
	fmtString := "%s %v\t%v"
	if len(self.Value.Value) > 0 {
		fmtString = "%s \"%v\"\t%v"
	}
	return fmt.Sprintf(fmtString, prefixedNames(self.Name), self.Value.Value, self.Usage)
}

// called by cli library, grabs variable from environment (if in env)
// and adds variable to flag set for parsing.
func (self DirectoryFlag) Apply(set *flag.FlagSet) {
	if self.EnvVar != "" {
		if envVal, ok := syscall.Getenv(self.EnvVar); ok {
			self.Value.Value = envVal
		}
	}
	eachName(self.Name, func(name string) {
		set.Var(&self.Value, self.Name, self.Usage)
	})
}

func eachName(longName string, fn func(string)) {
	parts := strings.Split(longName, ",")
	for _, name := range parts {
		name = strings.Trim(name, " ")
		fn(name)
	}
}

func isEnvVarSet(envVars string) bool {
	for _, envVar := range strings.Split(envVars, ",") {
		envVar = strings.TrimSpace(envVar)
		if env, ok := syscall.Getenv(envVar); ok {
			// TODO: Can't use this for bools as
			// set means that it was true or false based on
			// Bool flag type, should work for other types
			logger.Info("env", "env", env)
			return true
		}
	}
	return false
}

type TextMarshaler interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
}

// textMarshalerVal turns a TextMarshaler into a flag.Value
type textMarshalerVal struct {
	v TextMarshaler
}

func (v textMarshalerVal) String() string {
	if v.v == nil {
		return ""
	}
	text, _ := v.v.MarshalText()
	return string(text)
}

func (v textMarshalerVal) Set(s string) error {
	return v.v.UnmarshalText([]byte(s))
}

// TextMarshalerFlag wraps a TextMarshaler value.
type TextMarshalerFlag struct {
	Name string

	Category string
	Usage    string

	Required   bool
	Hidden     bool
	HasBeenSet bool

	Value       TextMarshaler
	Destination *TextMarshaler

	Aliases []string
	EnvVars []string

	Action func(*cli.Context, TextMarshaler) error
}

// IsSet returns whether or not the flag has been set through env or file
func (f *TextMarshalerFlag) IsSet() bool {
	return f.HasBeenSet
}

// Names returns the names of the flag
func (f *TextMarshalerFlag) Names() []string {
	return cli.FlagNames(f.Name, f.Aliases)
}

// IsRequired returns whether or not the flag is required
func (f *TextMarshalerFlag) IsRequired() bool {
	return f.Required
}

// IsVisible returns true if the flag is not hidden, otherwise false
func (f *TextMarshalerFlag) IsVisible() bool {
	return !f.Hidden
}

func (f *TextMarshalerFlag) String() string {
	return fmt.Sprintf("%s \"%v\"\t%v", prefixedNames(f.Name), f.Value, f.Usage)
}

// TakesValue returns true of the flag takes a value, otherwise false
func (f *TextMarshalerFlag) TakesValue() bool {
	return true
}

// GetUsage returns the usage string for the flag
func (f *TextMarshalerFlag) GetUsage() string {
	return f.Usage
}

// GetCategory returns the category for the flag
func (f *TextMarshalerFlag) GetCategory() string {
	return f.Category
}

// GetValue returns the flags value as string representation and an empty
// string if the flag takes no value at all.
func (f *TextMarshalerFlag) GetValue() TextMarshaler {
	return f.Value
}

// GetEnvVars returns the env vars for this flag
func (f *TextMarshalerFlag) GetEnvVars() []string {
	return f.EnvVars
}

func (f *TextMarshalerFlag) Apply(set *flag.FlagSet) error {
	if f.EnvVars[0] != "" && f.Value != nil {
		if envVal, ok := syscall.Getenv(f.EnvVars[0]); ok {
			var mode downloader.SyncMode
			switch envVal {
			case "full":
				mode = downloader.FullSync
			case "fast":
				mode = downloader.FastSync
			case "snap":
				mode = downloader.SnapSync
			case "light":
				mode = downloader.LightSync
			}
			f.Value = &mode
		}
	}
	eachName(f.Name, func(name string) {
		set.Var(textMarshalerVal{f.Value}, f.Name, f.Usage)
	})

	return nil
}

// Get returns the flag’s value in the given Context.
func (f *TextMarshalerFlag) Get(ctx *cli.Context) string {
	return ctx.Path(f.Name)
}

// RunAction executes flag action if set
func (f *TextMarshalerFlag) RunAction(c *cli.Context) error {
	if f.Action != nil {
		return f.Action(c, GlobalTextMarshaler(c, f.Name))
	}
	return nil
}

type WrappedTextMarshalerFlag struct {
	*TextMarshalerFlag
	set *flag.FlagSet
}

func NewWrappedTextMarshalerFlag(fl *TextMarshalerFlag) *WrappedTextMarshalerFlag {
	return &WrappedTextMarshalerFlag{TextMarshalerFlag: fl, set: nil}
}

func (f *WrappedTextMarshalerFlag) Apply(set *flag.FlagSet) error {
	f.set = set
	return f.TextMarshalerFlag.Apply(set)
}

func (f *WrappedTextMarshalerFlag) ApplyInputSourceValue(context *cli.Context, isc altsrc.InputSourceContext) error {
	if f.set != nil {
		if !context.IsSet(f.Name) && !isEnvVarSet(f.EnvVars[0]) {
			value, err := isc.String(f.TextMarshalerFlag.Name)
			if err != nil {
				return err
			}
			if value != "" {
				eachName(f.Name, func(name string) {
					f.set.Set(f.Name, value)
				})
			}
		}
	}
	return nil
}

// GlobalTextMarshaler returns the value of a TextMarshalerFlag from the global flag set.
func GlobalTextMarshaler(ctx *cli.Context, name string) TextMarshaler {
	val := ctx.Generic(name)
	if val == nil {
		return nil
	}
	return val.(textMarshalerVal).v
}

// BigFlag is a command line flag that accepts 256 bit big integers in decimal or
// hexadecimal syntax.
type BigFlag struct {
	Name  string
	Value *big.Int
	Usage string
}

// bigValue turns *big.Int into a flag.Value
type bigValue big.Int

func (b *bigValue) String() string {
	if b == nil {
		return ""
	}
	return (*big.Int)(b).String()
}

func (b *bigValue) Set(s string) error {
	int, ok := math.ParseBig256(s)
	if !ok {
		return errors.New("invalid integer syntax")
	}
	*b = (bigValue)(*int)
	return nil
}

func (f BigFlag) GetName() string {
	return f.Name
}

func (f BigFlag) String() string {
	fmtString := "%s %v\t%v"
	if f.Value != nil {
		fmtString = "%s \"%v\"\t%v"
	}
	return fmt.Sprintf(fmtString, prefixedNames(f.Name), f.Value, f.Usage)
}

func (f BigFlag) Apply(set *flag.FlagSet) {
	eachName(f.Name, func(name string) {
		set.Var((*bigValue)(f.Value), f.Name, f.Usage)
	})
}

// GlobalBig returns the value of a BigFlag from the global flag set.
func GlobalBig(ctx *cli.Context, name string) *big.Int {
	val := ctx.Generic(name)
	if val == nil {
		return nil
	}
	return (*big.Int)(val.(*bigValue))
}

func prefixFor(name string) (prefix string) {
	if len(name) == 1 {
		prefix = "-"
	} else {
		prefix = "--"
	}

	return
}

func prefixedNames(fullName string) (prefixed string) {
	parts := strings.Split(fullName, ",")
	for i, name := range parts {
		name = strings.Trim(name, " ")
		prefixed += prefixFor(name) + name
		if i < len(parts)-1 {
			prefixed += ", "
		}
	}
	return
}

// Expands a file path
// 1. replace tilde with users home dir
// 2. expands embedded environment variables
// 3. cleans the path, e.g. /a/b/../c -> /a/c
// Note, it has limitations, e.g. ~someuser/tmp will not be expanded
func expandPath(p string) string {
	if strings.HasPrefix(p, "~/") || strings.HasPrefix(p, "~\\") {
		if home := homeDir(); home != "" {
			p = home + p[1:]
		}
	}
	return path.Clean(os.ExpandEnv(p))
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

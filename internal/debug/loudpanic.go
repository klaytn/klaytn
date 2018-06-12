// +build go1.6

package debug

import "runtime/debug"

// LoudPanic panics in a way that gets all goroutine stacks printed on stderr.
func LoudPanic(x interface{}) {
	debug.SetTraceback("all")
	panic(x)
}
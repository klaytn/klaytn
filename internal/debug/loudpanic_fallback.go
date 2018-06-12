// +build !go1.6

package debug

// LoudPanic panics in a way that gets all goroutine stacks printed on stderr.
func LoudPanic(x interface{}) {
	panic(x)
}

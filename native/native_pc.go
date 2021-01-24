//+build windows linux darwin

package native

// Invoke invokes the specified feature's method for the current platform.
func Invoke(feature Feature) {
	feature.OnPc()
}
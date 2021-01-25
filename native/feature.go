// This package provides APIs for working with platform-specific code.
package native

// Feature is used to implement platform-specific logic.
// This interface provides methods to be called for each supported platform.
type Feature interface {
	// This method is called if the app is running in the browser.
	OnWeb()

	// This method is called if the app is running on either of the following platforms: Windows, Linux, Darwin(macOS).
	OnPc()
}

type FeatureImpl struct {

}

func (f *FeatureImpl) OnWeb() {

}

func (f *FeatureImpl) OnPc() {

}
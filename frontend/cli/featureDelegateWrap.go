package cli

//NativeFeatureDelegate is a redeclaration of frontend.FeatureDelegate to be exported into the CLI.
//Can be used to implement a feature in native code.
//The OnInvoke method takes in and returns string instead of interface{} (unlike the original interface)
//due to the gobind type restrictions.
type NativeFeatureDelegate interface {
	OnInvoke(string) string
}

//NativeFeatureDelegateWrap is a wrap around NativeFeatureDelegate implementing the frontend.FeatureDelegate,
//so native feature can be passed to the engine.
type NativeFeatureDelegateWrap struct {
	delegate NativeFeatureDelegate
}

//OnInvoke implements the frontend.FeatureDelegate interface.
func (f *NativeFeatureDelegateWrap) OnInvoke(data interface{}) interface{} {
	if s, ok := data.(string); ok {
		return f.delegate.OnInvoke(s)
	} else {
		return f.delegate.OnInvoke("")
	}
}

//NewNativeFeatureDelegateWrap creates a new instance of the NativeFeatureDelegateWrap with the specified NativeFeatureDelegate.
func NewNativeFeatureDelegateWrap(delegate NativeFeatureDelegate) *NativeFeatureDelegateWrap {
	return &NativeFeatureDelegateWrap{
		delegate: delegate,
	}
}
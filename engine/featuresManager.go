package engine

type FeatureCode int

const (
	FeatureUnknown = iota
	FeatureClipboardManager
)

//FeaturesManager is responsible for keeping track of what frontend features are available and invoking them.
//Frontend should register all features that it provides. Then those features become available and can be invoked.
type FeaturesManager struct {
	features map[FeatureCode]interface{}
}

//IsFeatureAvailable checks if feature with the specified FeatureCode is available (provided by the frontend).
func (f *FeaturesManager) IsFeatureAvailable(code FeatureCode) bool {
	_, ok := f.features[code]
	return ok
}

//RegisterFeatureDelegate should be called by the frontend.
func (f *FeaturesManager) RegisterFeatureDelegate(code FeatureCode, delegate interface{}) {
	f.features[code] = delegate
}

//GetFeature returns the feature with the given code.
//Returns nil if feature is not available.
func (f *FeaturesManager) GetFeature(code FeatureCode) interface{} {
	return f.features[code]
}

func newFeaturesManager() *FeaturesManager {
	return &FeaturesManager{
		features: map[FeatureCode]interface{}{},
	}
}
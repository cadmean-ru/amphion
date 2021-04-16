package engine

import "github.com/cadmean-ru/amphion/frontend"

//FeaturesManager is responsible for keeping track of what frontend features are available and invoking them.
//Frontend should register all features that it provides. Then those features become available and can be invoked.
type FeaturesManager struct {
	features map[FeatureCode]frontend.FeatureDelegate
}

//InvokeFeature invokes the feature with the specified FeatureCode if is available.
//Passes the specified data to the feature delegate and returns the result.
func (f *FeaturesManager) InvokeFeature(code FeatureCode, data interface{}) interface{} {
	if feat, ok := f.features[code]; ok {
		return feat.OnInvoke(data)
	}
	return nil
}

//IsFeatureAvailable checks if feature with the specified FeatureCode is available (provided by the frontend).
func (f *FeaturesManager) IsFeatureAvailable(code FeatureCode) bool {
	_, ok := f.features[code]
	return ok
}

//RegisterFeatureDelegate should be called by the frontend.
func (f *FeaturesManager) RegisterFeatureDelegate(code FeatureCode, delegate frontend.FeatureDelegate) {
	f.features[code] = delegate
}

func newFeaturesManager() *FeaturesManager {
	return &FeaturesManager{features: map[FeatureCode]frontend.FeatureDelegate{}}
}
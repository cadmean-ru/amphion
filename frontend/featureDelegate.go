package frontend

//FeatureDelegate defines interface for a feature, that is provided by the frontend, to be used by the engine.
//The OnInvoke method is called by the engine. It takes in and returns interface{},
//so any types can be used for the specific feature.
type FeatureDelegate interface {
	OnInvoke(data interface{}) interface{}
}

//FeatureFunc is an implementation of the FeatureDelegate interface that takes in a function that implements the desired feature.
type FeatureFunc struct {
	delegate func(interface{}) interface{}
}

//OnInvoke implementation of the interface
func (f *FeatureFunc) OnInvoke(data interface{}) interface{} {
	return f.delegate(data)
}

//NewFeatureFunc Creates a new instance of FeatureFunc with the specified delegate function.
func NewFeatureFunc(f func(interface{}) interface{}) *FeatureFunc {
	return &FeatureFunc{
		delegate: f,
	}
}


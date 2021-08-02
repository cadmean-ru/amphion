package engine

type AppDelegate interface {
	OnAppLoaded()
	OnAppStopping()
}

type AppDelegateImpl struct {

}

func (d *AppDelegateImpl) OnAppLoaded() {

}

func (d *AppDelegateImpl) OnAppStopping() {

}
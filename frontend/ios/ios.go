package ios

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/frontend"
	"github.com/cadmean-ru/amphion/frontend/cli"
	"github.com/cadmean-ru/amphion/rendering"
	"gopkg.in/yaml.v3"
)

type Frontend struct {
	f cli.FrontendCLI
	handler frontend.CallbackHandler
	resMan  frontend.ResourceManager
}

func (f *Frontend) Init() {
	f.f.Init()
}

func (f *Frontend) Run() {
	f.f.Run()
}

func (f *Frontend) Reset() {
	f.f.Reset()
}

func (f *Frontend) SetCallback(handler frontend.CallbackHandler) {
	f.handler = handler
}

func (f *Frontend) GetInputManager() frontend.InputManager {
	return &InputManager{}
}

func (f *Frontend) GetRenderer() rendering.Renderer {
	return &MetalRenderer{}
}

func (f *Frontend) GetContext() frontend.Context {
	return frontend.Context{}
}

func (f *Frontend) GetPlatform() common.Platform {
	return common.PlatformFromString("ios")
}

func (f *Frontend) CommencePanic(reason, msg string) {
	f.f.CommencePanic(reason, msg)
}

func (f *Frontend) ReceiveMessage(message frontend.Message) {

}

func (f *Frontend) GetResourceManager() frontend.ResourceManager {
	return f.resMan
}

func (f *Frontend) GetApp() *frontend.App {
	app := &frontend.App{}

	err := yaml.Unmarshal(f.f.GetAppData(), &app)
	if err != nil {
		f.f.CommencePanic("Load app failed", err.Error())
		return app
	}

	return app
}

func (f *Frontend) GetLaunchArgs() a.SiMap {
	return a.SiMap{}
}

func NewFrontend(f cli.FrontendCLI, rm cli.ResourceManagerCLI) frontend.Frontend {
	return &Frontend{
		f: f,
		resMan: &ResourceManager{
			ResourceManagerImpl: frontend.NewResourceManagerImpl(),
			rm: rm,
		},
	}
}

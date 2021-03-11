package ios

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/frontend"
	"github.com/cadmean-ru/amphion/frontend/cli"
	"github.com/cadmean-ru/amphion/rendering"
	"gopkg.in/yaml.v3"
)

type Frontend struct {
	f        cli.FrontendCLI
	handler  frontend.CallbackHandler
	resMan   frontend.ResourceManager
	renderer *MetalRenderer
	msgChan  chan frontend.Message
}

func (f *Frontend) Init() {
	f.f.Init()

	f.renderer.Prepare()
}

func (f *Frontend) Run() {
	f.f.Run()

	fmt.Println("Running")

	for msg := range f.msgChan {
		switch msg.Code {
		case frontend.MessageRender:
			fmt.Println("Message render")
			f.renderer.PerformRendering()
		case frontend.MessageExec:
			if msg.Data != nil {
				if action, ok := msg.Data.(func()); ok {
					action()
				}
			}
		case frontend.MessageExit:
			close(f.msgChan)
			return
		}
	}
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
	return f.renderer
}

func (f *Frontend) GetContext() frontend.Context {
	ctx := frontend.Context{}

	c := f.f.GetContext()

	ctx.ScreenInfo = common.NewScreenInfo(int(c.ScreenSize.X), int(c.ScreenSize.Y))
	fmt.Printf("Got context: %+v\n", ctx)

	return ctx
}

func (f *Frontend) GetPlatform() common.Platform {
	return common.PlatformFromString("ios")
}

func (f *Frontend) CommencePanic(reason, msg string) {
	f.f.CommencePanic(reason, msg)
}

func (f *Frontend) ReceiveMessage(message frontend.Message) {
	f.msgChan <- message
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

func NewFrontend(f cli.FrontendCLI, rm cli.ResourceManagerCLI, r cli.RendererCLI) frontend.Frontend {
	return &Frontend{
		f: f,
		resMan: &ResourceManager{
			ResourceManagerImpl: frontend.NewResourceManagerImpl(),
			rm:                  rm,
		},
		renderer: &MetalRenderer{r: r},
		msgChan:  make(chan frontend.Message, 10),
	}
}

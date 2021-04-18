//+build android

package android

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/frontend"
	"github.com/cadmean-ru/amphion/frontend/cli"
	"github.com/cadmean-ru/amphion/rendering"
	"gopkg.in/yaml.v2"
)

type Frontend struct {
	frontendDelegate cli.FrontendDelegate
	handler          *cli.CallbackHandler
	resMan           frontend.ResourceManager
	renderer         *rendering.ARenderer
	msgChan          chan frontend.Message
}

func (f *Frontend) Init() {
	f.frontendDelegate.Init()

	f.renderer.SetManagementMode(rendering.FrontendManaged)
}

func (f *Frontend) Run() {
	f.frontendDelegate.Run()

	for msg := range f.msgChan {
		switch msg.Code {
		case frontend.MessageRender:
			f.frontendDelegate.ExecuteOnRenderingThread(cli.NewExecDelegate(f.renderer.PerformRendering))
		case frontend.MessageExec:
			if msg.Data != nil {
				if action, ok := msg.Data.(func()); ok {
					f.frontendDelegate.ExecuteOnMainThread(cli.NewExecDelegate(action))
				}
			}
		case frontend.MessageExit:
			close(f.msgChan)
			return
		}
	}
}

func (f *Frontend) Reset() {
	f.frontendDelegate.Reset()
}

func (f *Frontend) SetCallback(handler frontend.CallbackHandler) {
	f.handler = cli.NewCallbackHandler(handler)
	f.frontendDelegate.SetCallback(f.handler)
}

func (f *Frontend) GetInputManager() frontend.InputManager {
	return nil
}

func (f *Frontend) GetRenderer() *rendering.ARenderer {
	return f.renderer
}

func (f *Frontend) GetContext() frontend.Context {
	ctx := frontend.Context{}

	c := f.frontendDelegate.GetContext()

	ctx.ScreenInfo = common.NewScreenInfo(int(c.ScreenSize.X), int(c.ScreenSize.Y))
	fmt.Printf("Got context: %+v\n", ctx)

	return ctx
}

func (f *Frontend) GetPlatform() common.Platform {
	return common.PlatformFromString("android")
}

func (f *Frontend) CommencePanic(reason, msg string) {
	f.frontendDelegate.CommencePanic(reason, msg)
}

func (f *Frontend) ReceiveMessage(message frontend.Message) {
	f.msgChan <- message
}

func (f *Frontend) GetResourceManager() frontend.ResourceManager {
	return f.resMan
}

func (f *Frontend) GetApp() *frontend.App {
	app := &frontend.App{}

	err := yaml.Unmarshal(f.frontendDelegate.GetAppData(), &app)
	if err != nil {
		f.frontendDelegate.CommencePanic("Load app failed", err.Error())
		return app
	}

	return app
}

func (f *Frontend) GetLaunchArgs() a.SiMap {
	return a.SiMap{}
}

func NewFrontend(f cli.FrontendDelegate, rm cli.ResourceManagerDelegate, rd cli.RendererDelegate) frontend.Frontend {
	return &Frontend{
		frontendDelegate: f,
		resMan:           cli.NewResourceManagerImpl(rm),
		renderer:         rendering.NewARenderer(rd),
		msgChan:          make(chan frontend.Message, 10),
	}
}

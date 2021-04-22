//+build android

package android

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/dispatch"
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
	msgChan          *dispatch.MessageQueue
	mainDispatcher   dispatch.WorkDispatcher
}

func (f *Frontend) Init() {
	f.frontendDelegate.Init()

	f.renderer.SetManagementMode(rendering.FrontendManaged)
}

func (f *Frontend) Run() {
	f.frontendDelegate.Run()

	for {
		msg := f.msgChan.DequeueBlocking()

		if msg.What == frontend.MessageExit {
			break
		}

		switch msg.What {
		case frontend.MessageExec, dispatch.MessageWorkExec:
			if msg.AnyData != nil {
				if action, ok := msg.AnyData.(dispatch.WorkItem); ok {
					f.mainDispatcher.Execute(action)
				}
			}
		}
	}

	f.msgChan.Close()
}

func (f *Frontend) SendMessage(message *dispatch.Message) {
	f.msgChan.Enqueue(message)
}

func (f *Frontend) SetEngineDispatcher(disp dispatch.MessageDispatcher) {
	f.frontendDelegate.SetCallbackDispatcher(disp)
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

func (f *Frontend) GetMessageDispatcher() dispatch.MessageDispatcher {
	return f
}

func (f *Frontend) GetWorkDispatcher() dispatch.WorkDispatcher {
	return f.mainDispatcher
}


func NewFrontend(f cli.FrontendDelegate, rm cli.ResourceManagerDelegate, rd cli.RendererDelegate) frontend.Frontend {
	return &Frontend{
		frontendDelegate: f,
		resMan:           cli.NewResourceManagerImpl(rm),
		renderer:         rendering.NewARenderer(rd, f.GetRenderingThreadDispatcher()),
		msgChan:          dispatch.NewMessageQueue(1000),
	}
}

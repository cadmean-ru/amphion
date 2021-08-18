//+build ios

package ios

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/dispatch"
	"github.com/cadmean-ru/amphion/frontend"
	"github.com/cadmean-ru/amphion/frontend/cli"
	"github.com/cadmean-ru/amphion/rendering"
	"gopkg.in/yaml.v3"
)

type Frontend struct {
	*dispatch.LooperImpl
	f        cli.FrontendDelegate
	resMan   frontend.ResourceManager
	renderer *rendering.ARenderer
	msgQueue *dispatch.MessageQueue
	mainDisp dispatch.WorkDispatcher
}

func (f *Frontend) Init() {
	f.f.Init()

	f.renderer.Prepare()
}

func (f *Frontend) Run() {
	f.f.Run()

	fmt.Println("Running")

	var execHandler = func(msg *dispatch.Message) {
		if msg.AnyData != nil {
			if action, ok := msg.AnyData.(dispatch.WorkItem); ok {
				f.mainDisp.Execute(action)
			}
		}
	}

	f.SetMessageHandler(frontend.MessageExec, execHandler)
	f.SetMessageHandler(dispatch.MessageWorkExec, execHandler)

	f.Loop()
}

func (f *Frontend) GetRenderer() *rendering.ARenderer {
	return f.renderer
}

func (f *Frontend) GetContext() frontend.Context {
	ctx := frontend.Context{}

	c := f.f.GetContext()

	ctx.ScreenInfo = common.NewScreenInfo(int(c.ScreenSize.X), int(c.ScreenSize.Y))
	fmt.Printf("Got context: %+v\n", ctx)

	return ctx
}

func (f *Frontend) SendMessage(message *dispatch.Message) {
	f.msgQueue.Enqueue(message)
}

func (f *Frontend) SetEngineDispatcher(disp dispatch.MessageDispatcher) {
	f.f.SetCallbackDispatcher(disp)
}

func (f *Frontend) GetMessageDispatcher() dispatch.MessageDispatcher {
	return f
}

func (f *Frontend) GetWorkDispatcher() dispatch.WorkDispatcher {
	return f.mainDisp
}

func (f *Frontend) GetPlatform() common.Platform {
	return common.PlatformFromString("ios")
}

func (f *Frontend) CommencePanic(reason, msg string) {
	f.f.CommencePanic(reason, msg)
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

func NewFrontend(f cli.FrontendDelegate, rm cli.ResourceManagerDelegate, rd cli.RendererDelegate) frontend.Frontend {
	fmt.Println("new frontend")
	return &Frontend{
		f:          f,
		resMan:     cli.NewResourceManagerImpl(rm),
		renderer:   rendering.NewARenderer(rd, f.GetRenderingThreadDispatcher()),
		msgQueue:   dispatch.NewMessageQueue(1000),
		mainDisp:   f.GetMainThreadDispatcher(),
		LooperImpl: dispatch.NewLooperImpl(100),
	}
}

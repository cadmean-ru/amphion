package atest

import (
	"errors"
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/dispatch"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/frontend"
	"github.com/cadmean-ru/amphion/rendering"
	"time"
)

var instance *Frontend

// The testing frontend.
type Frontend struct {
	disp dispatch.MessageDispatcher
	clicked bool
	resMan  *ResourceManager
	renderer *rendering.ARenderer
}

func (f *Frontend) Run() {
	f.disp.SendMessage(dispatch.NewMessage(frontend.CallbackReady))
}

func (f *Frontend) GetRenderer() *rendering.ARenderer {
	return f.renderer
}

func (f *Frontend) GetContext() frontend.Context {
	return frontend.Context{
		ScreenInfo: common.NewScreenInfo(500, 500),
		Domain:     "cadmean.ru",
	}
}

func (f *Frontend) CommencePanic(reason, msg string) {
	panic(fmt.Sprintf("Test failed: %s - %s", reason, msg))
}

func (f *Frontend) GetMessageDispatcher() dispatch.MessageDispatcher {
	return f
}

func (f *Frontend) GetWorkDispatcher() dispatch.WorkDispatcher {
	return f
}

func (f *Frontend) GetResourceManager() frontend.ResourceManager {
	return f.resMan
}

func (f *Frontend) GetApp() *frontend.App {
	return nil
}

func (f *Frontend) GetLaunchArgs() a.SiMap {
	return a.SiMap{}
}

func (f *Frontend) Init() {
	instance = f
	f.renderer.Prepare()
	fmt.Println("This is a testing frontend")
}

func (f *Frontend) GetPlatform() common.Platform {
	return common.PlatformFromString("testing")
}

func (f *Frontend) SetEngineDispatcher(disp dispatch.MessageDispatcher) {
	f.disp = disp
}

func (f *Frontend) Execute(item dispatch.WorkItem) {
	item.Execute()
}

func (f *Frontend) SendMessage(message *dispatch.Message) {
	fmt.Printf("Testing frontend received message: %+v\n", message)
}

func (f *Frontend) simulateCallback(code int, data string) {
	callback := dispatch.NewMessageWithStringData(code, data)
	f.disp.SendMessage(callback)
}

func (f *Frontend) simulateClick(x, y int, button engine.MouseButton) {
	if f.clicked {
		return
	}

	data := fmt.Sprintf("%d;%d;%d", x, y, button)
	f.simulateCallback(frontend.CallbackMouseDown, data)
	f.clicked = true
	time.Sleep(100)
	f.simulateCallback(frontend.CallbackMouseUp, data)
	f.clicked = false
}

func newTestingFrontend() *Frontend {
	f := &Frontend{
		resMan: newResourceManager(),
	}
    f.renderer = rendering.NewARenderer(&RendererDelegate{}, f)
    return f
}

type ResourceManager struct {
	resources map[a.ResId]string
	idgen     *common.IdGenerator
}

func (r *ResourceManager) RegisterResource(path string) {
	r.resources[a.ResId(r.idgen.NextId())] = path
}

func (r *ResourceManager) IdOf(path string) a.ResId {
	for id, p := range r.resources {
		if p == path {
			return id
		}
	}

	return -1
}

func (r *ResourceManager) PathOf(id a.ResId) string {
	return r.resources[id]
}

func (r *ResourceManager) FullPathOf(id a.ResId) string {
	return "res/" + r.resources[id]
}

func (r *ResourceManager) ReadFile(id a.ResId) ([]byte, error) {
	return nil, errors.New("file reading not yet supported on testing frontend")
}

func newResourceManager() *ResourceManager {
	return &ResourceManager{
		resources: make(map[a.ResId]string),
		idgen:     common.NewIdGenerator(),
	}
}

type RendererDelegate struct {

}

func (r *RendererDelegate) OnPrepare() {

}

func (r *RendererDelegate) OnPerformRenderingStart() {

}

func (r *RendererDelegate) OnPerformRenderingEnd() {

}

func (r *RendererDelegate) OnClear() {

}

func (r *RendererDelegate) OnStop() {

}
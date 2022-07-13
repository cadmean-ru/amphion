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

// HeadlessFrontend the testing frontend.
type HeadlessFrontend struct {
	disp     dispatch.MessageDispatcher
	clicked  bool
	resMan   *HeadlessResourceManager
	renderer *rendering.ARenderer
}

func (f *HeadlessFrontend) Run() {
	f.disp.SendMessage(dispatch.NewMessage(frontend.CallbackReady))
}

func (f *HeadlessFrontend) GetRenderer() *rendering.ARenderer {
	return f.renderer
}

func (f *HeadlessFrontend) GetContext() frontend.Context {
	return frontend.Context{
		ScreenInfo: common.NewScreenInfo(500, 500),
		Domain:     "cadmean.ru",
	}
}

func (f *HeadlessFrontend) CommencePanic(reason, msg string) {
	panic(fmt.Sprintf("Test failed: %s - %s", reason, msg))
}

func (f *HeadlessFrontend) GetMessageDispatcher() dispatch.MessageDispatcher {
	return f
}

func (f *HeadlessFrontend) GetWorkDispatcher() dispatch.WorkDispatcher {
	return f
}

func (f *HeadlessFrontend) GetResourceManager() frontend.ResourceManager {
	return f.resMan
}

func (f *HeadlessFrontend) GetApp() *frontend.App {
	return nil
}

func (f *HeadlessFrontend) GetLaunchArgs() a.SiMap {
	return a.SiMap{}
}

func (f *HeadlessFrontend) Init() {
	f.renderer.Prepare()
	fmt.Println("This is a testing frontend")
}

func (f *HeadlessFrontend) GetPlatform() common.Platform {
	return common.PlatformFromString("testing")
}

func (f *HeadlessFrontend) SetEngineDispatcher(disp dispatch.MessageDispatcher) {
	f.disp = disp
}

func (f *HeadlessFrontend) Execute(item dispatch.WorkItem) {
	item.Execute()
}

func (f *HeadlessFrontend) SendMessage(message *dispatch.Message) {
	fmt.Printf("Testing frontend received message: %+v\n", message)
}

func (f *HeadlessFrontend) SimulateCallback(message *dispatch.Message) {
	f.disp.SendMessage(message)
}

func (f *HeadlessFrontend) simulateClick(x, y int, button engine.MouseButton) {
	if f.clicked {
		return
	}

	data := fmt.Sprintf("%d;%d;%d", x, y, button)
	f.SimulateCallback(dispatch.NewMessageWithStringData(frontend.CallbackMouseDown, data))
	f.clicked = true
	time.Sleep(100)
	f.SimulateCallback(dispatch.NewMessageWithStringData(frontend.CallbackMouseUp, data))
	f.clicked = false
}

func NewHeadlessFrontend() *HeadlessFrontend {
	f := &HeadlessFrontend{
		resMan: newResourceManager(),
	}
	f.renderer = rendering.NewARenderer(&HeadlessRendererDelegate{}, f)
	return f
}

type HeadlessResourceManager struct {
	resources map[a.ResId]string
	idgen     *common.IdGenerator
}

func (r *HeadlessResourceManager) RegisterResource(path string) {
	r.resources[a.ResId(r.idgen.NextId())] = path
}

func (r *HeadlessResourceManager) IdOf(path string) a.ResId {
	for id, p := range r.resources {
		if p == path {
			return id
		}
	}

	return -1
}

func (r *HeadlessResourceManager) PathOf(id a.ResId) string {
	return r.resources[id]
}

func (r *HeadlessResourceManager) FullPathOf(id a.ResId) string {
	return "res/" + r.resources[id]
}

func (r *HeadlessResourceManager) ReadFile(id a.ResId) ([]byte, error) {
	return nil, errors.New("file reading not yet supported on testing frontend")
}

func newResourceManager() *HeadlessResourceManager {
	return &HeadlessResourceManager{
		resources: make(map[a.ResId]string),
		idgen:     common.NewIdGenerator(),
	}
}

type HeadlessRendererDelegate struct {
}

func (r *HeadlessRendererDelegate) OnPrepare() {

}

func (r *HeadlessRendererDelegate) OnCreatePrimitiveRenderingContext(ctx *rendering.PrimitiveRenderingContext) {

}

func (r *HeadlessRendererDelegate) OnPerformRenderingStart() {

}

func (r *HeadlessRendererDelegate) OnPerformRenderingEnd() {

}

func (r *HeadlessRendererDelegate) OnClear() {

}

func (r *HeadlessRendererDelegate) OnStop() {

}

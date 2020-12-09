// +build js

package rendering

import (
	"github.com/cadmean-ru/amphion/common"
	"fmt"
	"syscall/js"
)

type P5Renderer struct {
	commandsBufferJs js.Value
	execCommandsJs   js.Value
	postPrimitivesJs js.Value
	primitives       map[int64]*primitiveContainer
	idgen            *common.IdGenerator
	shouldDelete     bool
}

func (r *P5Renderer) Prepare() {
	r.commandsBufferJs = js.Global().Get("commandsBuffer")
	r.postPrimitivesJs = js.Global().Get("postPrimitives")
	r.execCommandsJs   = js.Global().Get("execCommands")
	r.Clear()
}

func (r *P5Renderer) AddPrimitive() int64 {
	id := r.idgen.NextId()
	r.primitives[id] = newPrimitiveContainer(id)
	return id
}

func (r *P5Renderer) SetPrimitive(id int64, primitive PrimitiveBuilder, shouldRedraw bool) {
	if !shouldRedraw {
		return
	}

	if pc, ok := r.primitives[id]; ok {
		if pc.status != primitiveStatusNew {
			pc.status = primitiveStatusRedraw
		}
		pc.primitive = primitive.BuildPrimitive()
	}
}

func (r *P5Renderer) RemovePrimitive(id int64) {
	if p, ok := r.primitives[id]; ok {
		r.shouldDelete = true
		p.status = primitiveStatusDeleting
	}
}

func (r *P5Renderer) PerformRendering() {
	var idsToDelete []int64

	if r.shouldDelete {
		idsToDelete = make([]int64, 1, len(r.primitives))
		for _, p := range r.primitives {
			if p.status == primitiveStatusDeleting {
				idsToDelete = append(idsToDelete, p.id)
			}
		}
	}

	commands := make([]Command, 0, len(r.primitives))

	for _, p := range r.primitives {
		if p.primitive == nil {
			panic(fmt.Sprintf("Primitive was created, but it's data was never set. Primitive id: %d", p.id))
		}

		switch p.status {
		case primitiveStatusNew:
			commands = append(commands, newAddCommand(p.id, p.primitive))
			break
		case primitiveStatusRedraw:
			commands = append(commands, newRedrawCommand(p.id, p.primitive))
			break
		case primitiveStatusDeleting:
			commands = append(commands, newRemoveCommand(p.id))
			break
		}

		p.status = primitiveStatusIdle
	}

	r.postCommands(commands)

	if idsToDelete != nil {
		for _, id := range idsToDelete {
			delete(r.primitives, id)
		}
	}

	r.shouldDelete = false
}

func (r *P5Renderer) Clear() {
	r.primitives = make(map[int64]*primitiveContainer)
	r.idgen = common.NewIdGenerator()
}

func (r *P5Renderer) postCommands(commands []Command) {
	totalSize := 1
	for _, c := range commands {
		totalSize += c.GetLength()
	}
	commandsBuf := make([]byte, totalSize)
	i := 0
	for _, c := range commands {
		c1 := c.EncodeToByteArray()
		_ = common.CopyByteArray(c1, commandsBuf, i, c.GetLength())
		//fmt.Println(c1)
		i += c.GetLength()
	}
	commandsBuf[totalSize-1] = CommandEndOfCommands
	//fmt.Println(commandsBuf)
	//bytesCopied := js.CopyBytesToJS(r.commandsBufferJs, commandsBuf)
	js.CopyBytesToJS(r.commandsBufferJs, commandsBuf)
	//fmt.Println(bytesCopied)
	r.execCommandsJs.Invoke()
}

//func (r *P5Renderer) postPrimitives(ps []common.Mappable) {
//	primitivesNative := make([]interface{}, len(ps))
//	for i, p := range ps {
//		primitivesNative[i] = p.ToMap()
//	}
//	jsPrimitives := js.ValueOf(primitivesNative)
//	r.postPrimitivesJs.Invoke(jsPrimitives)
//}

func NewRenderer() Renderer {
	return &P5Renderer{}
}
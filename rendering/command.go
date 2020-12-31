package rendering

import (
	"github.com/cadmean-ru/amphion/common/a"
)

const (
	CommandAddPrimitive = 0
	CommandRedrawPrimitive = 1
	CommandRemovePrimitive = 2
	CommandEndOfCommands = 255
)

type Command struct {
	code     byte
	dataSize int
	data     []byte
}

func (c Command) GetLength() int {
	return 5 + c.dataSize
}

func (c Command) EncodeToByteArray() []byte {
	arr := make([]byte, 5 + c.dataSize)
	arr[0] = c.code
	dataSizeBytes := a.IntToByteArray(int32(c.dataSize))
	_ = a.CopyByteArray(dataSizeBytes, arr, 1, 4)
	_ = a.CopyByteArray(c.data, arr, 5, c.dataSize)
	return arr
}

func newCommand(code byte, data []byte) Command {
	return Command{
		code:     code,
		dataSize: len(data),
		data:     data,
	}
}

func newAddCommand(id int64, primitive a.ByteArrayEncodable) Command {
	pBytes := primitive.EncodeToByteArray()
	idBytes := a.Int64ToByteArray(id)
	data := make([]byte, len(idBytes) + len(pBytes))
	_ = a.CopyByteArray(idBytes, data, 0, len(idBytes))
	_ = a.CopyByteArray(pBytes, data, len(idBytes), len(pBytes))
	return newCommand(CommandAddPrimitive, data)
}

func newRedrawCommand(id int64, primitive a.ByteArrayEncodable) Command {
	pBytes := primitive.EncodeToByteArray()
	idBytes := a.Int64ToByteArray(id)
	data := make([]byte, len(idBytes) + len(pBytes))
	_ = a.CopyByteArray(idBytes, data, 0, len(idBytes))
	_ = a.CopyByteArray(pBytes, data, len(idBytes), len(pBytes))
	return newCommand(CommandRedrawPrimitive, data)
}

func newRemoveCommand(id int64) Command {
	idBytes := a.Int64ToByteArray(id)
	return newCommand(CommandRemovePrimitive, idBytes)
}

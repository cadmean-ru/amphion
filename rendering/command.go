package rendering

import (
	"github.com/cadmean-ru/amphion/common"
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
	dataSizeBytes := common.IntToByteArray(int32(c.dataSize))
	_ = common.CopyByteArray(dataSizeBytes, arr, 1, 4)
	_ = common.CopyByteArray(c.data, arr, 5, c.dataSize)
	return arr
}

func newCommand(code byte, data []byte) Command {
	return Command{
		code:     code,
		dataSize: len(data),
		data:     data,
	}
}

func newAddCommand(id int64, primitive common.ByteArrayEncodable) Command {
	pBytes := primitive.EncodeToByteArray()
	idBytes := common.Int64ToByteArray(id)
	data := make([]byte, len(idBytes) + len(pBytes))
	_ = common.CopyByteArray(idBytes, data, 0, len(idBytes))
	_ = common.CopyByteArray(pBytes, data, len(idBytes), len(pBytes))
	return newCommand(CommandAddPrimitive, data)
}

func newRedrawCommand(id int64, primitive common.ByteArrayEncodable) Command {
	pBytes := primitive.EncodeToByteArray()
	idBytes := common.Int64ToByteArray(id)
	data := make([]byte, len(idBytes) + len(pBytes))
	_ = common.CopyByteArray(idBytes, data, 0, len(idBytes))
	_ = common.CopyByteArray(pBytes, data, len(idBytes), len(pBytes))
	return newCommand(CommandRedrawPrimitive, data)
}

func newRemoveCommand(id int64) Command {
	idBytes := common.Int64ToByteArray(id)
	return newCommand(CommandRemovePrimitive, idBytes)
}

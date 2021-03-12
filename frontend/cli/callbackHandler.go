package cli

import "github.com/cadmean-ru/amphion/frontend"

type CallbackHandler struct {
	handler frontend.CallbackHandler
}

func (h *CallbackHandler) HandleCallback(code int, data string) {
	h.handler(frontend.NewCallback(code, data))
}

func NewCallbackHandler(handler frontend.CallbackHandler) *CallbackHandler {
	return &CallbackHandler{handler: handler}
}
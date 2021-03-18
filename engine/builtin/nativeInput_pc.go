// +build windows, linux, darwin, android

package builtin

import (
	"github.com/cadmean-ru/amphion/engine"
)

func (n *NativeInputView) onInitWeb(_ engine.InitContext) {

}

func (n *NativeInputView) onStartWeb() {
}

func (n *NativeInputView) onStopWeb() {
}

func (n *NativeInputView) onDrawWeb(_ engine.DrawingContext) {
}

func (n *NativeInputView) setTextWeb(t string) {
}

func (n *NativeInputView) getTextWeb() string {
	return ""
}

func (n *NativeInputView) setHintWeb(h string) {
}

func (n *NativeInputView) getHintWeb() string {
	return ""
}
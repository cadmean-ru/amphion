package engine

import (
	"fmt"
	"testing"
)

type testMessageListener struct {
	num int
}

func (l *testMessageListener) OnMessage(m Message) bool {
	fmt.Println(fmt.Sprintf("Messag received in listener %d. Message: %s", l.num, m.Data))
	return true
}

func newTestMessageListener(num int) *testMessageListener {
	return &testMessageListener{
		num: num,
	}
}

type testMessageListener2 struct {
	s string
}

func (l *testMessageListener2) OnMessage(m Message) bool {
	fmt.Println(fmt.Sprintf("Messag received in listener %s. Message: %s", l.s, m.Data))
	return true
}

func newTestMessageListener2(s string) *testMessageListener2 {
	return &testMessageListener2{
		s: s,
	}
}

var root = newTestMessageListener(1)
var l1 = newTestMessageListener(2)
var l2 = newTestMessageListener(3)
var l3 = newTestMessageListener(4)
var l4 = newTestMessageListener2("hello")
var l5 = newTestMessageListener2("bruh")
var dispatcher = NewMessageDispatcher(root)

func testSetupMessageListeners() {
	dispatcher.AddListener(l1, root)
	dispatcher.AddListener(l2, root)
	dispatcher.AddListener(l3, l2)
	dispatcher.AddListener(l4, root)
	dispatcher.AddListener(l5, l4)
}

func TestMessageDispatcher_Dispatch(t *testing.T) {
	testSetupMessageListeners()
	dispatcher.Dispatch(NewMessage("test", 123, "test Data"))
}

func TestMessageDispatcher_DispatchDown(t *testing.T) {
	testSetupMessageListeners()
	dispatcher.DispatchDown(l2, NewMessage("test", 123, "Message from 3"))
	dispatcher.DispatchDown(l4, NewMessage("test", 123, "Message from hello"))
}
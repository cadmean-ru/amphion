package engine
//
//import "errors"
//
//type NavGraph struct {
//	nodes       map[string]NavNode
//	transitions map[string]*NavTransition
//
//	currentNode NavNode
//}
//
//func (g *NavGraph) Navigate(transition string) error {
//	var t *NavTransition
//	var ok bool
//	if t, ok = g.transitions[transition]; !ok {
//		return errors.New("transition not found")
//	}
//
//	if t.from != g.currentNode {
//		return errors.New("invalid navigation state")
//	}
//
//	instance.eventChan<-NewAmphionEvent(g, EventNavigate, t)
//
//	return nil
//}
//
//func (g *NavGraph) handleNavigateEvent(event AmphionEvent) bool {
//	return false
//}
//
//func newNavGraph() *NavGraph {
//	var ng = &NavGraph{
//		nodes:       make(map[string]NavNode),
//		transitions: make(map[string]*NavTransition),
//		currentNode: nil,
//	}
//	instance.BindEventHandler(EventNavigate, ng.handleNavigateEvent)
//	return ng
//}
//
//type NavTransition struct {
//	name string
//	from NavNode
//	to   NavNode
//}
//
//func (t *NavTransition) GetName() string {
//	return t.name
//}
//
//type NavNode struct {
//	name string
//}
//
//func (n *NavNode) GetName() string {
//	return n.name
//}
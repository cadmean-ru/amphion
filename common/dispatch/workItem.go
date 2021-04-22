package dispatch

type WorkItem interface {
	Execute()
}

type WorkItemFunc struct {
	work func()
}

func (w *WorkItemFunc) Execute() {
	if w.work == nil {
		return
	}

	w.work()
}

func NewWorkItemFunc(work func()) *WorkItemFunc {
	return &WorkItemFunc{
		work: work,
	}
}
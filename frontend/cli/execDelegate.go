package cli

type ExecDelegate struct {
	delegate func()
}

func (e *ExecDelegate) Execute() {
	e.delegate()
}

func NewExecDelegate(delegate func()) *ExecDelegate {
	return &ExecDelegate{delegate: delegate}
}
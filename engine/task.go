package engine

type Task struct {
	run         TaskRunner
	callback    TaskCallback
	errCallback TaskErrorCallback
}

func NewTask(runner TaskRunner, callback TaskCallback, errorCallback TaskErrorCallback) *Task {
	return &Task{
		run:         runner,
		callback:    callback,
		errCallback: errorCallback,
	}
}

type TaskRunner func() (interface{}, error)

type TaskCallback func(res interface{})

type TaskErrorCallback func(err error)

type TasksRoutine struct {
	taskChan chan *Task
}

func (r *TasksRoutine) start() {
	go r.run()
}

func (r *TasksRoutine) run() {
	defer instance.recover()

	for task := range r.taskChan {
		if task == nil {
			close(r.taskChan)
			return
		}

		if task.run == nil {
			continue
		}

		res, err := task.run()

		if err != nil && task.errCallback != nil {
			task.errCallback(err)
		} else if task.callback != nil {
			task.callback(res)
		}
	}
}

func (r *TasksRoutine) stop() {
	r.taskChan<-nil
}

func (r *TasksRoutine) RunTask(task *Task) {
	r.taskChan<-task
}

func newTasksRoutine() *TasksRoutine {
	return &TasksRoutine{
		taskChan: make(chan *Task, 100),
	}
}

type TaskBuilder struct {
	task *Task
}

func (b *TaskBuilder) Run(runner TaskRunner) *TaskBuilder {
	b.task.run = runner
	return b
}

func (b *TaskBuilder) Then(callback TaskCallback) *TaskBuilder {
	b.task.callback = callback
	return b
}

func (b *TaskBuilder) Err(callback TaskErrorCallback) *TaskBuilder {
	b.task.errCallback = callback
	return b
}

func (b *TaskBuilder) Build() *Task {
	return b.task
}

func NewTaskBuilder() *TaskBuilder {
	return &TaskBuilder{
		task: &Task{},
	}
}
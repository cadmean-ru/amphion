package engine

//Task is used to run asynchronous code.
//It contains the function to be executed asynchronously as well as callbacks to be called after the task is executed.
type Task struct {
	run             TaskRunner
	callback        TaskCallback
	errCallback     TaskErrorCallback
	finallyCallback TaskFinallyCallback
}

//NewTask creates a new task given the function that should be executed asynchronously and all the callbacks.
//All callbacks are optional.
//Prefer NewTaskBuilder to create a new task object.
func NewTask(runner TaskRunner, callback TaskCallback, errorCallback TaskErrorCallback, finally TaskFinallyCallback) *Task {
	return &Task{
		run:             runner,
		callback:        callback,
		errCallback:     errorCallback,
		finallyCallback: finally,
	}
}

type TaskRunner func() (interface{}, error)

type TaskCallback func(res interface{})

type TaskErrorCallback func(err error)

type TaskFinallyCallback func()

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

		if task.finallyCallback != nil {
			task.finallyCallback()
		}
	}
}

func (r *TasksRoutine) stop() {
	r.taskChan <- nil
}

func (r *TasksRoutine) RunTask(task *Task) {
	r.taskChan <- task
}

func newTasksRoutine() *TasksRoutine {
	return &TasksRoutine{
		taskChan: make(chan *Task, 100),
	}
}

//TaskBuilder struct for creating task objects easier with function call chain.
type TaskBuilder struct {
	task *Task
}

//Run specifies the main task's function.
func (b *TaskBuilder) Run(runner TaskRunner) *TaskBuilder {
	b.task.run = runner
	return b
}

//Then specifies task's then callback, that is called after successful task completion.
func (b *TaskBuilder) Then(callback TaskCallback) *TaskBuilder {
	b.task.callback = callback
	return b
}

//Err specifies task's err callback, that is called if task completed with error.
func (b *TaskBuilder) Err(callback TaskErrorCallback) *TaskBuilder {
	b.task.errCallback = callback
	return b
}

//Finally specifies task's finally callback, that is after task execution is completed.
func (b *TaskBuilder) Finally(callback TaskFinallyCallback) *TaskBuilder {
	b.task.finallyCallback = callback
	return b
}

//Build creates the task.
func (b *TaskBuilder) Build() *Task {
	return b.task
}

//NewTaskBuilder creates new TaskBuilder. This is the preferred way to create a new task object.
func NewTaskBuilder() *TaskBuilder {
	return &TaskBuilder{
		task: &Task{},
	}
}

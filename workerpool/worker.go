package workerpool

type Task interface {
	Perform() (interface{}, error)
}

type Result struct {
	RequestedTask *Task
	Error         error
	TaskResult    interface{}
}

type TaskChannel chan Task
type TaskQueue chan chan Task
type ResultChannel chan Result

type Worker struct {
	WorkerId int
	Tasks    TaskChannel
	Queue    TaskQueue
	Results  ResultChannel
}

func NewWorker(id int, queue TaskQueue, results ResultChannel) *Worker {
	return &Worker{id, make(chan Task), queue, results}
}

func (w *Worker) Start() {
	go func() {
		for {
			w.Queue <- w.Tasks
			task := <-w.Tasks
			result := Result{RequestedTask: &task}
			result.TaskResult, result.Error = task.Perform()
			w.Results <- result
		}
	}()
}

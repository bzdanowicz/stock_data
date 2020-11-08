package workerpool

import (
	"sync"
)

type Task interface {
	Execute() (interface{}, error)
	GetParameter() string
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
	WorkerId  int
	Tasks     TaskChannel
	Queue     TaskQueue
	Results   ResultChannel
	waitGroup *sync.WaitGroup
}

func NewWorker(id int, queue TaskQueue, results ResultChannel, wg *sync.WaitGroup) *Worker {
	return &Worker{id, make(chan Task), queue, results, wg}
}

func (w *Worker) Start() {
	go func() {
		for {
			w.Queue <- w.Tasks
			task := <-w.Tasks
			result := Result{RequestedTask: &task}
			result.TaskResult, result.Error = task.Execute()
			w.waitGroup.Done()
			w.Results <- result
		}
	}()
}

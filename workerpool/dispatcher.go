package workerpool

import (
	"sync"
)

type Dispatcher struct {
	Workers     []*Worker
	WorkChannel TaskChannel
	WorkQueue   TaskQueue
	Results     ResultChannel
	waitGroup   sync.WaitGroup
}

func NewDispatcher(workersNumber int, resultBuffer int) *Dispatcher {
	return &Dispatcher{
		Workers:     make([]*Worker, workersNumber),
		WorkChannel: make(TaskChannel),
		WorkQueue:   make(TaskQueue),
		Results:     make(ResultChannel, resultBuffer),
	}
}

func (dispatcher *Dispatcher) Start() {
	size := len(dispatcher.Workers)
	for i := 0; i < size; i++ {
		worker := NewWorker(i, dispatcher.WorkQueue, dispatcher.Results, &dispatcher.waitGroup)
		worker.Start()
		dispatcher.Workers = append(dispatcher.Workers, worker)
	}

	go dispatcher.run()
}

func (dispatcher *Dispatcher) run() {
	for {
		job := <-dispatcher.WorkChannel
		jobChannel := <-dispatcher.WorkQueue
		jobChannel <- job
	}
}

func (dispatcher *Dispatcher) Enqueue(task Task) {
	dispatcher.waitGroup.Add(1)
	dispatcher.WorkChannel <- task
}

func (dispatcher *Dispatcher) GetResult() Result {
	if len(dispatcher.Results) == 0 {
		return Result{}
	}
	return <-dispatcher.Results
}

func (dispatcher *Dispatcher) WaitAllFinished() {
	dispatcher.waitGroup.Wait()
}

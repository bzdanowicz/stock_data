package workerpool

type Dispatcher struct {
	Workers     []*Worker
	WorkChannel TaskChannel
	WorkQueue   TaskQueue
	Results     ResultChannel
}

func NewDispatcher(size int) *Dispatcher {
	return &Dispatcher{
		Workers:     make([]*Worker, size),
		WorkChannel: make(TaskChannel),
		WorkQueue:   make(TaskQueue),
		Results:     make(ResultChannel),
	}
}

func (dispatcher *Dispatcher) Start() {
	size := len(dispatcher.Workers)
	for i := 0; i < size; i++ {
		worker := NewWorker(i, dispatcher.WorkQueue, dispatcher.Results)
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
	dispatcher.WorkChannel <- task
}

func (dispatcher *Dispatcher) GetResult() Result {
	return <-dispatcher.Results
}

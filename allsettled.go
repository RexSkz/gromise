package gromise

type AllSettledValue struct {
	Status GromiseStatus
	Value  interface{}
	Reason error
}

type AllSettledResult struct {
	finishedCh chan bool
	timeoutCh  chan bool
	timeout    bool
	values     []*AllSettledValue
}

func newAllSettledResult() *AllSettledResult {
	result := &AllSettledResult{
		finishedCh: make(chan bool, 1),
		timeoutCh:  make(chan bool, 1),
		timeout:    false,
		values:     []*AllSettledValue{},
	}
	return result
}

// Run all the functions in parallel and wait for all of them to finish.
// Notice that gromise can't terminate the goroutines since it's the Golang's
// limitation, so if there are any goroutines that are still running after
// the timeout, they will be left running, which may cause goroutine leak.
func (r *AllSettledResult) Await() ([]*AllSettledValue, error) {
	select {
	case <-r.finishedCh:
		return r.values, nil
	case <-r.timeoutCh:
		return r.values, ErrorTimeout
	}
}

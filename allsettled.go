package gromise

type AllSettledValue struct {
	Status GromiseStatus
	Value  interface{}
	Reason error
}

type AllSettledResult struct {
	finished chan bool
	timeout  chan bool
	values   []*AllSettledValue
}

func newAllSettledResult() *AllSettledResult {
	result := &AllSettledResult{
		finished: make(chan bool),
		timeout:  make(chan bool),
		values:   []*AllSettledValue{},
	}
	return result
}

func (r *AllSettledResult) Await() ([]*AllSettledValue, error) {
	select {
	case <-r.finished:
		return r.values, nil
	case <-r.timeout:
		return r.values, ErrorTimeout
	}
}

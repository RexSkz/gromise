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

func (r *AllSettledResult) Await() ([]*AllSettledValue, error) {
	select {
	case <-r.finishedCh:
		return r.values, nil
	case <-r.timeoutCh:
		return r.values, ErrorTimeout
	}
}

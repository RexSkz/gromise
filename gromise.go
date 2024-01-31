package gromise

import (
	"context"
	"errors"
	"time"
)

type Gromise struct {
	// TODO: implement max concurrency
	// The default value is `runtime.NumCPU()`.
	// maxConcurrency int

	timeoutMs int
}

func New(timeoutMs int) *Gromise {
	return &Gromise{
		// maxConcurrency: runtime.NumCPU(),
		timeoutMs: timeoutMs,
	}
}

// func (g *Gromise) SetMaxConcurrency(n int) error {
// 	g.maxConcurrency = n
// 	return nil
// }

type Executor func() (interface{}, error)

func (g *Gromise) AllSettled(fns []Executor) *AllSettledResult {
	result := newAllSettledResult()

	go func() {
		if len(fns) == 0 {
			result.finishedCh <- true
			return
		}

		result.values = make([]*AllSettledValue, len(fns))
		for index := range result.values {
			result.values[index] = &AllSettledValue{
				Status: StatusPending,
			}
		}

		goroutinesToWait := len(fns)
		for index, itemFn := range fns {
			go func(itemFn func() (interface{}, error), index int) {
				defer func() {
					goroutinesToWait--
					if r := recover(); r != nil {
						result.values[index].Status = StatusRejected
						switch x := r.(type) {
						case string:
							result.values[index].Reason = errors.New(x)
						case error:
							result.values[index].Reason = x
						default:
							result.values[index].Reason = ErrorUnknownPanic
						}
					}
				}()

				if itemResult, err := itemFn(); err != nil {
					if !result.timeout {
						result.values[index].Status = StatusRejected
						result.values[index].Reason = err
					}
				} else {
					if !result.timeout {
						result.values[index].Status = StatusFulfilled
						result.values[index].Value = itemResult
					}
				}
			}(itemFn, index)
		}

		ctxWithTimeout, clearTimeout := context.WithTimeout(context.Background(), time.Duration(g.timeoutMs)*time.Millisecond)
		for {
			select {
			case <-ctxWithTimeout.Done():
				result.timeoutCh <- true
				result.timeout = true
				clearTimeout()
				for index := range result.values {
					if result.values[index].Status == StatusPending {
						result.values[index].Status = StatusRejected
						result.values[index].Reason = ErrorTimeout
					}
				}
				return
			default:
				if goroutinesToWait <= 0 {
					result.finishedCh <- true
					clearTimeout()
					return
				}
			}
		}
	}()

	return result
}

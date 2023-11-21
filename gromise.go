package gromise

import (
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
			result.finished <- true
			return
		}

		goroutinesToWait := len(fns)
		now := time.Now()

		result.values = make([]*AllSettledValue, len(fns))
		for index := range result.values {
			result.values[index] = &AllSettledValue{
				Status: StatusPending,
			}
		}

		for index, fn := range fns {
			go func(fn func() (interface{}, error), index int) {
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

				if r, err := fn(); err != nil {
					result.values[index].Status = StatusRejected
					result.values[index].Reason = err
				} else {
					result.values[index].Status = StatusFulfilled
					result.values[index].Value = r
				}
			}(fn, index)
		}

		for {
			if goroutinesToWait <= 0 {
				result.finished <- true
				break
			}
			if time.Since(now).Milliseconds() > int64(g.timeoutMs) {
				result.timeout <- true
				break
			}
		}
	}()

	return result
}

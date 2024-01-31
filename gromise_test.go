package gromise

import (
	"errors"
	"testing"
	"time"

	"gopkg.in/go-playground/assert.v1"
)

func TestNewGromise(t *testing.T) {
	fns := []Executor{
		func() (interface{}, error) {
			time.Sleep(100 * time.Millisecond)
			return 1, nil
		},
		func() (interface{}, error) {
			time.Sleep(200 * time.Millisecond)
			return 2, nil
		},
		func() (interface{}, error) {
			time.Sleep(300 * time.Millisecond)
			return 3, nil
		},
	}

	prev := time.Now()
	results, err := New(1000).AllSettled(fns).Await()
	after := time.Now()

	if err != nil {
		t.Errorf("err should be nil, got %v", err)
	}

	if len(results) != 3 {
		t.Errorf("the length of results should be 3, got %d", len(results))
	}

	assert.Equal(t, results, []*AllSettledValue{
		{
			Status: StatusFulfilled,
			Value:  1,
		},
		{
			Status: StatusFulfilled,
			Value:  2,
		},
		{
			Status: StatusFulfilled,
			Value:  3,
		},
	})

	elapsedTime := after.Sub(prev)
	// 350 ms is 300ms + 50ms epsilon
	if elapsedTime > 350*time.Millisecond {
		t.Errorf("fns should be executed concurrently, but used %d ms", elapsedTime.Milliseconds())
	}
}

func TestError(t *testing.T) {
	fns := []Executor{
		func() (interface{}, error) {
			time.Sleep(100 * time.Millisecond)
			return 1, nil
		},
		func() (interface{}, error) {
			return nil, errors.New("you shall not pass")
		},
	}

	results, err := New(1000).AllSettled(fns).Await()
	if err != nil {
		t.Errorf("err should be nil, got %v", err)
	}

	assert.Equal(t, results, []*AllSettledValue{
		{
			Status: StatusFulfilled,
			Value:  1,
		},
		{
			Status: StatusRejected,
			Reason: errors.New("you shall not pass"),
		},
	})
}

func TestPanic(t *testing.T) {
	fns := []Executor{
		func() (interface{}, error) {
			time.Sleep(100 * time.Millisecond)
			return 1, nil
		},
		func() (interface{}, error) {
			panic("blahblah")
		},
		func() (interface{}, error) {
			panic(errors.New("an error"))
		},
		func() (interface{}, error) {
			panic(123)
		},
	}

	results, err := New(1000).AllSettled(fns).Await()
	if err != nil {
		t.Errorf("err should be nil, got %v", err)
	}

	assert.Equal(t, results, []*AllSettledValue{
		{
			Status: StatusFulfilled,
			Value:  1,
		},
		{
			Status: StatusRejected,
			Reason: errors.New("blahblah"),
		},
		{
			Status: StatusRejected,
			Reason: errors.New("an error"),
		},
		{
			Status: StatusRejected,
			Reason: ErrorUnknownPanic,
		},
	})
}

func TestTimeout(t *testing.T) {
	fns := []Executor{
		func() (interface{}, error) {
			time.Sleep(1000 * time.Millisecond)
			return 1, nil
		},
		func() (interface{}, error) {
			return 2, nil
		},
	}

	result, err := New(100).AllSettled(fns).Await()
	if err != ErrorTimeout {
		t.Errorf("err should be ErrorTimeout, got %v", err)
	}

	assert.Equal(t, result, []*AllSettledValue{
		{
			Status: StatusRejected,
			Reason: ErrorTimeout,
		},
		{
			Status: StatusFulfilled,
			Value:  2,
		},
	})
}

func TestEmptyFns(t *testing.T) {
	fns := []Executor{}
	results, err := New(1000).AllSettled(fns).Await()

	if err != nil {
		t.Errorf("err should be nil, got %v", err)
	}

	if len(results) != 0 {
		t.Errorf("the length of results should be 0, got %d", len(results))
	}
}

func TestTimeoutNotAffectResult(t *testing.T) {
	fns := []Executor{
		func() (interface{}, error) {
			time.Sleep(1000 * time.Millisecond)
			return 1, nil
		},
	}

	result, _ := New(100).AllSettled(fns).Await()

	// If fn is still running, the result status will be changed to
	// StatusFulfilled after 1 second, which is not expected.
	// The expected behaviour should be StatusReject.
	time.Sleep(1000 * time.Millisecond)

	assert.Equal(t, result, []*AllSettledValue{
		{
			Status: StatusRejected,
			Reason: ErrorTimeout,
		},
	})
}

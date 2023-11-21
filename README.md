# gromise

A library to execute goroutines like `Promise.allSettled` in JavaScript, make some scenarios (e.g. BFF data aggregation) easier.

## Usage

```go
import (
    "errors"
    "github.com/rexskz/gromise"
)

// 1. define a list of functions
fns := []gromise.Executor{
    func() (interface{}, error) {
        // business logic that may cause lots of time
        return data, nil
    },
    func() (interface{}, error) {
        // a function that can return error
        return nil, errors.New("this is an error")
    },
    func() (interface{}, error) {
        // panic will be converted to error
        panic("fatal error")
    },
    // ...other functions
}

// 2. execute all functions concurrently and block until all
// functions are finished, throw error, panic, or timeout
timeoutMs := 1000
results, err := gromise.New(timeoutMs).AllSettled(fns).Await()

// 3. the results (if not timeout) will be like:
assert.Equal(t, results, []*gromise.AllSettledValue{
    {
        Status: gromise.StatusFulfilled,
        Value:  data,
    },
    {
        Status: gromise.StatusRejected,
        Reason: errors.New("this is an error"),
    },
    {
        Status: gromise.StatusRejected,
        Reason: errors.New("fatal error"),
    },
})

// 4. if timeout, the error will be:
assert.Equal(t, err, gromise.ErrTimeout)
```

## Test & Coverage

```bash
./scripts/run_tests.sh
```

You can find the coverage report in `./coverage/index.html`.

## Why It's Useful

Sometimes we need to execute a list of functions concurrently and block until all functions are finished, e.g. BFF data aggregation. In JavaScript, we can use `Promise.all` to achieve this, but in Go, we have to use `sync.WaitGroup` or `chan`, and handle the `panic` manually, which is not so convenient. This library is to make this scenario easier.

## License

MIT

package gromise

import (
	"errors"
)

type GromiseStatus string

const (
	StatusPending   GromiseStatus = "pending"
	StatusFulfilled GromiseStatus = "fulfilled"
	StatusRejected  GromiseStatus = "rejected"
)

var ErrorTimeout = errors.New("gromise: timeout")
var ErrorUnknownPanic = errors.New("unknown panic")

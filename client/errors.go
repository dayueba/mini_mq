package client

import "errors"

var ErrStopped = errors.New("stopped")

var ErrNotConnected = errors.New("not connected")

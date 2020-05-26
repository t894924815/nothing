package stat

import (
	"io"
)

type TrafficMeter interface {
	Count(sent uint64, recv uint64)
	Query() (sent uint64, recv uint64)
	io.Closer
}

type Authenticator interface {
	CheckHash(hash string) bool
	io.Closer
}
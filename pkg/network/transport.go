package network

import "time"

// Transport is an interface that abstracts the underlying network transport.
// It could be backed by any kind (Thrift, GRPC, plain TCP,..)
type Transport interface {
	Dial(string, time.Duration) error
	Listen(string) error
	Close()
}

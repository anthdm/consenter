package consensus

import pb "github.com/anthdm/consenter/pkg/protos"

// Engine is an interface abstraction for an algorithm agnostic consensus engine.
type Engine interface {
	// AddTransaction will be called each time the server sees a tx for the
	// first time.
	AddTransaction(*pb.Transaction)
}

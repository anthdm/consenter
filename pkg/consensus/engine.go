package consensus

import (
	"crypto/ecdsa"

	pb "github.com/anthdm/consenter/pkg/protos"
)

// Engine is an interface abstraction for an algorithm agnostic consensus engine.
type Engine interface {
	// Configurate will be called on server startup, where the server will pass
	// its relay channel, which can be used to relay generated blocks and
	// consensus messages into the network and the private key.
	Configurate(chan<- *pb.Message, *ecdsa.PrivateKey)
	// AddTransaction will be called each time the server sees a tx for the
	// first time.
	AddTransaction(*pb.Transaction)
}

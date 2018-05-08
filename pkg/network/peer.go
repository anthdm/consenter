package network

import pb "github.com/anthdm/consenter/pkg/protos"

// Peer represents a remote node in the network its an interface the may be
// backed by any concrete transport.
type Peer interface {
	Send(*pb.Message) error
	Disconnect(error)
	Endpoint() string
}

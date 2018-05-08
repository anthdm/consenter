package network

import (
	"net"

	"github.com/anthdm/consenter/pkg/common/codec"
	pb "github.com/anthdm/consenter/pkg/protos"
)

// TCPPeer represents a remote node backed by TCP transport.
type TCPPeer struct {
	// underlying TCP connection
	conn  net.Conn
	errCh chan error
}

// NewTCPPeer returns a new TCPPeer object.
func NewTCPPeer(conn net.Conn) *TCPPeer {
	return &TCPPeer{
		conn: conn,
	}
}

// Send implements the Peer interface.
func (p *TCPPeer) Send(msg *pb.Message) error {
	select {
	case err := <-p.errCh:
		return err
	default:
		return codec.EncodeProto(p.conn, msg)
	}
}

// Disconnect implements the Peer interface.
func (p *TCPPeer) Disconnect(err error) {
	p.conn.Close()
	close(p.errCh)
}

// Endpoint implements the Peer interface.
func (p *TCPPeer) Endpoint() string {
	return p.conn.RemoteAddr().String()
}

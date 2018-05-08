package network

import (
	"net"
	"strings"
	"time"

	"github.com/anthdm/consenter/pkg/common/codec"
	pb "github.com/anthdm/consenter/pkg/protos"
	log "github.com/sirupsen/logrus"
)

// TCPTransport represents network transportation backed by plain TCP.
type TCPTransport struct {
	// Reference to the top level server.
	srv *Server
	// Underlying TCP listener.
	listener net.Listener
}

// NewTCPTransport return a new TCPTransport.
func NewTCPTransport(s *Server) *TCPTransport {
	return &TCPTransport{
		srv: s,
	}
}

// Listen implements the Transport inteface.
func (t *TCPTransport) Listen(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	t.listener = ln
	log.Infof("server.tcp accepting new connections on 0.0.0.0%s", addr)

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				if strings.Contains("close", err.Error()) {
					return
				}
				log.Warnf("server.tcp accept error: %s", err)
				continue
			}
			go t.handleConn(conn)
		}
	}()
	return nil
}

// Dial implements the Transport interface.
func (t *TCPTransport) Dial(addr string, d time.Duration) error {
	conn, err := net.DialTimeout("tcp", addr, d)
	if err != nil {
		return err
	}
	go t.handleConn(conn)
	return nil
}

// Close implements the Transport interface.
func (t *TCPTransport) Close() {
	if t.listener != nil {
		t.listener.Close()
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	var (
		err  error
		peer = NewTCPPeer(conn)
	)
	t.srv.addPeer <- peer

	for {
		msg := pb.Message{}
		if err = codec.DecodeProto(conn, &msg); err != nil {
			break
		}
		t.srv.protoCh <- messageTuple{
			peer: peer,
			msg:  &msg,
		}
	}
}

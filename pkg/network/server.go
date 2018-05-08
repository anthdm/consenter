package network

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/anthdm/consenter/pkg/common"
	"github.com/anthdm/consenter/pkg/consensus"
	pb "github.com/anthdm/consenter/pkg/protos"
	"github.com/anthdm/consenter/pkg/storage"
	log "github.com/sirupsen/logrus"
)

var errServerShutdown = errors.New("server shutting down")

// ServerConfig holds the server configuration.
type ServerConfig struct {
	// The listen address of the server.
	ListenAddr int

	// A list of seed nodes to bootstrap the initial network.
	BootstrapNodes []string

	// The number of seconds it may take for dialing outbound connections.
	// Note that you need to pass N * time.Second to get the time.Duration
	// in seconds.
	DialTimeout time.Duration

	// When set to true this node will act as a faulty node in the network.
	Faulty bool

	// Whether this node will act as a consensus node (block producer).
	Consensus bool

	// PrivateKey of the server. This can be left empty if consensus is set
	// to false.
	PrivateKey *ecdsa.PrivateKey
}

type (
	// Server represents a remote node in the p2p network.
	Server struct {
		// Configuration of the server.
		ServerConfig

		// Underlying transport of the network for exchanging messages.
		transport Transport

		// Engine is the attached consensus engine algorithm, responsible for
		// proposing blocks.
		engine consensus.Engine

		// Tuple used for message communication between the server and
		// its transport. It holds both the message and the peer.
		protoCh chan messageTuple

		// relayCh is used to communicate between the server and its underlying
		// consensus engine.
		relayCh chan *pb.Message

		// RelayCache holds transaction hashes that this server already relayed
		// to its peers.
		relayCache *storage.MemStore

		// Peers is a map of current connected peers to the server.
		peers   map[Peer]bool
		addPeer chan Peer
		delPeer chan peerDrop

		// Waitgroup for orchestrate a gracefull shutdown.
		wg sync.WaitGroup

		// Field members to orchestrate a gracefull server shutdown.
		lock    sync.Mutex
		quit    chan struct{}
		running bool
	}

	// messageTuple holds information between the message send through the
	// network and the peer.
	messageTuple struct {
		peer Peer
		msg  *pb.Message
	}

	// peerDrop is used for disconnecting peers. It holds the peer that need
	// to be disconnected along with the reason.
	peerDrop struct {
		peer   Peer
		reason error
	}
)

// NewServer returns a new Server object.
func NewServer(cfg ServerConfig, engine consensus.Engine) *Server {
	s := &Server{
		ServerConfig: cfg,
		peers:        make(map[Peer]bool),
		addPeer:      make(chan Peer),
		delPeer:      make(chan peerDrop),
		protoCh:      make(chan messageTuple),
		relayCache:   storage.NewMemStore(),
		relayCh:      make(chan *pb.Message),
	}
	if engine != nil {
		s.engine = engine
		s.engine.Configurate(s.relayCh, s.PrivateKey)
	}
	return s
}

// Start attempts to start running the server.
func (s *Server) Start() error {
	s.lock.Lock()
	if s.running {
		return errors.New("server already running")
	}
	s.running = true
	s.lock.Unlock()

	log.Info("starting p2p server..")
	ts := NewTCPTransport(s)
	if err := s.listen(ts); err != nil {
		return err
	}
	go s.run()
	go s.bootstrapNetwork()
	go s.generateTxLoop()
	s.wg.Add(1)
	s.wg.Wait()
	return nil
}

func (s *Server) bootstrapNetwork() {
	for _, addr := range s.BootstrapNodes {
		if err := s.transport.Dial(addr, s.DialTimeout); err != nil {
			// Do 1 more attempt to allow docker containers to connect.
			time.Sleep(5 * time.Second)
			if err := s.transport.Dial(addr, s.DialTimeout); err != nil {
				log.Warnf("failed to dial (%s) reason: %s", addr, err)
			}
		}
	}
}

func (s *Server) listen(ts Transport) error {
	if err := ts.Listen(fmt.Sprintf(":%d", s.ListenAddr)); err != nil {
		return err
	}
	s.transport = ts
	return nil
}

func (s *Server) run() {
running:
	for {
		select {
		case <-s.quit:
			break running
		case msg := <-s.relayCh:
			s.Relay(msg)
		case t := <-s.protoCh:
			if err := s.handleMessage(t.peer, t.msg); err != nil {
				log.Warnf("failed processing message: %s", err)
			}
		case p := <-s.addPeer:
			s.peers[p] = true
			log.WithFields(log.Fields{
				"endpoint": p.Endpoint(),
			}).Info("new peer connected")
		case t := <-s.delPeer:
			delete(s.peers, t.peer)
			log.WithFields(log.Fields{
				"endpoint": t.peer.Endpoint(),
				"reason":   t.reason,
			}).Warn("peer disconnected")
		}
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	if !s.running {
		return
	}
	for peer := range s.peers {
		peer.Disconnect(errServerShutdown)
		delete(s.peers, peer)
	}
	if s.transport != nil {
		s.transport.Close()
	}
	s.running = false
	s.wg.Done()
}

// Relay will forward any given message to the connected peers.
func (s *Server) Relay(msg *pb.Message) {
	for peer := range s.peers {
		go func(peer Peer) {
			if err := peer.Send(msg); err != nil {
				log.Warnf("failed to relay message to peer (%s) reason: %s",
					peer.Endpoint(), err)
			}
		}(peer)
	}
}

func (s *Server) handleMessage(peer Peer, msg *pb.Message) error {
	switch p := msg.Payload.(type) {
	case *pb.Message_Transaction:
		// We already seen and relayed this tx.
		if s.relayCache.Has(p.Transaction.Hash()) {
			return nil
		}
		log.Infof("receiving new tx: %s",
			hex.EncodeToString(p.Transaction.Hash()))

		s.relayCache.Put(p.Transaction.Hash(), nil)
		s.Relay(msg)
		s.addTransaction(p.Transaction)
	}
	return nil
}

func (s *Server) addTransaction(tx *pb.Transaction) {
	if s.engine != nil {
		s.engine.AddTransaction(tx)
	}
}

// generateTxLoop will create and relay random transactions to all connected
// peers, simulating transactions created by clients to the node.
func (s *Server) generateTxLoop() {
	for {
		tx := pb.NewTransaction()
		msg := &pb.Message{
			Payload: &pb.Message_Transaction{
				Transaction: tx,
			},
		}
		s.Relay(msg)
		s.relayCache.Put(tx.Hash(), nil)
		s.addTransaction(tx)
		d := time.Duration(common.RandInt(1, 4)) * time.Second
		time.Sleep(d)
	}
}

package solo

import (
	"crypto/ecdsa"
	"time"

	pb "github.com/anthdm/consenter/pkg/protos"
)

// Engine represents a single consensus engine ^^. This is used as implementation
// example.
type Engine struct {
	blockGenerationInterval time.Duration
	privKey                 *ecdsa.PrivateKey
	relayCh                 chan<- *pb.Message
	transactions            []*pb.Transaction
}

// NewEngine returns a new "Solo" consensus engine.
func NewEngine(interval time.Duration) *Engine {
	return &Engine{
		blockGenerationInterval: interval,
		transactions:            []*pb.Transaction{},
	}
}

// Configurate implements the Engine interface.
func (e *Engine) Configurate(relayCh chan<- *pb.Message, priv *ecdsa.PrivateKey) {
	e.privKey = priv
	e.relayCh = relayCh
	go e.run()
}

func (e *Engine) run() {
	var (
		timer = time.NewTimer(e.blockGenerationInterval)
		index uint32
	)
	for {
		select {
		case <-timer.C:
			block := pb.NewBlock(index)
			block.Transactions = e.transactions
			e.relayCh <- &pb.Message{
				Payload: &pb.Message_Block{
					Block: block,
				},
			}
			index++
			e.transactions = []*pb.Transaction{}
			timer.Reset(e.blockGenerationInterval)
		}
	}
}

// AddTransaction implements the Engine interface.
func (e *Engine) AddTransaction(tx *pb.Transaction) {
	// Assume this tx is valid.
	e.transactions = append(e.transactions, tx)
}

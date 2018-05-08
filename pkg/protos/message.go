package message

import (
	"math/rand"
	"time"

	"github.com/anthdm/consenter/pkg/common"
	proto "github.com/golang/protobuf/proto"
)

// NewTransaction will create a new random Transaction.
func NewTransaction() *Transaction {
	return &Transaction{
		Nonce: rand.Uint64(),
	}
}

// Hash computes the double sha256 hash.
func (tx *Transaction) Hash() []byte {
	b, err := proto.Marshal(tx)
	if err != nil {
		panic(err)
	}
	return common.Hash256(b)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

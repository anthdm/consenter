package storage

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemStore(t *testing.T) {
	s := NewMemStore()
	for i := 0; i < 1000; i++ {
		key := []byte(fmt.Sprintf("%d", i))
		val := key
		assert.Nil(t, s.Put(key, val))

		newVal, err := s.Get(key)
		assert.Nil(t, err)
		assert.Equal(t, val, newVal)
	}
}

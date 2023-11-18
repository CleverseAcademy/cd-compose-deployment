package providers

import (
	"crypto/sha512"
	"encoding/base64"
	"hash"
	"sync"

	"github.com/pkg/errors"
)

type Entropy struct {
	h hash.Hash
	sync.RWMutex
}

type IEntropy interface {
	Update([]byte) error
	HashBase64() string
}

func NewEntropyGenerator(initialEntropy []byte) *Entropy {
	h := sha512.New()
	h.Write(initialEntropy)
	return &Entropy{
		h: h,
	}
}

func (e *Entropy) Update(data []byte) error {
	e.Lock()
	defer e.Unlock()

	_, err := e.h.Write(data)
	return errors.Wrap(err, "Hash.Write")
}

func (e *Entropy) Base64Get() string {
	e.RLock()
	defer e.RUnlock()

	hashedBytes := e.h.Sum(nil)

	return base64.StdEncoding.EncodeToString(hashedBytes)
}

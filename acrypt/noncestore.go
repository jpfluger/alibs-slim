package acrypt

import "sync"

type NonceStore struct {
	mu     sync.Mutex
	nonces map[string]struct{}
}

func NewNonceStore() *NonceStore {
	return &NonceStore{nonces: make(map[string]struct{})}
}

func (ns *NonceStore) Add(nonce string) bool {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	if _, exists := ns.nonces[nonce]; exists {
		return false
	}

	ns.nonces[nonce] = struct{}{}
	return true
}

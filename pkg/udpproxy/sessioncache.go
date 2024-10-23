package udpproxy

import "sync"

// sessionCache tracks connection sessions
type sessionCache struct {
	mu   sync.RWMutex
	data map[string]*session
}

// newSessionCache creates a usable sessionCache.
func newSessionCache() sessionCache {
	return sessionCache{
		data: make(map[string]*session),
	}
}

// Read returns the associated session for an addr
func (sc *sessionCache) Read(addr string) (*session, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	v, ok := sc.data[addr]
	return v, ok
}

// Upsert overrides the session for a given addr.
func (sc *sessionCache) Upsert(addr string, session *session) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.data[addr] = session
}

// Delete removes the session for the given addr.
func (sc *sessionCache) Delete(addr string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	delete(sc.data, addr)
}

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

// read returns the associated session for an addr
func (sc *sessionCache) read(addr string) (*session, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	v, ok := sc.data[addr]
	return v, ok
}

// insert adds a session for a given addr.
func (sc *sessionCache) insert(addr string, session *session) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.data[addr] = session
}

// delete removes the session for the given addr.
func (sc *sessionCache) delete(addr string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	delete(sc.data, addr)
}

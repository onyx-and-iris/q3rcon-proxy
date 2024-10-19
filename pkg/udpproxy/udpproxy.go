package udpproxy

import (
	"net"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// Option is a functional option type that allows us to configure the Client.
type Option func(*Client)

// WithStaleTimeout is a functional option to set the stale session timeout
func WithStaleTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		if timeout < time.Minute {
			log.Warnf("cannot set stale session timeout to less than 1 minute.. defaulting to 5 minutes")
			return
		}

		c.timeout = timeout
	}
}

type Client struct {
	laddr *net.UDPAddr
	raddr *net.UDPAddr

	proxyConn *net.UDPConn

	mutex    sync.RWMutex
	sessions map[string]*session

	timeout time.Duration
}

func New(port, target string, options ...Option) (*Client, error) {
	laddr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		return nil, err
	}

	raddr, err := net.ResolveUDPAddr("udp", target)
	if err != nil {
		return nil, err
	}

	c := &Client{
		laddr:    laddr,
		raddr:    raddr,
		mutex:    sync.RWMutex{},
		sessions: map[string]*session{},
		timeout:  5 * time.Minute,
	}

	for _, o := range options {
		o(c)
	}

	return c, nil
}

func (c *Client) ListenAndServe() error {
	var err error
	c.proxyConn, err = net.ListenUDP("udp", c.laddr)
	if err != nil {
		return err
	}

	go c.pruneSessions()

	buf := make([]byte, 2048)
	for {
		n, caddr, err := c.proxyConn.ReadFromUDP(buf)
		if err != nil {
			log.Error(err)
		}

		session, ok := c.sessions[caddr.String()]
		if !ok {
			session, err = newSession(caddr, c.raddr, c.proxyConn)
			if err != nil {
				log.Error(err)
				continue
			}

			c.sessions[caddr.String()] = session
		}

		go session.proxyTo(buf[:n])
	}
}

func (c *Client) pruneSessions() {
	ticker := time.NewTicker(1 * time.Minute)

	// the locks here could be abusive and i dont even know if this is a real
	// problem but we definitely need to clean up stale sessions
	for range ticker.C {
		for _, session := range c.sessions {
			c.mutex.RLock()
			if time.Since(session.updateTime) > c.timeout {
				delete(c.sessions, session.caddr.String())
				log.Tracef("session for %s deleted", session.caddr)
			}
			c.mutex.RUnlock()
		}
	}
}

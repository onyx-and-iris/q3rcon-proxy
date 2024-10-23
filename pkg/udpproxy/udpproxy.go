package udpproxy

import (
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

// Option is a functional option type that allows us to configure the Client.
type Option func(*Client)

// WithSessionTimeout is a functional option to set the session timeout
func WithSessionTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		if timeout < time.Minute {
			log.Warnf("cannot set stale session timeout to less than 1 minute.. defaulting to 5 minutes")
			return
		}

		c.sessionTimeout = timeout
	}
}

type Client struct {
	laddr *net.UDPAddr
	raddr *net.UDPAddr

	proxyConn *net.UDPConn

	sessionCache   sessionCache
	sessionTimeout time.Duration
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
		laddr:          laddr,
		raddr:          raddr,
		sessionCache:   newSessionCache(),
		sessionTimeout: 5 * time.Minute,
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

		session, ok := c.sessionCache.read(caddr.String())
		if !ok {
			session, err = newSession(caddr, c.raddr, c.proxyConn)
			if err != nil {
				log.Error(err)
				continue
			}

			c.sessionCache.insert(caddr.String(), session)
		}

		go session.proxyTo(buf[:n])
	}
}

func (c *Client) pruneSessions() {
	ticker := time.NewTicker(1 * time.Minute)

	for range ticker.C {
		for _, session := range c.sessionCache.data {
			if time.Since(session.updateTime) > c.sessionTimeout {
				c.sessionCache.delete(session.caddr.String())
				log.Tracef("session for %s deleted", session.caddr)
			}
		}
	}
}

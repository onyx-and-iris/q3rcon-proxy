package udpproxy

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type session struct {
	serverConn *net.UDPConn
	proxyConn  *net.UDPConn
	caddr      *net.UDPAddr
	updateTime time.Time

	validator
}

func newSession(caddr *net.UDPAddr, raddr *net.UDPAddr, proxyConn *net.UDPConn) (*session, error) {
	serverConn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return nil, err
	}

	session := &session{
		serverConn: serverConn,
		proxyConn:  proxyConn,
		caddr:      caddr,
		updateTime: time.Now(),
		validator:  newValidator(),
	}

	go session.listen()

	return session, nil
}

func (s *session) listen() error {
	buf := make([]byte, 2048)
	for {
		n, err := s.serverConn.Read(buf)
		if err != nil {
			log.Error(err)
			continue
		}

		go s.proxyFrom(buf[:n])
	}
}

func (s *session) proxyFrom(buf []byte) error {
	if !s.isValidResponsePacket(buf) {
		err := errors.New("not a rcon or query response packet")
		log.Error(err.Error())
		return err
	}

	s.updateTime = time.Now()
	_, err := s.proxyConn.WriteToUDP(buf, s.caddr)
	if err != nil {
		log.Error(err)
		return err
	}

	if s.isRconResponsePacket(buf) {
		if s.isBadRconResponse(buf) {
			log.Infof("Response: Bad rcon from %s", s.caddr.IP)
		} else {
			log.Debugf("Response: %s", string(buf[len(s.rconResponseHeader):]))
		}
	}

	return nil
}

func (s *session) proxyTo(buf []byte) error {
	if !s.isValidRequestPacket(buf) {
		var err error
		if s.isChallengeRequestPacket(buf) {
			parts := strings.SplitN(string(buf), " ", 3)
			err = fmt.Errorf("invalid challenge from %s with GUID: %s", s.caddr.IP, parts[len(parts)-1])
		} else {
			err = errors.New("not a rcon or query request packet")
		}
		log.Error(err.Error())
		return err
	}

	s.updateTime = time.Now()
	_, err := s.serverConn.Write(buf)
	if err != nil {
		log.Error(err)
		return err
	}

	if s.isRconRequestPacket(buf) {
		parts := strings.SplitN(string(buf), " ", 3)
		log.Infof("From [%s] To [%s] Command: %s", s.caddr.IP, s.serverConn.RemoteAddr(), parts[len(parts)-1])
	}

	return nil
}

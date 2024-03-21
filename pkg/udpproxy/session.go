package udpproxy

import (
	"errors"
	"net"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type Session struct {
	serverConn *net.UDPConn
	proxyConn  *net.UDPConn
	caddr      *net.UDPAddr
	updateTime time.Time
}

func newSession(caddr *net.UDPAddr, raddr *net.UDPAddr, proxyConn *net.UDPConn) (*Session, error) {
	serverConn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return nil, err
	}

	session := &Session{
		serverConn: serverConn,
		proxyConn:  proxyConn,
		caddr:      caddr,
		updateTime: time.Now(),
	}

	go session.listen()

	return session, nil
}

func (s *Session) isRconRequestPacket(buf []byte) bool {
	return string(buf[:8]) == "\xff\xff\xff\xffrcon"
}

func (s *Session) isQueryRequestPacket(buf []byte) bool {
	return string(buf[:13]) == "\xff\xff\xff\xffgetstatus" || string(buf[:11]) == "\xff\xff\xff\xffgetinfo"
}

func (s *Session) isValidRequestPacket(buf []byte) bool {
	return s.isRconRequestPacket(buf) || s.isQueryRequestPacket(buf)
}

func (s *Session) isRconResponsePacket(buf []byte) bool {
	return string(buf[:9]) == "\xff\xff\xff\xffprint"
}

func (s *Session) isQueryResponsePacket(buf []byte) bool {
	return string(buf[:18]) == "\xff\xff\xff\xffstatusResponse" || string(buf[:16]) == "\xff\xff\xff\xffinfoResponse"
}

func (s *Session) isValidResponsePacket(buf []byte) bool {
	return s.isRconResponsePacket(buf) || s.isQueryResponsePacket(buf)
}

func (s *Session) listen() error {
	for {
		buf := make([]byte, 2048)
		n, err := s.serverConn.Read(buf)
		if err != nil {
			log.Error(err)
			continue
		}

		go s.proxyFrom(buf[:n])
	}
}

func (s *Session) proxyFrom(buf []byte) error {
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
		log.Debugf("Response: %s", string(buf[10:]))
	}

	return nil
}

func (s *Session) proxyTo(buf []byte) error {
	if !s.isValidRequestPacket(buf) {
		err := errors.New("not a rcon or query request packet")
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
		parts := strings.Split(string(buf), " ")
		log.Infof("From [%s] To [%s] Command: %s", s.caddr.IP.String(), s.serverConn.RemoteAddr().String(), strings.Join(parts[2:], " "))
	}

	return nil
}

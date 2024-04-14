package udpproxy

import "bytes"

type validator struct {
	rconRequestHeader       []byte
	getstatusRequestHeader  []byte
	getinfoRequestHeader    []byte
	rconResponseHeader      []byte
	getstatusResponseHeader []byte
	getinfoResponseHeader   []byte
	badRconIdentifier       []byte
}

func newValidator() validator {
	v := validator{}
	v.rconRequestHeader = []byte("\xff\xff\xff\xffrcon")
	v.getstatusRequestHeader = []byte("\xff\xff\xff\xffgetstatus")
	v.getinfoRequestHeader = []byte("\xff\xff\xff\xffgetinfo")
	v.rconResponseHeader = []byte("\xff\xff\xff\xffprint\n")
	v.getstatusResponseHeader = []byte("\xff\xff\xff\xffstatusResponse\n")
	v.getinfoResponseHeader = []byte("\xff\xff\xff\xffinfoResponse\n")
	v.badRconIdentifier = []byte("Bad rcon")
	return v
}

func (v *validator) compare(buf, c []byte) bool {
	return bytes.Equal(buf[:len(c)], c)
}

func (v *validator) isRconRequestPacket(buf []byte) bool {
	return v.compare(buf, v.rconRequestHeader)
}

func (v *validator) isQueryRequestPacket(buf []byte) bool {
	return v.compare(buf, v.getstatusRequestHeader) ||
		v.compare(buf, v.getinfoRequestHeader)
}

func (v *validator) isValidRequestPacket(buf []byte) bool {
	return v.isRconRequestPacket(buf) || v.isQueryRequestPacket(buf)
}

func (v *validator) isRconResponsePacket(buf []byte) bool {
	return v.compare(buf, v.rconResponseHeader)
}

func (v *validator) isQueryResponsePacket(buf []byte) bool {
	return v.compare(buf, v.getstatusResponseHeader) ||
		v.compare(buf, v.getinfoResponseHeader)
}

func (v *validator) isValidResponsePacket(buf []byte) bool {
	return v.isRconResponsePacket(buf) || v.isQueryResponsePacket(buf)
}

func (v *validator) isBadRconResponse(buf []byte) bool {
	return v.compare(buf[len(v.rconResponseHeader):], v.badRconIdentifier)
}

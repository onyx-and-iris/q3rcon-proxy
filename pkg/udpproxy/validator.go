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
	return validator{
		rconRequestHeader:       []byte("\xff\xff\xff\xffrcon"),
		getstatusRequestHeader:  []byte("\xff\xff\xff\xffgetstatus"),
		getinfoRequestHeader:    []byte("\xff\xff\xff\xffgetinfo"),
		rconResponseHeader:      []byte("\xff\xff\xff\xffprint\n"),
		getstatusResponseHeader: []byte("\xff\xff\xff\xffstatusResponse\n"),
		getinfoResponseHeader:   []byte("\xff\xff\xff\xffinfoResponse\n"),
		badRconIdentifier:       []byte("Bad rcon"),
	}
}

func (v *validator) isRconRequestPacket(buf []byte) bool {
	return bytes.Equal(buf[:len(v.rconRequestHeader)], v.rconRequestHeader)
}

func (v *validator) isQueryRequestPacket(buf []byte) bool {
	return bytes.Equal(buf[:len(v.getstatusRequestHeader)], v.getstatusRequestHeader) || bytes.Equal(buf[:len(v.getinfoRequestHeader)], v.getinfoRequestHeader)
}

func (v *validator) isValidRequestPacket(buf []byte) bool {
	return v.isRconRequestPacket(buf) || v.isQueryRequestPacket(buf)
}

func (v *validator) isRconResponsePacket(buf []byte) bool {
	return bytes.Equal(buf[:len(v.rconResponseHeader)], v.rconResponseHeader)
}

func (v *validator) isQueryResponsePacket(buf []byte) bool {
	return bytes.Equal(buf[:len(v.getstatusResponseHeader)], v.getstatusResponseHeader) || bytes.Equal(buf[:len(v.getinfoResponseHeader)], v.getinfoResponseHeader)
}

func (v *validator) isValidResponsePacket(buf []byte) bool {
	return v.isRconResponsePacket(buf) || v.isQueryResponsePacket(buf)
}

func (v *validator) isBadRconResponse(buf []byte) bool {
	return bytes.Equal(buf[len(v.rconResponseHeader):len(v.rconResponseHeader)+len(v.badRconIdentifier)], v.badRconIdentifier)
}

package udpproxy

type validator struct {
}

func (v *validator) isRconRequestPacket(buf []byte) bool {
	return string(buf[:8]) == "\xff\xff\xff\xffrcon"
}

func (v *validator) isQueryRequestPacket(buf []byte) bool {
	return string(buf[:13]) == "\xff\xff\xff\xffgetstatus" || string(buf[:11]) == "\xff\xff\xff\xffgetinfo"
}

func (v *validator) isValidRequestPacket(buf []byte) bool {
	return v.isRconRequestPacket(buf) || v.isQueryRequestPacket(buf)
}

func (v *validator) isRconResponsePacket(buf []byte) bool {
	return string(buf[:9]) == "\xff\xff\xff\xffprint"
}

func (v *validator) isQueryResponsePacket(buf []byte) bool {
	return string(buf[:18]) == "\xff\xff\xff\xffstatusResponse" || string(buf[:16]) == "\xff\xff\xff\xffinfoResponse"
}

func (v *validator) isValidResponsePacket(buf []byte) bool {
	return v.isRconResponsePacket(buf) || v.isQueryResponsePacket(buf)
}

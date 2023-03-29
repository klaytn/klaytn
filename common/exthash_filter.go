package common

import (
	"errors"
)

type rlpDS struct {
	kind   uint8
	data   []byte
	length int
}

func RlpPaddingFilter(src []byte) (reRlp []byte, err error) {
	if len(src) <= 32 {
		return src, err
	}

	obj, err := rlpDecodeStructure(src)
	if err != nil {
		return src, nil
	}

	for _, v := range obj {
		reRlp = append(reRlp, v.data...)
	}
	// fmt.Printf("\n")

	return reRlp, err
}

func rlpCheckHeader(src []byte) (err error) {
	var totalLen int

	strLen := len(src)

	switch {
	case src[0] == 0x80:
		totalLen, err = 1, nil
	case src[0] < 0x80:
		_, _, totalLen, err = parseHeader(src, 1)
	case src[0] < 0xb8:
		_, _, totalLen, err = parseHeader(src, 2)
	case src[0] < 0xc0:
		_, _, totalLen, err = parseHeader(src, 3)
	case src[0] < 0xf8:
		_, _, totalLen, err = parseHeader(src, 4)
	default:
		_, _, totalLen, err = parseHeader(src, 5)
	}
	if strLen != totalLen {
		err = errors.New("rlp decode length error at simple RLP CheckHeader")
	}
	return err
}

func rlpDecodeStructure(src []byte) (reObj []rlpDS, err error) {
	var obj rlpDS

	tmpLen := 0
	for i := 0; i < len(src); {
		switch {
		case src[i] == 0x80:
			obj.kind = 0
			obj.data, obj.length, tmpLen, err = src[i:i+1], 1, 1, nil
		case src[i] < 0x80:
			obj.kind = 1
			obj.data, obj.length, tmpLen, err = rule1parse(src[i:])
		case src[i] < 0xb8:
			obj.kind = 2
			obj.data, obj.length, tmpLen, err = rule2parse(src[i:])
		case src[i] < 0xc0:
			obj.kind = 2
			obj.data, obj.length, tmpLen, err = rule3parse(src[i:])
		case src[i] < 0xf8:
			obj.kind = 3
			obj.data, obj.length, tmpLen, err = rule4parse(src[i:])
		default:
			obj.kind = 3
			obj.data, obj.length, tmpLen, err = rule5parse(src[i:])
		}
		if err != nil {
			break
		}

		reObj = append(reObj, obj)
		i += tmpLen
	}
	return reObj, err
}

func rule1parse(src []byte) (reStr []byte, reLen, totalLen int, err error) {
	totalLen = 0
	srcLen := len(src)
	for ; totalLen < srcLen && src[totalLen] < 0x80; totalLen += 1 {
	}

	tmpStr := src[:totalLen]
	reLen = len(tmpStr)

	reStr = append(reStr, tmpStr[:]...)
	reLen = len(tmpStr)
	return reStr, reLen, totalLen, err
}

func rule2parse(src []byte) (reStr []byte, reLen int, totalLen int, err error) {
	dataIdx, _, totalLen, err := parseHeader(src, 2)
	if err != nil {
		return reStr, reLen, totalLen, err
	}

	tmpStr := ExtPaddingFilter(src[dataIdx:totalLen])
	reStr, reLen = makepacket(tmpStr, 2)

	return reStr, reLen, totalLen, err
}

// reLen : padding 제거된후 전체 길이
// totalLen : 원래 전체 길이
func rule3parse(src []byte) (reStr []byte, reLen, totalLen int, err error) {
	dataIdx, _, totalLen, err := parseHeader(src, 3)
	if err != nil {
		return reStr, reLen, totalLen, err
	}

	tmpStr := ExtPaddingFilter(src[dataIdx:totalLen])
	reStr, reLen = makepacket(tmpStr, 3)

	return reStr, reLen, totalLen, err
}

func rule4parse(src []byte) (reStr []byte, reLen, totalLen int, err error) {
	var tmpPacket, tmpReStr []byte
	var packetList [][]byte
	var tmpPacketLen int

	dataIdx, _, totalLen, err := parseHeader(src, 4)
	if err != nil {
		return reStr, reLen, totalLen, err
	}

	tmpReLen, tmpTotalLen := int(0), int(0)
	for i := dataIdx; i < totalLen; {
		switch {
		case src[i] == 0x80:
			tmpReStr, tmpReLen, tmpTotalLen, err = src[i:i+1], 1, 1, nil
		case src[i] < 0x80:
			tmpReStr, tmpReLen, tmpTotalLen, err = rule1parse(src[i:])
		case src[i] < 0xb8:
			tmpReStr, tmpReLen, tmpTotalLen, err = rule2parse(src[i:])
		case src[i] < 0xc0:
			tmpReStr, tmpReLen, tmpTotalLen, err = rule3parse(src[i:])
		case src[i] < 0xf8:
			tmpReStr, tmpReLen, tmpTotalLen, err = rule4parse(src[i:])
		default:
			tmpReStr, tmpReLen, tmpTotalLen, err = rule5parse(src[i:])
		}
		if err != nil || i+tmpTotalLen > totalLen {
			return src[:totalLen], totalLen, totalLen, nil
		}

		packetList = append(packetList, tmpReStr)
		tmpPacket = append(tmpPacket, tmpReStr[:]...)
		tmpPacketLen += tmpReLen
		i += tmpTotalLen
	}
	reStr, reLen = makepacket(tmpPacket, 4)

	return reStr, reLen, totalLen, err
}

func rule5parse(src []byte) (reStr []byte, reLen, totalLen int, err error) {
	var tmpPacket, tmpReStr []byte
	// var packetList [][]byte
	var tmpPacketLen int

	dataIdx, _, totalLen, err := parseHeader(src, 5)
	if err != nil {
		return reStr, reLen, totalLen, err
	}

	for i := dataIdx; i < totalLen; {
		tmpReLen, tmpTotalLen := int(0), int(0)
		switch {
		case src[i] == 0x80:
			tmpReStr, tmpReLen, tmpTotalLen, err = src[i:i+1], 1, 1, nil
		case src[i] < 0x80:
			tmpReStr, tmpReLen, tmpTotalLen, err = rule1parse(src[i:])
		case src[i] < 0xb8:
			tmpReStr, tmpReLen, tmpTotalLen, err = rule2parse(src[i:])
		case src[i] < 0xc0:
			tmpReStr, tmpReLen, tmpTotalLen, err = rule3parse(src[i:])
		case src[i] < 0xf8:
			tmpReStr, tmpReLen, tmpTotalLen, err = rule4parse(src[i:])
		default:
			tmpReStr, tmpReLen, tmpTotalLen, err = rule5parse(src[i:])
		}
		if err != nil || i+tmpTotalLen > totalLen {
			return src[:totalLen], totalLen, totalLen, nil
		}

		tmpPacket = append(tmpPacket, tmpReStr[:]...)
		tmpPacketLen += tmpReLen
		i += tmpTotalLen
	}
	reStr, reLen = makepacket(tmpPacket, 5)

	return reStr, reLen, totalLen, err
}

func parseHeader(src []byte, flag int) (dataIdx, dataLen, totalLen int, err error) {
	var hexIdx int
	srcLen := len(src)
	if srcLen < 1 {
		err = errors.New("rlp decode length too small - 1... at simple RLP Decoder")
		return dataIdx, dataLen, totalLen, err
	}

	switch {
	case flag == 2 || flag == 4:
		if flag == 2 {
			hexIdx = 0x80
		} else {
			hexIdx = 0xc0
		}

		dataLen = int(src[0]) - hexIdx
		totalLen = dataLen + 1
		dataIdx = 1
	case flag == 3 || flag == 5:
		if flag == 3 {
			hexIdx = 0xb7
		} else {
			hexIdx = 0xf7
		}

		if srcLen < int(src[0])-hexIdx {
			err = errors.New("rlp decode length too small - 2... at simple RLP Decoder")
			return dataIdx, dataLen, totalLen, err
		}
		i := 1
		for ; i < srcLen && i <= int(src[0])-hexIdx; i += 1 {
			dataLen *= 0x100
			dataLen += int(src[i])
		}
		totalLen = dataLen + i
		dataIdx = i
	}
	if totalLen > srcLen || dataLen > srcLen || dataLen <= 0 || dataIdx == totalLen {
		err = errors.New("rlp decode length error at simple RLP Decoder")
	} else if (hexIdx == 0x80 && dataLen >= 55) || (hexIdx == 0xb7 && dataLen < 55) || (hexIdx == 0xc0 && dataLen >= 55) || (hexIdx == 0xf7 && dataLen < 55) {
		err = errors.New("rlp decode length error at simple RLP Decoder")
	}
	return dataIdx, dataLen, totalLen, err
}

func makepacket(data []byte, flag int) (packet []byte, packetLen int) {
	var tmpByte byte
	var procData []byte
	var tmpHeader []byte

	if tmpData, err := RlpPaddingFilter(data); err == nil {
		procData = tmpData
	} else {
		procData = data
	}
	dataLen := len(procData)

	if dataLen <= 55 {
		if flag == 2 {
			tmpByte = uint8(0x80 + dataLen)
		} else { // 4
			tmpByte = uint8(0xc0 + dataLen)
		}
		packet = append(packet, tmpByte)
		packet = append(packet, procData[:dataLen]...)
	} else {
		div := 0xffffff
		dlen := dataLen
		for ; div >= 0; div /= 0x100 {
			if dlen > div && div > 0 {
				tmpByte = uint8(dlen / div)
				tmpHeader = append(tmpHeader, tmpByte)
				div %= (div + 1)
			} else if div == 0 {
				tmpByte = uint8(dlen)
				tmpHeader = append(tmpHeader, tmpByte)
				break
			}
		}
		tmpHeaderLen := len(tmpHeader)
		if flag == 3 {
			tmpByte = uint8(0xb7 + tmpHeaderLen)
		} else { // 5
			tmpByte = uint8(0xf7 + tmpHeaderLen)
		}
		packet = append(packet, tmpByte)
		packet = append(packet, tmpHeader...)
		packet = append(packet, procData[:dataLen]...)
	}
	packetLen = len(packet)
	return packet, packetLen
}

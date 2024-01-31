package utils

import (
	"bytes"
	"encoding/binary"
)

// IntToBytes 将int类型的数转化为字节并以小端存储
func IntToBytes(intNum int) []byte {
	uint16Num := uint16(intNum)
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, uint16Num)
	return buf.Bytes()
}

// BytesToInt 将以小端存储的长为1/2字节的数转化成int类型的数
func BytesToInt(bytesArr []byte) int {
	var intNum int
	if len(bytesArr) == 1 {
		bytesArr = append(bytesArr, byte(0))
		intNum = int(binary.LittleEndian.Uint16(bytesArr))
	} else if len(bytesArr) == 2 {
		intNum = int(binary.LittleEndian.Uint16(bytesArr))
	}

	return intNum
}

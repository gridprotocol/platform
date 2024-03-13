package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
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

// duration(month) to time stamp
func DurToTS(month string) (string, error) {
	// string to uint
	m, err := strconv.ParseInt(month, 10, 64) // 第二个参数表示基数（这里是十进制），第三个参数表示目标类型的位数（这里是64位）
	if err != nil {
		return "", errors.Errorf("string to int64 error: %s", err)
	}

	//fmt.Println("month:", m)

	// get seconds
	sec := m * 30 * 86400
	//fmt.Println("seconds:", sec)

	// get current time stamp
	now := time.Now().Unix()
	//fmt.Println("current timestamp:", now)

	// get expire
	expire := now + sec
	//fmt.Println("expire timestamp:", expire)

	// int to string
	expireS := fmt.Sprintf("%d", expire)

	return expireS, nil
}

// string to uint64
func StringToUint64(s string) (uint64, error) {
	u, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return u, nil
}

func Uint64ToString(u uint64) string {
	res := strconv.FormatUint(u, 10) //uint64转字符串
	return res
}

// string to int64
func StringToInt64(s string) (int64, error) {
	u, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return u, nil
}

func Int64ToString(i int64) string {
	res := strconv.FormatInt(i, 10) //int64转字符串
	return res
}

package ic78civCmd

import (
	"fmt"
	"strconv"

	"github.com/albenik/bcd"
)

const maxReadBuff = 100

type civCommand struct {
	transiverAddr   byte
	requestFreque   []byte
	requestMode     []byte
	requestATT      []byte
	requestAFLevel  []byte
	requestRFLevel  []byte
	requestSQLLevel []byte
	requestPreamp   []byte
}

type cmdValue byte

const (
	controllerAddrCmd cmdValue = 0xe1
	preambleCmd       cmdValue = 0xfe
	sendFreqCmd       cmdValue = 0x00
	sendModeCmd       cmdValue = 0x01
	readFreqCmd       cmdValue = 0x03
	readModeCmd       cmdValue = 0x04
	setFreqCmd        cmdValue = 0x05
	setModeCmd        cmdValue = 0x06
	attCmd            cmdValue = 0x11
	afrfsqlCmd        cmdValue = 0x14
	afSubCmd          cmdValue = 0x01
	rfSubCmd          cmdValue = 0x02
	sqlSubCmd         cmdValue = 0x03
	preampCmd         cmdValue = 0x16
	preampSubCmd      cmdValue = 0x02
	readAddrCmd       cmdValue = 0x19
	endMsgCmd         cmdValue = 0xfd
	okCode            cmdValue = 0xfb
	ngCode            cmdValue = 0xfa
)

type commandName int

const (
	freqRead commandName = iota
	taddr
	mode
	att
	af
	rf
	sql
	preamp
	status
)

var myic78civCommand *civCommand

func getTransiverAddr() byte {
	return myic78civCommand.transiverAddr
}

func addElementToFirstIndex(x []byte, y byte) []byte {
	x = append([]byte{y}, x...)
	return x
}

func printByte(data []byte) {
	for _, value := range data {
		fmt.Printf("%#x ", value)
	}
	fmt.Println()
}

func intToArr(data int, size uint8) []byte {
	arr := make([]byte, len(strconv.Itoa(data)), size)
	for i := len(arr) - 1; data > 0; i-- {
		arr[i] = byte(data % 10)
		data = int(data / 10)
	}
	return arr
}

func byteArrToBCD(arr []byte, size uint8) []byte {
	buf := make([]byte, size)
	dig := len(arr)/2 - 1
	for i := 0; i < len(arr)-1; i = i + 2 {
		buf[dig] = bcd.FromUint8((arr[i] * 10) + arr[i+1])
		dig--
	}
	return buf
}

func intFreqToBcdArr(freq int) []byte {
	arr := intToArr(freq, 10)
	arr = append(arr, 0x00)
	for len(arr) < 10 {
		arr = addElementToFirstIndex(arr, 0)
	}
	return byteArrToBCD(arr, 5)
}

func bcdToInt(buff []byte) uint32 {
	return bcd.ToUint32(buff)
}

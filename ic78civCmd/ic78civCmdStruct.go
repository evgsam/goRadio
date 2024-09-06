package ic78civCmd

import (
	"fmt"
	datastruct "goRadio/dataStruct"
	"goRadio/menu"
	"strconv"
	"sync"
	"time"

	"github.com/albenik/bcd"
	"go.bug.st/serial"
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
)

func DataPollingGorutine(port serial.Port, serialAcces *sync.Mutex) {
	ch := make(chan *datastruct.RadioSettings, 10)
	go menu.Menu(ch)
	for {
		serialAcces.Lock()
		port.ResetInputBuffer()
		adr, err := requestTransiverAddr(port)
		if err != nil {
			for err != nil {
				ch <- &datastruct.RadioSettings{
					Err:    err,
					Status: "Error",
				}
				adr, err = requestTransiverAddr(port)
				time.Sleep(50 * time.Millisecond)
			}
		}
		myic78civCommand = newIc78civCommand(adr)
		mode, _ := requestMode(port, myic78civCommand)
		att, _ := requestATT(port, myic78civCommand)
		preamp, _ := requestPreamp(port, myic78civCommand)
		freq, _ := requestFreque(port, myic78civCommand)
		af, _ := requestAFLevel(port, myic78civCommand)
		rf, _ := requestRFLevel(port, myic78civCommand)
		sql, _ := requestSQLLevel(port, myic78civCommand)

		port.ResetInputBuffer()
		serialAcces.Unlock()

		ch <- &datastruct.RadioSettings{
			Err:    err,
			Status: "Connect",
			Mode:   mode,
			ATT:    att,
			Preamp: preamp,
			Freque: freq,
			AF:     af,
			RF:     rf,
			SQL:    sql,
			TrAddr: adr,
		}
		time.Sleep(1 * time.Second)
	}
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

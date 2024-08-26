package ic78civCmd

import (
	"bytes"
	"fmt"
	"goRadio/serialDataExchange"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/albenik/bcd"
	"go.bug.st/serial"
)

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
	FREQ commandName = iota
	TADDR
	MODE
	ATT
	AF
	RF
	SQL
	PREAMP
)

func DataPollingGorutine(port serial.Port, serialAcces *sync.Mutex) {
	for {
		serialAcces.Lock()
		port.ResetInputBuffer()
		transiverAddr := requestTransiverAddr(port)
		if transiverAddr == 0x00 {
			for requestTransiverAddr(port) != 0 {
				time.Sleep(50 * time.Millisecond)
			}
		}
		serialAcces.Unlock()
		time.Sleep(3 * time.Second)
	}
}

func newIc78civCommand(transiverAddr byte) *civCommand {
	ic78civCommand := &civCommand{
		transiverAddr:   transiverAddr,
		requestFreque:   []byte{byte(preambleCmd), byte(preambleCmd), transiverAddr, byte(controllerAddrCmd), byte(readFreqCmd), byte(endMsgCmd)},
		requestMode:     []byte{byte(preambleCmd), byte(preambleCmd), transiverAddr, byte(controllerAddrCmd), byte(readModeCmd), byte(endMsgCmd)},
		requestATT:      []byte{byte(preambleCmd), byte(preambleCmd), transiverAddr, byte(controllerAddrCmd), byte(attCmd), byte(endMsgCmd)},
		requestAFLevel:  []byte{byte(preambleCmd), byte(preambleCmd), transiverAddr, byte(controllerAddrCmd), byte(afrfsqlCmd), byte(afSubCmd), byte(endMsgCmd)},
		requestRFLevel:  []byte{byte(preambleCmd), byte(preambleCmd), transiverAddr, byte(controllerAddrCmd), byte(afrfsqlCmd), byte(rfSubCmd), byte(endMsgCmd)},
		requestSQLLevel: []byte{byte(preambleCmd), byte(preambleCmd), transiverAddr, byte(controllerAddrCmd), byte(afrfsqlCmd), byte(sqlSubCmd), byte(endMsgCmd)},
		requestPreamp:   []byte{byte(preambleCmd), byte(preambleCmd), transiverAddr, byte(controllerAddrCmd), byte(preampCmd), byte(preampSubCmd), byte(endMsgCmd)},
	}
	return ic78civCommand
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

func setFreque(freq int) {
	buf := make([]byte, 5)
	arr := make([]byte, len(strconv.Itoa(freq)), 10)
	for i := len(arr) - 1; freq > 0; i-- {
		arr[i] = byte(freq % 10)
		freq = int(freq / 10)
	}
	for len(arr) < 10 {
		arr = addElementToFirstIndex(arr, 0)
	}
	dig := 5
	for i := 0; i < 10; i = i + 2 {
		dig--
		buf[dig] = (arr[i] * 10) + arr[i+1]
	}
	println(buf)

}

func bcdToInt(buff []byte) uint32 {
	return bcd.ToUint32(buff)
}

func commandSend(port serial.Port, p *civCommand, c commandName) []byte {
	correctMsg := false
	readBuff := make([]byte, 30)
	dataBuff := make([]byte, 7)
	var arg []byte
	var cmd byte
	switch c {
	case FREQ:
		arg = p.requestFreque
		cmd = byte(readFreqCmd)
	case TADDR:
		arg = []byte{byte(preambleCmd), byte(preambleCmd), 0x00, byte(controllerAddrCmd), byte(readAddrCmd), 0x00, byte(endMsgCmd)}
		cmd = byte(readAddrCmd)
	case MODE:
		arg = p.requestMode
		cmd = byte(readModeCmd)
	case ATT:
		arg = p.requestATT
		cmd = byte(attCmd)
	case AF:
		arg = p.requestAFLevel
		cmd = byte(afrfsqlCmd)
	case RF:
		arg = p.requestRFLevel
		cmd = byte(afrfsqlCmd)
	case SQL:
		arg = p.requestSQLLevel
		cmd = byte(afrfsqlCmd)
	case PREAMP:
		arg = p.requestPreamp
		cmd = byte(preampCmd)
	}
	n := 0
	for !correctMsg {
		port.ResetInputBuffer()
		time.Sleep(time.Duration(100) * time.Millisecond)
		serialDataExchange.WriteSerialPort(port, arg)
		time.Sleep(time.Duration(100) * time.Millisecond)
		_ = serialDataExchange.ReadSerialPort(port, readBuff)
		for _, value := range readBuff {
			if value == 0xfd {
				n++
			}
		}
		if n < 2 {
			n = 0
			for i, _ := range readBuff {
				readBuff[i] = 0x00
			}
		} else {
			correctMsg = true
		}
	}
	for i := 0; i < n; i++ {
		idxCmd := slices.Index(readBuff, cmd)
		idxEnd := slices.Index(readBuff, byte(endMsgCmd))
		if idxEnd != -1 {
			if bytes.Equal(readBuff[:idxEnd+1], arg[:len(arg)]) {
				readBuff = readBuff[idxEnd+1 : len(readBuff)]
			} else {
				dataBuff = readBuff[idxCmd:idxEnd]
				dataBuff = append(make([]byte, 0), dataBuff[1:]...)
			}
		}
	}
	return dataBuff
}

func requestTransiverAddr(port serial.Port) byte {
	return commandSend(port, nil, TADDR)[1]
}

func requestMode(port serial.Port, p *civCommand) string {
	buff := commandSend(port, p, MODE)
	var mode string
	switch buff[0] {
	case 0x00:
		mode = "LSB"
	case 0x01:
		mode = "USB"
	case 0x02:
		mode = "AM"
	case 0x04:
		mode = "RTTY"
	case 0x07:
		mode = "CW"
	}
	return mode
}

func requestPreamp(port serial.Port, p *civCommand) string {
	var preamp string
	switch append(make([]byte, 0), commandSend(port, p, PREAMP)[1:]...)[0] {
	case 0x00:
		preamp = "OFF"
	case 0x01:
		preamp = "P.AMP"
	}
	return preamp
}

func requestATT(port serial.Port, p *civCommand) string {
	var att string
	switch commandSend(port, p, ATT)[0] {
	case 0x00:
		att = "NO"
	case 0x20:
		att = "YES"
	}
	return att
}

func requestFreque(port serial.Port, p *civCommand) uint32 {
	buff := commandSend(port, p, FREQ)
	buffRevers := make([]byte, len(buff))
	j := 0
	for i := len(buff) - 1; i > -1; i-- {
		buffRevers[j] = buff[i]
		j++
	}

	return bcdToInt(buffRevers) / 1000
}

func requestAFLevel(port serial.Port, p *civCommand) uint32 {
	return (bcdToInt(append(make([]byte, 0), commandSend(port, p, AF)[1:]...)) * 100) / 254
}

func requestSQLLevel(port serial.Port, p *civCommand) uint32 {
	return (bcdToInt(append(make([]byte, 0), commandSend(port, p, SQL)[1:]...)) * 100) / 254
}

func requestRFLevel(port serial.Port, p *civCommand) uint32 {
	return bcdToInt(append(make([]byte, 0), commandSend(port, p, RF)[1:]...))
}

func IC78connect(port serial.Port, serialAcces *sync.Mutex) {
	serialAcces.Lock()
	fmt.Println("IC78 Connect")
	port.ResetInputBuffer()
	myic78civCommand := newIc78civCommand(requestTransiverAddr(port))
	fmt.Printf("Transiver Addr: %#x \n", myic78civCommand.transiverAddr)
	fmt.Printf("Transiver Freque: %d Hz \n", requestFreque(port, myic78civCommand))
	fmt.Println("Transiver Mode:", requestMode(port, myic78civCommand))
	fmt.Println("ATT status:", requestATT(port, myic78civCommand))
	fmt.Printf("AF level: %d % \n", requestAFLevel(port, myic78civCommand))
	fmt.Printf("RF level: %d \n", requestRFLevel(port, myic78civCommand))
	fmt.Printf("SQL level: %d % \n", requestSQLLevel(port, myic78civCommand))
	fmt.Println("Preamp status:", requestPreamp(port, myic78civCommand))
	//setFreque(30569)

	serialAcces.Unlock()
}

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
		fmt.Println("HELLO")
		port.ResetInputBuffer()
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
	if len(buff) > 2 {
		buffRevers := make([]byte, len(buff))
		j := 0
		for i := len(buff) - 1; i > -1; i-- {
			buffRevers[j] = buff[i]
			j++
		}
		return bcd.ToUint32(buffRevers)
	}
	return bcd.ToUint32(buff)
}

func requestTransiverAddr(port serial.Port) byte {
	requestTAddres := []byte{byte(preambleCmd), byte(preambleCmd), 0x00, byte(controllerAddrCmd), byte(readAddrCmd), 0x00, byte(endMsgCmd)}
	correctMsg := false
	buff := make([]byte, 30)
	var transiverAddr byte
	n := 0
	attempt := 0
	for !correctMsg {
		if attempt > 20 {
			break
		}
		attempt++
		port.ResetInputBuffer()
		time.Sleep(time.Duration(100) * time.Millisecond)
		serialDataExchange.WriteSerialPort(port, requestTAddres)
		time.Sleep(time.Duration(100) * time.Millisecond)
		_ = serialDataExchange.ReadSerialPort(port, buff)
		for _, value := range buff {
			if value == 0xfd {
				n++
			}
		}
		if n < 2 {
			n = 0
			for i, _ := range buff {
				buff[i] = 0x00
			}
		} else {
			correctMsg = true
		}
	}
	for i := 0; i < n; i++ {
		idx := slices.Index(buff, 0xfd)
		if idx != -1 {
			if bytes.Equal(buff[:idx+1], requestTAddres[:len(requestTAddres)]) {
				buff = buff[idx+1 : len(buff)]
			} else {
				transiverAddr = buff[idx-1]
			}

		}
	}
	if attempt > 99 {
		fmt.Println("attempt>100!")
	}

	return transiverAddr
}

func requestMode(port serial.Port, p *civCommand) string {
	correctMsg := false
	buff := make([]byte, 30)
	var mode string
	var modeByte byte
	n := 0
	for !correctMsg {
		port.ResetInputBuffer()
		time.Sleep(time.Duration(100) * time.Millisecond)
		serialDataExchange.WriteSerialPort(port, p.requestMode)
		time.Sleep(time.Duration(100) * time.Millisecond)
		_ = serialDataExchange.ReadSerialPort(port, buff)
		for _, value := range buff {
			if value == 0xfd {
				n++
			}
		}
		if n < 2 {
			n = 0
			for i, _ := range buff {
				buff[i] = 0x00
			}
		} else {
			correctMsg = true
		}
	}
	for i := 0; i < n; i++ {
		idx := slices.Index(buff, p.endMsg)
		if idx != -1 {
			if bytes.Equal(buff[:idx+1], p.requestMode[:len(p.requestMode)]) {
				buff = buff[idx+1 : len(buff)]
			} else {
				modeByte = buff[idx-2]
			}

		}
	}
	switch modeByte {
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
	correctMsg := false
	buff := make([]byte, 30)
	var preamp string
	var preampByte byte
	n := 0
	for !correctMsg {
		port.ResetInputBuffer()
		time.Sleep(time.Duration(100) * time.Millisecond)
		serialDataExchange.WriteSerialPort(port, p.requestPreamp)
		time.Sleep(time.Duration(100) * time.Millisecond)
		_ = serialDataExchange.ReadSerialPort(port, buff)
		for _, value := range buff {
			if value == 0xfd {
				n++
			}
		}
		if n < 2 {
			n = 0
			for i, _ := range buff {
				buff[i] = 0x00
			}
		} else {
			correctMsg = true
		}
	}
	for i := 0; i < n; i++ {
		idx := slices.Index(buff, p.endMsg)
		if idx != -1 {
			if bytes.Equal(buff[:idx+1], p.requestPreamp[:len(p.requestPreamp)]) {
				buff = buff[idx+1 : len(buff)]
			} else {
				preampByte = buff[idx-1]
			}

		}
	}
	switch preampByte {
	case 0x00:
		preamp = "OFF"
	case 0x01:
		preamp = "P.AMP"
	}
	return preamp

}

func requestATT(port serial.Port, p *civCommand) string {
	correctMsg := false
	buff := make([]byte, 30)
	var att string
	var attByte byte
	n := 0
	for !correctMsg {
		port.ResetInputBuffer()
		time.Sleep(time.Duration(100) * time.Millisecond)
		serialDataExchange.WriteSerialPort(port, p.requestATT)
		time.Sleep(time.Duration(100) * time.Millisecond)
		_ = serialDataExchange.ReadSerialPort(port, buff)
		for _, value := range buff {
			if value == 0xfd {
				n++
			}
		}
		if n < 2 {
			n = 0
			for i, _ := range buff {
				buff[i] = 0x00
			}
		} else {
			correctMsg = true
		}
	}
	for i := 0; i < n; i++ {
		idx := slices.Index(buff, p.endMsg)
		if idx != -1 {
			if bytes.Equal(buff[:idx+1], p.requestATT[:len(p.requestATT)]) {
				buff = buff[idx+1 : len(buff)]
			} else {
				attByte = buff[idx-1]
			}

		}
	}
	switch attByte {
	case 0x00:
		att = "NO"
	case 0x20:
		att = "YES"
	}
	return att
}

func requestFreque(port serial.Port, p *civCommand) uint32 {
	correctMsg := false
	buff := make([]byte, 30)
	freque := make([]byte, 5)
	n := 0
	for !correctMsg {
		port.ResetInputBuffer()
		time.Sleep(time.Duration(100) * time.Millisecond)
		serialDataExchange.WriteSerialPort(port, p.requestFreque)
		time.Sleep(time.Duration(100) * time.Millisecond)
		_ = serialDataExchange.ReadSerialPort(port, buff)
		for _, value := range buff {
			if value == 0xfd {
				n++
			}
		}
		if n < 2 {
			n = 0
			for i, _ := range buff {
				buff[i] = 0x00
			}
		} else {
			correctMsg = true
		}
	}
	for i := 0; i < n; i++ {
		idx := slices.Index(buff, p.endMsg)
		if idx != -1 {
			if bytes.Equal(buff[:idx+1], p.requestFreque[:len(p.requestFreque)]) {
				buff = buff[idx+1 : len(buff)]
			} else {
				freque = buff[idx-5 : idx]
			}

		}
	}
	return bcdToInt(freque) / 1000
}

func requestAFLevel(port serial.Port, p *civCommand) uint32 {
	correctMsg := false
	buff := make([]byte, 30)
	level := make([]byte, 2)
	n := 0
	for !correctMsg {
		port.ResetInputBuffer()
		time.Sleep(time.Duration(100) * time.Millisecond)
		serialDataExchange.WriteSerialPort(port, p.requestAFLevel)
		time.Sleep(time.Duration(100) * time.Millisecond)
		_ = serialDataExchange.ReadSerialPort(port, buff)
		for _, value := range buff {
			if value == 0xfd {
				n++
			}
		}
		if n < 2 {
			n = 0
			for i, _ := range buff {
				buff[i] = 0x00
			}
		} else {
			correctMsg = true
		}
	}
	for i := 0; i < n; i++ {
		idx := slices.Index(buff, p.endMsg)
		if idx != -1 {
			if bytes.Equal(buff[:idx+1], p.requestAFLevel[:len(p.requestAFLevel)]) {
				buff = buff[idx+1 : len(buff)]
			} else {
				level = buff[idx-2 : idx]
			}

		}
	}
	return (bcdToInt(level) * 100) / 254
}

func commandSend(port serial.Port, p *civCommand, c commandName) []byte {
	correctMsg := false
	readBuff := make([]byte, 30)
	dataBuff := make([]byte, 5)
	var arg []byte
	var cmd byte
	switch c {
	case FREQ:
		arg = p.requestFreque
		cmd = byte(readFreqCmd)
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
			}

		}
	}
	return dataBuff
}

func requestSQLLevel(port serial.Port, p *civCommand) uint32 {
	correctMsg := false
	buff := make([]byte, 30)
	level := make([]byte, 2)
	n := 0
	for !correctMsg {
		port.ResetInputBuffer()
		time.Sleep(time.Duration(100) * time.Millisecond)
		serialDataExchange.WriteSerialPort(port, p.requestSQLLevel)
		time.Sleep(time.Duration(100) * time.Millisecond)
		_ = serialDataExchange.ReadSerialPort(port, buff)
		for _, value := range buff {
			if value == 0xfd {
				n++
			}
		}
		if n < 2 {
			n = 0
			for i, _ := range buff {
				buff[i] = 0x00
			}
		} else {
			correctMsg = true
		}
	}
	for i := 0; i < n; i++ {
		idx := slices.Index(buff, p.endMsg)
		if idx != -1 {
			if bytes.Equal(buff[:idx+1], p.requestSQLLevel[:len(p.requestSQLLevel)]) {
				buff = buff[idx+1 : len(buff)]
			} else {
				level = buff[idx-2 : idx]
			}

		}
	}
	return (bcdToInt(level) * 100) / 254
}

func requestRFLevel(port serial.Port, p *civCommand) uint32 {
	correctMsg := false
	buff := make([]byte, 30)
	level := make([]byte, 2)
	n := 0
	for !correctMsg {
		port.ResetInputBuffer()
		time.Sleep(time.Duration(100) * time.Millisecond)
		serialDataExchange.WriteSerialPort(port, p.requestRFLevel)
		time.Sleep(time.Duration(100) * time.Millisecond)
		_ = serialDataExchange.ReadSerialPort(port, buff)
		for _, value := range buff {
			if value == 0xfd {
				n++
			}
		}
		if n < 2 {
			n = 0
			for i, _ := range buff {
				buff[i] = 0x00
			}
		} else {
			correctMsg = true
		}
	}
	for i := 0; i < n; i++ {
		idx := slices.Index(buff, p.endMsg)
		if idx != -1 {
			if bytes.Equal(buff[:idx+1], p.requestRFLevel[:len(p.requestRFLevel)]) {
				buff = buff[idx+1 : len(buff)]
			} else {
				level = buff[idx-2 : idx]
			}

		}
	}

	return bcdToInt(level)
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

package ic78civCmd

import (
	"bytes"
	"errors"
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
		var err error
		var adr byte
		adr, err = requestTransiverAddr(port)
		if err != nil {
			for err != nil {
				fmt.Println("error:", err)
				adr, err = requestTransiverAddr(port)
				time.Sleep(50 * time.Millisecond)
			}
		}
		fmt.Printf("transiver connected, addr:= %#x \n", adr)
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
	return bcd.ToUint32(buff)
}

func commandSend(port serial.Port, p *civCommand, c commandName) ([]byte, error) {
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
	attempt := 0
	for !correctMsg {
		if attempt < 10 {
			attempt++
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
		} else {
			return readBuff, errors.New("can't connect Transiver")
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
	return dataBuff, nil
}

func requestTransiverAddr(port serial.Port) (byte, error) {
	addr, er := commandSend(port, nil, TADDR)
	if er != nil {
		return 0x00, er
	} else {
		return addr[1], nil
	}

}

func requestMode(port serial.Port, p *civCommand) (string, error) {
	buff, err := commandSend(port, p, MODE)
	var mode string
	if err != nil {
		return "error", err
	}
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
	return mode, err
}

func requestPreamp(port serial.Port, p *civCommand) (string, error) {
	buff, err := commandSend(port, p, PREAMP)
	var preamp string
	if err != nil {
		return "error", err
	}
	buff = append(make([]byte, 0), buff[1:]...)
	switch buff[0] {
	case 0x00:
		preamp = "OFF"
	case 0x01:
		preamp = "P.AMP"
	}
	return preamp, err
}

func requestATT(port serial.Port, p *civCommand) (string, error) {
	buff, err := commandSend(port, p, ATT)
	var att string
	if err != nil {
		return "error", err
	}
	switch buff[0] {
	case 0x00:
		att = "NO"
	case 0x20:
		att = "YES"
	}
	return att, err
}

func requestFreque(port serial.Port, p *civCommand) (uint32, error) {
	buff, err := commandSend(port, p, FREQ)
	if err != nil {
		return 0, err
	}
	buffRevers := make([]byte, len(buff))
	j := 0
	for i := len(buff) - 1; i > -1; i-- {
		buffRevers[j] = buff[i]
		j++
	}
	return bcdToInt(buffRevers) / 1000, err
}

func requestAFLevel(port serial.Port, p *civCommand) (uint32, error) {
	buff, err := commandSend(port, p, AF)
	if err != nil {
		return 0, err
	}
	buff = append(make([]byte, 0), buff[1:]...)
	return (bcdToInt(buff) * 100) / 254, err
}

func requestSQLLevel(port serial.Port, p *civCommand) (uint32, error) {
	buff, err := commandSend(port, p, SQL)
	if err != nil {
		return 0, err
	}
	buff = append(make([]byte, 0), buff[1:]...)
	return (bcdToInt(buff) * 100) / 254, err
}

func requestRFLevel(port serial.Port, p *civCommand) (uint32, error) {
	buff, err := commandSend(port, p, RF)
	if err != nil {
		return 0, err
	}
	buff = append(make([]byte, 0), buff[1:]...)
	return bcdToInt(buff), err
}

func IC78connect(port serial.Port, serialAcces *sync.Mutex) error {
	serialAcces.Lock()
	fmt.Println("IC78 Connect")
	port.ResetInputBuffer()
	var myic78civCommand *civCommand
	addr, err := requestTransiverAddr(port)
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		myic78civCommand = newIc78civCommand(addr)
		fmt.Printf("Transiver Addr: %#x \n", myic78civCommand.transiverAddr)
	}
	freq, err := requestFreque(port, myic78civCommand)
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Printf("Transiver Freque: %d Hz \n", freq)
	}
	mode, err := requestMode(port, myic78civCommand)
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Println("Transiver Mode:", mode)
	}
	att, err := requestATT(port, myic78civCommand)
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Println("ATT status:", att)
	}
	af, err := requestAFLevel(port, myic78civCommand)
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Printf("AF level: %d % \n", af)
	}
	rf, err := requestRFLevel(port, myic78civCommand)
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Printf("RF level: %d % \n", rf)
	}
	sql, err := requestSQLLevel(port, myic78civCommand)
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Printf("SQL level: %d % \n", sql)
	}
	preamp, err := requestPreamp(port, myic78civCommand)
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Println("Preamp status:", preamp)
	}

	serialAcces.Unlock()
	return nil

	//setFreque(30569)

}

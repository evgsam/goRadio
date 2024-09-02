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

func setAfRfSql(port serial.Port, p *civCommand, c commandName, level int) error {
	levelBuf := intToArr((level*255)/100, 4)
	for len(levelBuf) < 4 {
		levelBuf = addElementToFirstIndex(levelBuf, 0)
	}
	var subcmd byte
	switch c {
	case af:
		subcmd = byte(afSubCmd)
	case rf:
		subcmd = byte(rfSubCmd)
	case sql:
		subcmd = byte(sqlSubCmd)
	}
	buf := make([]byte, 7)
	buf = []byte{byte(preambleCmd), byte(preambleCmd), p.transiverAddr, byte(controllerAddrCmd), byte(afrfsqlCmd),
		subcmd, byteArrToBCD(levelBuf, 2)[1], byteArrToBCD(levelBuf, 2)[0], byte(endMsgCmd)}
	_, readBuff, err := sendDataToTransiver(port, buf)
	if err != nil {
		return err
	}
	if slices.Index(readBuff, byte(ngCode)) > 0 {
		return errors.New("transceiver sent NG")
	}
	return nil
}

func setMode(port serial.Port, p *civCommand, mode string) error {
	var arg byte
	switch mode {
	case "LSB":
		arg = 0x00
	case "USB":
		arg = 0x01
	case "AM":
		arg = 0x02
	case "RTTY":
		arg = 0x04
	case "CW":
		arg = 0x07
	}
	buf := make([]byte, 7)
	buf = []byte{byte(preambleCmd), byte(preambleCmd), p.transiverAddr, byte(controllerAddrCmd), byte(setModeCmd), arg, byte(endMsgCmd)}
	_, readBuff, err := sendDataToTransiver(port, buf)
	if err != nil {
		return err
	}
	if slices.Index(readBuff, byte(ngCode)) > 0 {
		return errors.New("transceiver sent NG")
	}
	return nil
}

func setFreque(port serial.Port, p *civCommand, freq int) error {
	freqBuf := intFreqToBcdArr(freq)
	buf := make([]byte, 11)
	buf = []byte{byte(preambleCmd), byte(preambleCmd), p.transiverAddr, byte(controllerAddrCmd), byte(setFreqCmd)}
	for i := 0; i < len(freqBuf); i++ {
		buf = append(buf, freqBuf[i])

	}
	buf = append(buf, byte(endMsgCmd))
	_, readBuff, err := sendDataToTransiver(port, buf)
	if err != nil {
		return err
	}
	if slices.Index(readBuff, byte(ngCode)) > 0 {
		return errors.New("transceiver sent NG")
	}
	return nil
}

func setPreamp(port serial.Port, p *civCommand, preamp string) error {
	var cmd byte
	switch preamp {
	case "OFF":
		cmd = 0x00
	case "P.AMP":
		cmd = 0x01
	}
	buf := make([]byte, 7)
	buf = []byte{byte(preambleCmd), byte(preambleCmd), p.transiverAddr, byte(controllerAddrCmd), byte(preampCmd), byte(preampSubCmd), cmd, byte(endMsgCmd)}
	_, readBuff, err := sendDataToTransiver(port, buf)
	if err != nil {
		return err
	}
	if slices.Index(readBuff, byte(ngCode)) > 0 {
		return errors.New("transceiver sent NG")
	}
	return nil

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

func sendDataToTransiver(port serial.Port, arg []byte) (int, []byte, error) {
	n := 0
	attempt := 0
	readBuff := make([]byte, 30)
	correctMsg := false
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
			return n, readBuff, errors.New("can't connect Transiver")
		}
	}
	return n, readBuff, nil
}

func commandSend(port serial.Port, p *civCommand, c commandName) ([]byte, error) {
	readBuff := make([]byte, 30)
	dataBuff := make([]byte, 7)
	var arg []byte
	var cmd byte
	switch c {
	case freqRead:
		arg = p.requestFreque
		cmd = byte(readFreqCmd)
	case taddr:
		arg = []byte{byte(preambleCmd), byte(preambleCmd), 0x00, byte(controllerAddrCmd), byte(readAddrCmd), 0x00, byte(endMsgCmd)}
		cmd = byte(readAddrCmd)
	case mode:
		arg = p.requestMode
		cmd = byte(readModeCmd)
	case att:
		arg = p.requestATT
		cmd = byte(attCmd)
	case af:
		arg = p.requestAFLevel
		cmd = byte(afrfsqlCmd)
	case rf:
		arg = p.requestRFLevel
		cmd = byte(afrfsqlCmd)
	case sql:
		arg = p.requestSQLLevel
		cmd = byte(afrfsqlCmd)
	case preamp:
		arg = p.requestPreamp
		cmd = byte(preampCmd)
	}
	n, readBuff, err := sendDataToTransiver(port, arg)
	if err != nil {
		return readBuff, err
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
	addr, er := commandSend(port, nil, taddr)
	if er != nil {
		return 0x00, er
	} else {
		return addr[1], nil
	}

}

func requestMode(port serial.Port, p *civCommand) (string, error) {
	buff, err := commandSend(port, p, mode)
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
	buff, err := commandSend(port, p, preamp)
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
	buff, err := commandSend(port, p, att)
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
	buff, err := commandSend(port, p, freqRead)
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
	buff, err := commandSend(port, p, af)
	if err != nil {
		return 0, err
	}
	buff = append(make([]byte, 0), buff[1:]...)
	return (bcdToInt(buff) * 100) / 254, err
}

func requestSQLLevel(port serial.Port, p *civCommand) (uint32, error) {
	buff, err := commandSend(port, p, sql)
	if err != nil {
		return 0, err
	}
	buff = append(make([]byte, 0), buff[1:]...)
	return (bcdToInt(buff) * 100) / 254, err
}

func requestRFLevel(port serial.Port, p *civCommand) (uint32, error) {
	buff, err := commandSend(port, p, rf)
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
	afLevel, err := requestAFLevel(port, myic78civCommand)
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Printf("AF level: %d % \n", afLevel)
	}
	rfLevel, err := requestRFLevel(port, myic78civCommand)
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Printf("RF level: %d % \n", rfLevel)
	}
	sqlLevel, err := requestSQLLevel(port, myic78civCommand)
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Printf("SQL level: %d % \n", sqlLevel)
	}
	preamp, err := requestPreamp(port, myic78civCommand)
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Println("Preamp status:", preamp)
	}
	err = setFreque(port, myic78civCommand, 3501)
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Println("freque set")
	}
	err = setMode(port, myic78civCommand, "AM")
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Println("mode set")
	}
	err = setAfRfSql(port, myic78civCommand, af, 93)
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Println("af level set")
	}

	err = setAfRfSql(port, myic78civCommand, rf, 99)
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Println("rf level set")
	}
	err = setAfRfSql(port, myic78civCommand, sql, 69)
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Println("sql level set")
	}

	err = setPreamp(port, myic78civCommand, "P.AMP")
	if err != nil {
		serialAcces.Unlock()
		return err
	} else {
		fmt.Println("P.AMP set")
	}

	serialAcces.Unlock()
	return nil

}

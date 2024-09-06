package ic78civCmd

import (
	"bytes"
	"errors"
	"goRadio/serialDataExchange"
	"slices"
	"time"

	"go.bug.st/serial"
)

func sendDataToTransiver(port serial.Port, arg []byte) (int, []byte, error) {
	n := 0
	attempt := 0
	readBuff := make([]byte, maxReadBuff)
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
	readBuff := make([]byte, maxReadBuff)
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
		if idxCmd > idxEnd {
			readBuff = readBuff[idxCmd+1 : len(readBuff)]
		}
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
